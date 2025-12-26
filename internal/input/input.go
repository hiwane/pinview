package input

import (
	"os"

	"github.com/mattn/go-tty"
)

type Input struct {
	tty *tty.TTY
}

type Key struct {
	Rune rune
	Ctrl bool
}

func New(f *os.File) (*Input, error) {
	t, err := tty.Open()
	if err != nil {
		return nil, err
	}

	return &Input{
		tty: t,
	}, nil
}

func (i *Input) ReadRune() (Key, error) {
	k, err := i.tty.ReadRune()
	if err != nil {
		return Key{}, err
	}
	if 0 <= k && k <= 0x1f {
		return Key{Rune: rune('a' + k - 1), Ctrl: true}, nil
	} else {
		return Key{Rune: k, Ctrl: false}, nil
	}
}

func (i *Input) Close() error {
	return i.tty.Close()
}

func (i *Input) Runes() <-chan Key {
	ch := make(chan Key)

	go func() {
		defer close(ch)
		for {
			r, err := i.ReadRune()
			if err != nil {
				return
			}
			ch <- r
		}
	}()

	return ch
}

func (key Key) String() string {
	if key.Ctrl {
		return "C-" + string(key.Rune)
	}
	return string(key.Rune)
}
