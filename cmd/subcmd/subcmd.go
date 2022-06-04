// Package subcmd provides convenience functions for command line tools
// with "subcommands", similar to git's init, clone, add, etc.
package subcmd

import (
	"fmt"
	"os"
	"strings"
)

type SubCmd interface {
	Name() string
	Run(args []string)
}

func subCmdOptionsStr(subCmds []SubCmd) string {
	sb := strings.Builder{}
	i := 0
	for _, subCommand := range subCmds {
		sb.WriteRune('\'')
		sb.WriteString(subCommand.Name())
		sb.WriteRune('\'')
		if i+2 < len(subCmds) {
			sb.WriteString(", ")
		} else if i+1 < len(subCmds) {
			sb.WriteString(" or ")
		}
		i++
	}
	return sb.String()
}

// PickAndRun takes a list of possible subcommands,
// and runs the one that was specified by the first argument to the program,
// according to SubCmd.Name()
func PickAndRun(subCmds []SubCmd) error {
	subCommandsStr := subCmdOptionsStr(subCmds)

	args := os.Args[1:]
	if len(os.Args) < 1 {
		return fmt.Errorf("no subcommand, expected %s\n", subCommandsStr)
	}

	subCmdsMap := make(map[string]SubCmd, len(subCmds))
	for _, subCommand := range subCmds {
		subCmdsMap[subCommand.Name()] = subCommand
	}

	subCmd, ok := subCmdsMap[strings.ToLower(args[0])]
	if !ok {
		return fmt.Errorf("unknown subcommand '%s', expected %s\n", subCommandsStr, args[0])
	}
	subCmd.Run(args[1:])
	return nil
}
