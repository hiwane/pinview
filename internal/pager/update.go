package pager

import "fmt"

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
	m.SetKey(key)
	fmt.Printf("Update key: %q, OffsetY: %d/%d\n", key, m.OffsetY, maxOffset)

	switch key {
	case 'q', 3: // 'q' または Ctrl-C
		return true
	case 'h':
		m.OffsetX--
	case 'j':
		m.OffsetY++
	case 'k':
		m.OffsetY--
	case 'l':
		m.OffsetX++
	case 't':
		m.SetHeader(m.header - 1)
	case 'T': // top
		m.SetHeader(m.header + 1)
	case 'b':
		m.SetFooter(m.footer - 1)
	case 'B': // bottom
		m.SetFooter(m.footer + 1)
	case '+':
		m.Height++
	case '-':
		m.Height--
	case 'g':
		m.OffsetY = 0
	case 'G':
		m.OffsetY = maxOffset
	case ' ', 0x4:
		m.OffsetY += m.bodyHeight()
	case 0x15:
		m.OffsetY -= m.bodyHeight()
	case 'r':
		m.Ruler = !m.Ruler
	}

	if m.OffsetY < 0 {
		m.OffsetY = 0
	} else if m.OffsetY > maxOffset {
		m.OffsetY = maxOffset
	}
	if m.OffsetX < 0 {
		m.OffsetX = 0
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
