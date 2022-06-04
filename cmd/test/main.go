package main

import (
	"flag"
	"fmt"
	"github.com/bhollier/morse/cmd/subcmd"
	"github.com/bhollier/morse/cmd/test/receive"
	"os"
)

func main() {
	err := subcmd.PickAndRun([]subcmd.SubCmd{receive.SubCmd})
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), err.Error())
		os.Exit(2)
	}
}
