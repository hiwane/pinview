package input

import (
	"os"

	"golang.org/x/term"
)

type Input struct {
	fd   int
	old  *term.State
	file *os.File
}

func New(f *os.File) (*Input, error) {
	state, err := term.MakeRaw(int(f.Fd()))
	if err != nil {
		return nil, err
	}

	return &Input{
		fd:   int(f.Fd()),
		old:  state,
		file: f,
	}, nil
}

func (i *Input) ReadKey() (byte, error) {
	var b [1]byte
	_, err := i.file.Read(b[:])
	return b[0], err
}

func (i *Input) Close() {
	if i.old != nil {
		term.Restore(i.fd, i.old)
	}
}
