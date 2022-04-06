package main

import (
	"fmt"
	"os"

	"github.com/matsuyoshi30/germanium/cli"
)

var exit = os.Exit

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exit(1)
	}
}
