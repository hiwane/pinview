package term

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func WatchResize(ctx context.Context, tty *os.File) <-chan Size {
	ch := make(chan Size)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)

	go func() {
		defer close(ch)
		defer signal.Stop(sigCh)

		var last Size

		// 初期サイズ
		size, err := GetSize(tty)
		if err == nil {
			last = size
			ch <- size
		} else {
			fmt.Println("Failed to get terminal size:", err)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-sigCh:
				size, err := GetSize(tty)
				if err != nil {
					continue
				}
				if size != last {
					last = size
					ch <- size
				}
			}
		}
	}()

	return ch
}
