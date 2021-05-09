package main

import (
	"bytes"
	_ "embed"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTail(t *testing.T) {
	given := `11111
22222
33333
44444
55555
66666
77777
88888
99999
00000`
	want := []string{"99999", "00000"}

	got, err := tailD(strings.NewReader(given), 2)
	if err != nil {
		t.Fatal(err)
	}

	ret := cmp.Diff(got, want)
	if len(ret) != 0 {
		t.Fatal("一致しない")
	}
}

func TestPrint(t *testing.T) {
	given := []string{"A", "B", "C"}
	want := `A
B
C
`
	buf := &bytes.Buffer{}
	if err := print(buf, given); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if got != want {
		t.Fatal("一致しない")
	}
}

//go:embed testdata/testdata_short.csv
var testData []byte

func BenchmarkTail(b *testing.B) {
	benchmarks := []struct {
		name     string
		tailFunc func(io.Reader, int) ([]string, error)
	}{
		{"Scanner", tail},
		{"ReadLine+Shift", tailB},
		{"ReadLine+Shift都度", tailC},
		{"ReadLine+Shift明示的", tailD},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				r := bytes.NewReader(testData)
				_, err := bm.tailFunc(r, 10)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}