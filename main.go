package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

const usage = `contextualize -- prefix lines with a previously-found expression

Usage: contextualize [-h, --help] expression [filename]

parameters:
  -h, --help    Display this message
   expression   RE2-style regex. Matched values will be used as the prefix
   filename     Optional filename or glob. If not provided, stdin will be read

The expression is only searched in the first 65535 characters of a line
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: expression required")
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	if args[0] == "-h" || args[0] == "--help" {
		fmt.Println(usage)
		return
	}

	exp, err := regexp.Compile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid expression: %s\n", err.Error())
		os.Exit(1)
	}
	args = args[1:]

	if len(args) > 0 {
		for _, g := range args {
			filenames, err := filepath.Glob(g)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: invalid glob '%s'\n", g)
				os.Exit(1)
			}
			for _, fn := range filenames {
				file, err := os.Open(fn)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed to open file '%s': %s\n", fn, err.Error())
					os.Exit(1)
				}
				if err := doContextualize(file, exp); err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: Failed parsing file '%s': %s\n", fn, err.Error())
				}
			}
		}
	} else {
		if err := doContextualize(os.Stdin, exp); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed stdin: %s\n", err.Error())
		}
	}
}

var (
	divider = []byte(": ")
	newline = []byte("\n")
)

func doContextualize(src io.Reader, exp *regexp.Regexp) error {
	var prefix []byte
	bsrc := bufio.NewReader(src)
	var line []byte
	var isPrefix bool
	var err error
	var lastWasPrefix = false
	for {
		line, isPrefix, err = bsrc.ReadLine()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		if !lastWasPrefix {
			if m := exp.FindSubmatch(line); m != nil {
				if len(m) > 1 {
					prefix = m[1]
				} else {
					prefix = m[0]
				}
			}

			os.Stdout.Write(prefix)
			os.Stdout.Write(divider)
		}
		os.Stdout.Write(line)
		if !isPrefix {
			os.Stdout.Write(newline)
		}
	}
}
