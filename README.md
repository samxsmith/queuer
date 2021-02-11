# Queuer


### Requires
- Go

### Install

```
git clone git@github.com:samxsmith/queuer.git
cd queuer
make install
```

### Usage
For first use, you need to init the state file:
```
touch ~/.queuer_state
```

#### Create a new queue
```
queuer new
> Path to store new queue? (dir must exist): ./my_quotes_queue
> Name your queue: quotes
```

You can have as many queues as you want.
I use it for quotes, a daily question to journal about & a friend to check in on.


#### Read from the front of the queue e.g. today's quote
```
queuer current quotes
OR
queuer c quotes
```

#### Move to the next in line
```
queuer next quotes
OR
queuer n quotes
```

#### Adding to the queue
The queue file is plain text.
Each line is treated as an element in the queue.

To add a new element, just add a line.

You can edit the queue with any text editor or run the following command to use your default:
```
queuer edit quotes
```

### Future Features
- Random element from queue
