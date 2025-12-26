package pager

import (
	"github.com/hiwane/pinview/internal/input"
)

type ModeType int

const (
	ModePager ModeType = iota
	ModeHelp
)

// Model は pager の状態を表す。
// この構造体は端末・stdin・stdout を一切知らない。
// 「何行を表示すべきか」を決めるための純粋な状態のみを持つ。
type Model struct {
	// Lines は入力全体（1行=1要素）
	Lines []string

	// 常に先頭に固定表示される行数
	header int
	footer int

	// Offset は本文部分の開始行インデックス
	// 常に Offset >= Pin になるよう制御される
	OffsetX int
	OffsetY int

	// Height は画面の高さ（表示可能な行数）
	Height int

	Width int

	Ruler bool

	key        input.Key
	sizeUpdate bool

	Mode ModeType
}

// New は pager.Model を安全な初期値で生成する。
func NewModel(lines []string) *Model {

	// Offset は pin 行の直後から開始する
	return &Model{
		Lines:  lines,
		header: 1,
		footer: 0,
		Height: 0,
		Mode:   ModePager,
	}
}

func (m *Model) SetRuler(on bool) {
	m.Ruler = on
}

func (m *Model) SetHeader(h int) {
	if h < 0 {
		h = 0
	}
	m.header = h
}

func (m *Model) SetFooter(f int) {
	if f < 0 {
		f = 0
	}
	m.footer = f
}

func (m *Model) bodyHeight() int {
	space := m.Height - m.header - m.footer
	if m.header > 0 { // separator
		space--
	}
	if m.footer > 0 { // separator
		space--
	}
	if space < 0 {
		space = 0
	}
	return space
}

func (m *Model) SetWidth(w int) {
	if w < 0 {
		w = 0
	}
	m.Width = w
}

func (m *Model) SetHeight(h int) {
	h -= 1
	if h < 0 {
		h = 0
	}
	m.Height = h
}

func (m *Model) SetSizeUpdate(f bool) {
	m.sizeUpdate = f
}

func (m *Model) SetKey(key input.Key) {
	m.key = key
}
