package term

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Size struct {
	Width  int
	Height int
}


// IsInteractive は stdout が端末かどうかを返す。
// パイプやリダイレクト時は false になる。
func IsInteractive() bool {
	return IsTTY(os.Stdout)
}

func IsTTY(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

// GetHeight は現在の端末の高さを取得する。
// 取得に失敗した場合でも必ず正の値を返す。
func GetSize(tty *os.File) (Size, error) {
	w, h, err := term.GetSize(int(tty.Fd()))
	return Size{Width: w, Height: h}, err
}

func ViewClearScreen(tty *os.File) {
	fmt.Fprint(tty, "\033[2J") // 画面全体をクリア
	fmt.Fprint(tty, "\033[0m") // 属性リセット
	fmt.Fprint(tty, "\033[H")  // カーソルを先頭に移動
}
