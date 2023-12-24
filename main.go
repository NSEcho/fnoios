package main

import (
	_ "embed"
	"fmt"
	"github.com/nsecho/fnoios/cmd"
	"os"
)

//go:embed script/script.js
var script string

func main() {
	if err := cmd.Execute(script); err != nil {
		fmt.Fprintf(os.Stderr, "errorr: %v\n", err)
		os.Exit(1)
	}
}
