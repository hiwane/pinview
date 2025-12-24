package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	headerLines := flag.Int("n", 1, "number of header lines to pin")
	flag.IntVar(headerLines, "header", 1, "number of header lines to pin")
	var showRuler bool
	flag.BoolVar(&showRuler, "ruler", false, "show ruler")
	flag.Parse()

	var scanner *bufio.Scanner

	// 引数なし → stdin
	if flag.NArg() == 0 {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
	}

	// 長い行対策
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pager, err := New(lines, *headerLines)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer pager.Close()
	pager.SetRuler(showRuler)
	pager.Run()
}
