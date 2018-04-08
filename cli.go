package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var helpText = `
Usage: ghead [-o=<num>] [-n=<num>] [file ...]

ghead is a tool to display lines of a file or standard input from head.

Options:

	-n=number Number of lines to show.

	-o=number Offset to heading.
`

const (
	defaultMaxSize = 10
	defaultOffset  = 0
)

const (
	ExitCodeOk = iota
	ExitCodeError
	ExitCodeParseFlagsError
)

var (
	maxSize int
	offset  int
)

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI) Run(args []string) int {
	flags := flag.NewFlagSet("head", flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprint(cli.errStream, helpText)
	}
	flags.IntVar(&maxSize, "n", defaultMaxSize, "Number of lines to show")
	flags.IntVar(&offset, "o", defaultOffset, "Offset to heading")
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	if flags.NArg() == 0 {
		err := head(os.Stdin, cli.outStream)
		if err != nil {
			fmt.Fprintln(cli.errStream, err)
			return ExitCodeError
		}
	} else {
		filenames := flags.Args()
		for i, filename := range filenames {
			if i >= 1 {
				fmt.Fprintln(cli.outStream)
			}

			input, err := os.Open(filename)
			if err != nil {
				fmt.Fprintln(cli.errStream, err)
				continue
			}
			defer input.Close()

			if flags.NArg() > 1 {
				fmt.Fprintf(cli.outStream, "==> %s <==\n", filename)
			}

			err = head(input, cli.outStream)
			if err != nil {
				fmt.Fprintln(cli.errStream, err)
				return ExitCodeError
			}
		}
	}

	return ExitCodeOk
}

func head(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for i := 1; scanner.Scan(); i++ {
		if i <= offset {
			continue
		}

		err := scanner.Err()
		if err != nil {
			return err
		}

		fmt.Fprintln(w, scanner.Text())

		if i == maxSize+offset {
			return nil
		}
	}

	return nil
}
