package term

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
)

// IsInteractive は stdout が端末かどうかを返す。
// パイプやリダイレクト時は false になる。
func IsInteractive() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// GetHeight は現在の端末の高さを取得する。
// 取得に失敗した場合でも必ず正の値を返す。
func GetHeight(tty *os.File) int {
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

func ViewClearScreen(tty *os.File) {
	fmt.Fprint(tty, "\033[2J") // 画面全体をクリア
	fmt.Fprint(tty, "\033[0m") // 属性リセット
	fmt.Fprint(tty, "\033[H")  // カーソルを先頭に移動
}
