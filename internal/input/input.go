package input

import (
	"os"

	"github.com/mattn/go-tty"
)

type Input struct {
	tty *tty.TTY
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

func (i *Input) ReadRune() (rune, error) {
	return i.tty.ReadRune()
}

func (i *Input) Close() error {
	return i.tty.Close()
}

func (i *Input) Runes() <-chan rune {
	ch := make(chan rune)

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
