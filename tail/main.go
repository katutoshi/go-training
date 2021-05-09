package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
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

	text, err := tailD(os.Stdin, Default)
	if err != nil {
		return ExitCodeError
	}

	if err := print(os.Stdout, text); err != nil {
		return ExitCodeError
	}

	return ExitCodeOK
}

func tail(r io.Reader, n int) ([]string, error) {
	scanner := bufio.NewScanner(r)

	buf := make([]string, 0)
	for scanner.Scan() {
		buf = append(buf, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	var offset int
	if len(buf) <= n {
		offset = 0
	} else {
		offset = len(buf) - n
	}

	return buf[offset:], nil
}

func tailB(r io.Reader, n int) ([]string, error) {
	buf := make([]string, 0)
	reader := bufio.NewReader(r)
	for {
		var b []byte
		var err error
		var line string
		prefix := true
		for prefix && err == nil {
			b, prefix, err = reader.ReadLine()
			line += string(b)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		buf = append(buf, line)
		if n*10 < len(buf) {
			shift := math.Floor(float64(len(buf) / 2))
			buf = buf[int(shift):]
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

func tailC(r io.Reader, n int) ([]string, error) {
	buf := make([]string, 0)
	reader := bufio.NewReader(r)
	for {
		var b []byte
		var err error
		var line string
		prefix := true
		for prefix && err == nil {
			b, prefix, err = reader.ReadLine()
			line += string(b)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		buf = append(buf, line)
		if len(buf) > n {
			buf = buf[1:]
		}
	}

	return buf, nil
}

func tailD(r io.Reader, n int) ([]string, error) {
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

func print(w io.Writer, text []string) error {
	for _, s := range text {
		if _, err := fmt.Fprintf(w, "%s\n", s); err != nil {
			return err
		}
	}

	return nil
}
