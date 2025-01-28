package main

import (
	_ "embed"
	"fmt"
	"github.com/nsecho/fnoios/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "errorr: %v\n", err)
		os.Exit(1)
	}
}
