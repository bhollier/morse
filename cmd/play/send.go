package main

import (
	"fmt"
	"github.com/bhollier/morse"
)

func SendSignals(sc chan<- morse.Signal, c morse.Code) {
	for _, s := range c {
		sc <- s
		if *printMorse {
			fmt.Print(s)
		}
	}
	if *printMorse && len(c) > 0 {
		fmt.Println()
	}
}
