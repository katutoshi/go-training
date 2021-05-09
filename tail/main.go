package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	ExitCodeOK    = 0
	ExitCodeError = 1

	Default = 10
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	if terminal.IsTerminal(0) {
		return ExitCodeOK
	}

	output, err := Tail(os.Stdin, Default)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	if err := Print(os.Stdout, output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func Tail(r io.Reader, n int) ([]string, error) {
	// tailする行数の何倍かをbufferとして持つ。倍率はテキトー
	bs := n * 10
	buf := make([]string, 0, bs)

	br := bufio.NewReader(r)
	for {
		var b []byte
		var err error
		var line string
		prefix := true
		for prefix && err == nil {
			b, prefix, err = br.ReadLine()
			line += string(b)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		buf = append(buf, line)
		if bs == len(buf) {
			tmp := make([]string, 0, bs)
			copy(tmp, buf[bs-n:])
			buf = tmp
		}
	}

	var offset int
	if len(buf) <= n {
		offset = 0
	} else {
		offset = len(buf) - n
	}

	return buf[offset:], nil
}

func Print(w io.Writer, text []string) error {
	for _, s := range text {
		if _, err := fmt.Fprintf(w, "%s\n", s); err != nil {
			return err
		}
	}

	return nil
}
