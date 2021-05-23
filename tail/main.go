package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	ExitCodeOK    = 0
	ExitCodeError = 1

	DefaultNumberOfLines = 10
)

func main() {
	os.Exit(Run(os.Args, os.Stdin, os.Stdout, os.Stderr))
}

func Run(args []string, r io.Reader, w io.Writer, ew io.Writer) int {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	var lineNumber int
	flags.IntVar(&lineNumber, "n", DefaultNumberOfLines, "-n=NUM\nnumber of lines")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintln(ew, err)
		return ExitCodeError
	}

	if terminal.IsTerminal(0) {
		return ExitCodeOK
	}

	output, err := Tail(r, lineNumber)
	if err != nil {
		fmt.Fprintln(ew, err)
		return ExitCodeError
	}

	if err := Print(w, output); err != nil {
		fmt.Fprintln(ew, err)
		return ExitCodeError
	}

	return ExitCodeOK
}

func Tail(r io.Reader, n int) ([]string, error) {
	// tailする行数の何倍かをバッファとして持つ。倍率はテキトー
	bs := n * 10
	buf := make([]string, 0, bs)
	br := bufio.NewReader(r)

	for {
		var b []byte
		var err error
		var line string
		isPrefix := true

		for isPrefix && err == nil {
			// Scannerだと行ごとに読める読める最大値が決まっているのでReadLineする
			// isPrefixがtrueだと行がバッファに収まりきってないのでisPrefixがfalseになるまで読む
			b, isPrefix, err = br.ReadLine()
			line += string(b)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		buf = append(buf, line)

		// メモリ節約のために一定サイズでバッファを作り直す
		// allocationの回数を減らすためにスライスをcopyしている
		if bs <= len(buf) {
			tmp := make([]string, n, bs)
			copy(tmp, buf[bs-n:])
			buf = tmp
		}
	}

	// バッファは取りたい末尾の行数より多めに作ってあるので取りたい範囲だけ取るようにする
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
