package main

import (
	"os"

	"github.com/ankit-songara/google-text-speech/cmd"
)

func main() {
	cli := &cmd.CLI{ErrStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
