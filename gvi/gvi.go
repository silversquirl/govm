package main

import (
	"flag"
	"fmt"
	"go.vktec.org.uk/govm"
	"io"
	"os"
)

func Main() int {
	flag.Parse()

	var input io.ReadSeeker
	if len(flag.Args()) > 0 {
		var err error
		input, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			return 1
		}
	} else {
		input = os.Stdin
	}

	vm := govm.New()

	vm.LoadFrom(input)
	if err := vm.Get("Main:"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := vm.Call(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(Main())
}
