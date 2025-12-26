package input

import (
	"os"

	"github.com/hiwane/pinview/internal/action"
	"github.com/mattn/go-tty"
)

type Input struct {
	tty        *tty.TTY
	Mode       InputMode
	CountBuf   int
	CommandBuf []rune
	InputBuf   []rune
}

type Key struct {
	Rune rune
	Ctrl bool
}

type InputMode int

const (
	ModeNormal InputMode = iota
	ModeCount
	ModeCommand
)

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
	if 1 <= k && k <= 26 {
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

func (i *Input) HandleKey(k Key) *action.Action {
	i.InputBuf = append(i.InputBuf, k.Rune)
	switch i.Mode {
	case ModeNormal:
		return i.handleNormalModeKey(k)
	case ModeCount:
		return i.handleCountModeKey(k)
	case ModeCommand:
		return i.handleCommandModeKey(k)
	}
	i.resetMode()
	return nil
}

func (i *Input) resetMode() {
	i.Mode = ModeNormal
	i.CountBuf = 0
	i.CommandBuf = nil
	i.InputBuf = nil
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func (i *Input) handleNormalModeKey(k Key) *action.Action {
	if isDigit(k.Rune) {
		i.Mode = ModeCount
		i.CountBuf = int(k.Rune - '0')
		return nil
	}
	// if k.Rune == ':' {
	// 	i.Mode = ModeCommand
	// 	i.CommandBuf = nil
	// 	return nil
	// }

	i.resetMode()
	return actionFromSingleKey(k, 1)
}

// 3g, 10j 等
func (i *Input) handleCountModeKey(k Key) *action.Action {
	if isDigit(k.Rune) {
		i.CountBuf = i.CountBuf*10 + int(k.Rune-'0')
		return nil
	}

	cnt := i.CountBuf
	if cnt == 0 {
		cnt = 1
	}
	i.resetMode()
	return actionFromSingleKey(k, cnt)
}

// コマンドモード. :で開始
// 使いみちがないので未実装
func (i *Input) handleCommandModeKey(k Key) *action.Action {
	switch k.Rune {
	// case '\n':
	// 	cmd := string(i.CommandBuf)
	// 	i.resetMode()
	// 	return actionFromCommand(cmd)
	case 0x1b: // ESC
		i.resetMode()
		return nil
	default:
		i.CommandBuf = append(i.CommandBuf, k.Rune)
		return nil
	}
}

func actionFromSingleKey(k Key, cnt int) *action.Action {
	if k.Ctrl {
		switch k.Rune {
		case 'd':
			return &action.Action{Type: action.ActPageDown, Count: cnt}
		case 'u':
			return &action.Action{Type: action.ActPageUp, Count: cnt}
		case 'c':
			return &action.Action{Type: action.ActQuit, Count: cnt}
		}
	} else {
		switch k.Rune {
		case 'j':
			return &action.Action{Type: action.ActDown, Count: cnt}
		case 'k':
			return &action.Action{Type: action.ActUp, Count: cnt}
		case 'd':
			return &action.Action{Type: action.ActPageDown, Count: cnt}
		case 'u':
			return &action.Action{Type: action.ActPageUp, Count: cnt}
		case 'g':
			return &action.Action{Type: action.ActTop, Count: cnt}
		case 'G':
			return &action.Action{Type: action.ActBottom, Count: cnt}
		case 'q':
			return &action.Action{Type: action.ActQuit, Count: cnt}
		case 'h':
			return &action.Action{Type: action.ActScrollLeft, Count: cnt}
		case 'l':
			return &action.Action{Type: action.ActScrollRight, Count: cnt}
		case 'T':
			return &action.Action{Type: action.ActIncreaseHeader, Count: cnt}
		case 't':
			return &action.Action{Type: action.ActDecreaseHeader, Count: cnt}
		case 'B':
			return &action.Action{Type: action.ActIncreaseFooter, Count: cnt}
		case 'b':
			return &action.Action{Type: action.ActDecreaseFooter, Count: cnt}
		case '+':
			return &action.Action{Type: action.ActIncreaseHeight, Count: cnt}
		case '-':
			return &action.Action{Type: action.ActDecreaseHeight, Count: cnt}
		}
	}
	return nil
}
