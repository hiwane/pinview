package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/hiwane/pinview/internal/draw"
	"github.com/hiwane/pinview/internal/input"
	"github.com/hiwane/pinview/internal/pager"
	"github.com/hiwane/pinview/internal/term"
)

func getLines() ([]string, error) {
	var scanner *bufio.Scanner
	// 引数なし → stdin
	if flag.NArg() == 0 {
		// パイプ...?
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "open file error:", flag.Arg(0))
			return nil, err
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
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		return nil, err
	}

	return lines, nil
}

func _main(header, footer int, showRuler bool) error {

	lines, err := getLines()
	if err != nil {
		fmt.Fprintln(os.Stderr, "getLines()", err)
		return err
	}

	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "open /dev/tty:", err)
		return err
	}

	in, err := input.New(tty)
	if err != nil {
		fmt.Fprintln(os.Stderr, "input init error:", err)
		return err
	}
	defer in.Close()

	h, err := term.GetHeight(tty)
	if err != nil {
		fmt.Fprintln(os.Stderr, "terminal init error:", err)
		return err
	}
	model := pager.NewModel(lines, h-1)
	model.SetRuler(showRuler)
	model.SetHeader(header)
	model.SetFooter(footer)
	for {
		term.ViewClearScreen(os.Stdout)

		lines := model.View()
		draw.Draw(lines)

		key, err := in.ReadRune()
		if err != nil {
			return err
		}

		if model.Update(key) {
			return nil
		}
	}
}

func main() {
	header := flag.Int("H", 1, "number of header lines to pin")
	footer := flag.Int("F", 0, "number of footer lines to pin (default 0)")

	var showRuler bool
	flag.BoolVar(&showRuler, "ruler", false, "show ruler")
	flag.Parse()

	if !term.IsInteractive() {
		fmt.Fprintln(os.Stderr, "pinview: stdout is not a terminal")
		os.Exit(1)
	}
	if flag.NArg() == 0 && term.IsTTY(os.Stdin) {
		fmt.Fprintln(os.Stderr, "pinview: no input file and stdin is a terminal")
		os.Exit(1)
	}
	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "pinview: too many arguments")
		os.Exit(1)
	}

	err := _main(*header, *footer, showRuler)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
