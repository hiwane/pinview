package pager

// Action は pager に対する操作を表す。
// 入力方法（キー・マウス等）とは切り離して定義する。
type Action int

const (
	ActNone Action = iota
	ActDown
	ActUp
	ActQuit
)

// Update は Action に基づいて Model の状態を更新する。
// 画面描画や入出力は行わない。
func (m *Model) Update(key rune) bool {

	maxOffset := m.maxScroll()
	m.key = key

	switch key {
	case 'q', 3: // 'q' または Ctrl-C
		return true
	case 'j':
		m.Offset++
	case 'k':
		m.Offset--
	case 'h':
		m.SetHeader(m.header - 1)
	case 'H':
		m.SetHeader(m.header + 1)
	case 'f':
		m.SetFooter(m.footer - 1)
	case 'F':
		m.SetFooter(m.footer + 1)
	case '+':
		m.Height++
	case '-':
		m.Height--
	case 'g':
		m.Offset = 0
	case 'G':
		m.Offset = maxOffset
	case ' ', 0x4:
		m.Offset += m.bodyHeight()
	case 0x15:
		m.Offset -= m.bodyHeight()
	case 'r':
		m.Ruler = !m.Ruler
	}

	if m.Offset < 0 {
		m.Offset = 0
	} else if m.Offset > maxOffset {
		m.Offset = maxOffset
	}

	return false
}

func (m *Model) maxScroll() int {
	maxScroll := len(m.Lines) - (m.Height - m.header - m.footer - 1)
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}
