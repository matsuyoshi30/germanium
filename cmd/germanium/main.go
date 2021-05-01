package main

import (
	"fmt"
	"os"

	"github.com/matsuyoshi30/germanium/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}
