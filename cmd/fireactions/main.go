package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stderr))
}

// run executes the command tree and returns the process exit code: 0 when the command
// succeeded, 1 when it reported an error.
//
// The exit code is not cosmetic. A CLI that prints its failure and exits 0 makes every
// caller read a failed start as a success: a `Restart=always` unit reports a healthy
// process, and a scripted gate (or a CI step) cannot detect the failure at all.
func run(args []string, stderr io.Writer) int {
	cmd := NewRootCommand()
	cmd.SetArgs(args)

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(stderr, "Error executing command: %s\n", err.Error())

		return 1
	}

	return 0
}
