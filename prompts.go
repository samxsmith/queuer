package main

import (
	"bufio"
	"fmt"
	"os"
)

func textLinePrompt(msg string) string {
	print(fmt.Sprintf("%v: ", msg))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	return input
}
