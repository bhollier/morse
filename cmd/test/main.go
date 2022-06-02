package main

import (
	"flag"
	"fmt"
	"github.com/bhollier/morse/cmd/test/receive"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(flag.CommandLine.Output(), "no subcommand, expected 'receive' or 'send'")
		os.Exit(2)
	}

	switch args[0] {
	case "receive":
		receive.Main(args[1:])
	case "send":
		fmt.Fprintln(flag.CommandLine.Output(), "todo")
		os.Exit(2)
	default:
		fmt.Fprintf(flag.CommandLine.Output(), "unknown subcommand '%s', "+
			"expected either 'receive' or 'send'\n", args[0])
		os.Exit(2)
	}
}
