package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	cli "github.com/urfave/cli/v2"
)

var app = cli.App{
	Name:  "queuer",
	Usage: "A queue rotator",
	Commands: []*cli.Command{
		{
			Name:  "new",
			Usage: "Create a new queue",
			Action: func(c *cli.Context) error {
				// prompt for location
				// TODO: path autocomplete would be great
				relPath := textLinePrompt("Path to store new queue? (dir must exist)")
				location, err := filepath.Abs(relPath)
				if err != nil {
					return fmt.Errorf("can't get absolute path: %s", err)
				}

				// check dir exists
				if _, err := os.Stat(filepath.Dir(location)); os.IsNotExist(err) {
					return fmt.Errorf("%s does not exist. dir must exist", filepath.Dir(location))
				}

				// prompt for name
				name := textLinePrompt("Name your queue")
				if name == "" {
					return fmt.Errorf("must provide name")
				}

				// read in state
				state, err := loadState()
				if err != nil {
					return fmt.Errorf("unable to load state: %s", err)
				}

				// save location in state
				if err = state.addQueue(name, location); err != nil {
					return fmt.Errorf("unable to add queue to state: %s", err)
				}

				// create file at location
				if err = createQueueFile(location); err != nil {
					return fmt.Errorf("unable to create queue file: %s", err)
				}

				// save state
				if err = state.save(); err != nil {
					return fmt.Errorf("unable to save state: %s", err)
				}

				return nil
			},
		},
		{
			Name:    "next",
			Aliases: []string{"n"},
			Usage:   "Move to next in queue",
			Action: func(c *cli.Context) error {
				args := c.Args()
				queueName := args.First()

				// read in state
				state, err := loadState()
				if err != nil {
					return fmt.Errorf("unable to load state: %s", err)
				}

				queueState, ok := state.getQueue(queueName)
				if !ok {
					return fmt.Errorf("queue <%s> not found", queueName)
				}

				// load queue
				queue, err := loadQueue(queueState.Location)
				if err != nil {
					return fmt.Errorf("unable to load queue: %s", err)
				}

				queueState.CurrentIndex++

				if queueState.CurrentIndex >= len(queue.data) {
					queueState.CurrentIndex = queueState.CurrentIndex % len(queue.data)
				}
				// always save, because we're incrementing
				state.save()

				// atIndex
				line := queue.atIndex(queueState.CurrentIndex)
				printNice(line)

				return nil
			},
		},
		{
			Name:    "current",
			Aliases: []string{"c"},
			Usage:   "View the current item",
			Action: func(c *cli.Context) error {
				args := c.Args()
				queueName := args.First()

				// read in state
				state, err := loadState()
				if err != nil {
					return fmt.Errorf("unable to load state: %s", err)
				}

				queueState, ok := state.getQueue(queueName)
				if !ok {
					return fmt.Errorf("queue <%s> not found", queueName)
				}

				// load queue
				queue, err := loadQueue(queueState.Location)
				if err != nil {
					return fmt.Errorf("unable to load queue: %s", err)
				}

				if queueState.CurrentIndex >= len(queue.data) {
					queueState.CurrentIndex = queueState.CurrentIndex % len(queue.data)
					state.save()
				}

				// atIndex
				line := queue.atIndex(queueState.CurrentIndex)
				printNice(line)

				return nil
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Usage:   "Edit a queue",
			Action: func(c *cli.Context) error {
				args := c.Args()
				queueName := args.First()

				// read in state
				state, err := loadState()
				if err != nil {
					return fmt.Errorf("unable to load state: %s", err)
				}

				queueState, ok := state.getQueue(queueName)
				if !ok {
					return fmt.Errorf("queue <%s> not found", queueName)
				}

				editor := os.Getenv("EDITOR")
				if editor == "" {
					return fmt.Errorf("your $EDITOR var must be set")
				}

				fmt.Println(editor, queueState.Location)

				// open file in editor
				cmd := exec.Command(editor, queueState.Location)

				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				return cmd.Run()
			},
		},
	},
}

func getQueueName(args []string) (string, bool) {
	if len(args) == 0 {
		return "", false
	}
	return args[0], true
}

func main() {
	if err := app.Run(os.Args); err != nil {
		println(err.Error())
	}
}
