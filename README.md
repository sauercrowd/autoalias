# autoalias

_This is a very quick first draft_

Sometimes it's hard to see the forest because of all the trees.

autoalias will automatically generate aliases for you of commands you use often.

```
$ hello world
$ hello world
[...]
$ hello world
Hey! I created a new alias for you
    hw=hello world
$ hw
```

## Setup

1. `go build`
2. Add autoalias binary to one of your PATH directories
3. Adapt setup.zsh
4. source setup.zsh from your .zshrc

Currently a command is considered to be used often if it's used 20 times within a week.
