package draw

import (
	"fmt"
	"os"
)

type TerminalDrawer struct {
	fp *os.File
}

func NewTerminalDrawer() *TerminalDrawer {
	return &TerminalDrawer{
		fp: os.Stdout,
	}
}

func (d *TerminalDrawer) viewClearScreen() {
	fmt.Fprint(d.fp, "\033[2J") // 画面全体をクリア
	fmt.Fprint(d.fp, "\033[0m") // 属性リセット
	fmt.Fprint(d.fp, "\033[H")  // カーソルを先頭に移動
}

func (d *TerminalDrawer) Render(lines []string) {
	d.viewClearScreen()
	fmt.Fprint(d.fp, "\033[0m") // リセット
	for _, line := range lines {
		fmt.Fprint(d.fp, line, "\r\n")
	}
}
