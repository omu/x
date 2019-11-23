package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/omu/zoo/usl"
)

// Usage should be commented
const usage = "Usage: usl USL [attributes...]"

func die(message ...interface{}) {
	fmt.Fprintln(os.Stderr, message...)
	os.Exit(1)
}

func wanted(defaultAttributes []string, attributes ...string) []string {
	if len(attributes) > 0 {
		return attributes
	}

	return defaultAttributes
}

// Dump should be commented
func Dump(us *usl.USL, attributes ...string) {
	m, ks := us.Map()

	for _, attribute := range wanted(ks, attributes...) {
		if value, ok := m[attribute]; ok && value != "" {
			fmt.Printf("%-24s %s\n", color.New(color.FgCyan, color.Bold).Sprint(attribute), value)
		}
	}
}

// Print should be commented
func Print(us *usl.USL, attributes ...string) {
	m, ks := us.Map()

	var values []string

	for _, attribute := range wanted(ks, attributes...) {
		if value, ok := m[attribute]; ok {
			values = append(values, value)
		}
	}

	fmt.Println(strings.Join(values, " "))
}

func main() {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		color.NoColor = true
	}

	if len(os.Args) < 2 {
		die(usage)
	}

	us, err := usl.Parse(os.Args[1])
	if err != nil {
		die(err)
	}

	if len(os.Args[2:]) == 0 {
		Dump(us, os.Args[2:]...)
	} else {
		Print(us, os.Args[2:]...)
	}
}
