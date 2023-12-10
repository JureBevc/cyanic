package main

import (
	"log/slog"
	"os"

	"github.com/JureBevc/cyanic/actions"
)

func main() {

	if len(os.Args) < 2 {
		actions.PrintHelp()
		return
	}

	command := os.Args[1]
	err := actions.ParseCommand(command, os.Args[2:])
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
