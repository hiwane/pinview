package pager

// Model は pager の状態を表す。
// この構造体は端末・stdin・stdout を一切知らない。
// 「何行を表示すべきか」を決めるための純粋な状態のみを持つ。
type Model struct {
	// Lines は入力全体（1行=1要素）
	Lines []string

	// Pin は常に先頭に固定表示される行数
	Pin int

	// Offset は本文部分の開始行インデックス
	// 常に Offset >= Pin になるよう制御される
	Offset int

	// Height は画面の高さ（表示可能な行数）
	Height int

	Width int

	Ruler bool

	key byte
}

// New は pager.Model を安全な初期値で生成する。
func NewModel(lines []string, pin, h int) *Model {

	// Offset は pin 行の直後から開始する
	return &Model{
		Lines:  lines,
		Pin:    pin,
		Height: h,
	}
}

func (m *Model) SetRuler(on bool) {
	m.Ruler = on
}
