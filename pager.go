package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type Pager struct {
	lines []string
	pin   int

	scroll int
	height int
	ruler  bool

	tty      *os.File
	oldState *term.State
}

func (p *Pager) SetRuler(on bool) {
	p.ruler = on
}

func New(lines []string, pin int) (*Pager, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	oldState, err := term.MakeRaw(int(tty.Fd()))
	if err != nil {
		tty.Close()
		return nil, err
	}

	h := getTermHeight(tty) - 1

	return &Pager{
		lines:    lines,
		pin:      pin,
		height:   h,
		tty:      tty,
		oldState: oldState,
	}, nil
}

func getTermHeight(tty *os.File) int {
	// カーソル位置保存
	fmt.Fprint(tty, "\0337")

	// 右下に移動
	fmt.Fprint(tty, "\033[999;999H")

	// カーソル位置問い合わせ
	fmt.Fprint(tty, "\033[6n")

	var buf []byte
	tmp := make([]byte, 1)
	for {
		_, err := tty.Read(tmp)
		if err != nil {
			break
		}
		buf = append(buf, tmp[0])
		if tmp[0] == 'R' {
			break
		}
	}

	// カーソル位置復元
	fmt.Fprint(tty, "\0338")

	// 解析: ESC[row;colR
	var row int
	fmt.Sscanf(string(buf), "\033[%d;", &row)

	if row > 0 {
		return row
	}

	// 1. tmux / shell が用意する値（最優先）
	if v := os.Getenv("LINES"); v != "" {
		if h, err := strconv.Atoi(v); err == nil && h > 0 {
			return h
		}
	}

	f, err := os.Open("/proc/self/fd/1")
	if err == nil {
		defer f.Close()
		if h, _, err := term.GetSize(int(f.Fd())); err == nil && h > 0 {
			return h
		}
	}

	// 2. ioctl（信用できる時だけ）
	if term.IsTerminal(int(os.Stdout.Fd())) {
		if h, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && h > 0 {
			return h
		}
	}

	// 3. fallback
	return 24
}

func (p *Pager) Close() {
	term.Restore(int(p.tty.Fd()), p.oldState)
	p.clear()
	p.tty.Close()
}

func (p *Pager) Run() {
	maxScroll := len(p.lines) - (p.height - p.pin - 1)
	if maxScroll < 0 {
		maxScroll = 0
	}

	for {
		p.draw()

		buf := make([]byte, 1)
		p.tty.Read(buf)

		switch buf[0] {
		case 'q':
			return
		case 'j':
			if p.scroll < maxScroll {
				p.scroll++
			}
		case 'k':
			if p.scroll > 0 {
				p.scroll--
			}
		case ' ':
			p.scroll += p.height - p.pin - 1
			if p.scroll > maxScroll {
				p.scroll = maxScroll
			}
		case 'G':
			p.scroll = maxScroll
		case 'g':
			p.scroll = 0
		}
	}
}

func (p *Pager) draw() {
	p.clear()

	// 固定ヘッダ
	for i := 0; i < p.pin && i < len(p.lines); i++ {
		fmt.Fprint(p.tty, p.lines[i], "\r\n")
	}

	if p.pin > 0 {
		fmt.Fprint(p.tty, strings.Repeat("-", 40), "\r\n")
	}

	start := p.pin + p.scroll
	end := start + (p.height - p.pin - 1)

	if end > len(p.lines) {
		end = len(p.lines)
	}

	for i := start; i < end-1; i++ {
		fmt.Fprint(p.tty, p.lines[i], "\r\n")
	}
	if p.ruler {
		p.print(fmt.Sprintf("start=(%d,%d,%d), pin=%d, len=%d, scroll=%d, LINE=%s",
			start, end, p.height, p.pin, len(p.lines), p.scroll, os.Getenv("LINES")), "\r\n")
	}
}

func (p *Pager) print(a ...any) (int, error) {
	return fmt.Fprint(os.Stdout, a...)
}

func (p *Pager) clear() {
	// fmt.Fprint(os.Stdout, "\033[H\033[2J")
	fmt.Fprint(os.Stdout, "\033[0m") // リセット
}
