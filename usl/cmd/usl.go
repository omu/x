package main

import (
	"fmt"
	"os"
	"strings"

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
	wanted := wanted(ks, attributes...)

	for _, attribute := range wanted {
		if value, ok := m[attribute]; ok && value != "" {
			fmt.Printf("%-16s %s\n", attribute, value)
		}
	}
}

// Print should be commented
func Print(us *usl.USL, attributes ...string) {
	m, ks := us.Map()
	wanted := wanted(ks, attributes...)

	var values []string
	for _, attribute := range wanted {
		if value, ok := m[attribute]; ok {
			values = append(values, value)
		}
	}

	fmt.Println(strings.Join(values[:], " "))
}

func main() {
	if len(os.Args) < 2 {
		die(usage)
	}

	rawusl := os.Args[1]

	us, err := usl.Parse(rawusl)

	if err != nil {
		die(err)
	}

	if len(os.Args[2:]) == 0 {
		Dump(us, os.Args[2:]...)
	} else {
		Print(us, os.Args[2:]...)
	}
}
