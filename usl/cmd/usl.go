package main

import (
	"fmt"
	"os"

	"github.com/omu/zoo/usl"
)

const Usage = "Usage: usl USL [attributes...]"

func die(message ...interface{}) {
	fmt.Fprintln(os.Stderr, message...)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		die(Usage)
	}

	rawusl := os.Args[1]

	us, err := usl.Parse(rawusl)

	if err != nil {
		die(err)
	}

	if len(os.Args[2:]) == 0 {
		us.Dump(os.Args[2:]...)
	} else {
		us.Print(os.Args[2:]...)
	}
}
