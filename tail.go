// Package tail only deals with files, this is incomplete, it needs to be able to deal with streams and pipes.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	filename = flag.String("file", "", "filename to tail")
	lines    = flag.Int("lines", 10, "number of lines to tail")
)

var diskBuf int64 = 4000

func Tail(m int, f *os.File) ([]string, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	// filesize
	size := stat.Size()
	buf := diskBuf
	count := 0
	var offset int64 = 0
	for {
		offset = size - buf
		if offset < 0 {
			offset = 0
			break
		}
		// initial read to the offset
		f.Seek(offset, 0)
		// create a scanner
		s := bufio.NewScanner(f)
		for s.Scan() {
			count++
		}
		buf = buf + diskBuf
		if count > m {
			fmt.Printf("count is %d\n", count)
			break
		}
	}
	// should now have the needed buffer size.
	f.Seek(offset, 0)
	// need to store lines in a slice of strings.
	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	return lines, nil
}

func main() {
	flag.StringVar(filename, "f", "", "filename to tail")
	flag.IntVar(lines, "n", 10, "ibleh")
	flag.Parse()
	if *filename == "" || *lines == 0 {
		fmt.Println("Please supply an argument")
		os.Exit(-1)
	}
	if *lines < 0 {
		log.Fatalf("%d lines cannot be less than zero", *lines)
	}
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("cannot handle file %q, %s", *filename, err)
	}
	defer f.Close()
	out, err := Tail(*lines, f)
	if err != nil {
		log.Fatal(err)
	}
	if *lines > len(out) {
		*lines = len(out)
	}
	for _, l := range out[:*lines] {
		fmt.Println(l)
	}
}
