package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const testInputFile = "testdata/in.txt"

var testCases = []struct {
	want   string
	option string
}{
	{want: "testdata/want_default_pass.txt", option: ""},
	{want: "testdata/want_option_n_pass.txt", option: " -n=3"},
	{want: "testdata/want_option_o_pass.txt", option: " -o=1"},
	{want: "testdata/want_option_n_o_pass.txt", option: " -n=4 -o=3"},
}

func TestRun(t *testing.T) {
	for _, tc := range testCases {
		outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
		cli := &CLI{
			outStream: outStream,
			errStream: errStream,
		}

		want, err := ioutil.ReadFile(tc.want)
		if err != nil {
			t.Fatalf("failed to open want file: %s", err)
		}

		command := fmt.Sprintf("ghead%s %s", tc.option, testInputFile)
		_ = cli.Run(strings.Split(command, " "))

		if got := outStream.String(); got != string(want) {
			t.Fatalf("faild heading with args: %s", command)
		}
	}
}
