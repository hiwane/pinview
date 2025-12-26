package pager

import (
	"fmt"
	"github.com/hiwane/pinview/internal/input"
)

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
func (m *Model) Update(keyEvent input.Key) bool {

	key := keyEvent.Rune
	maxOffset := m.maxScroll()
	m.SetKey(keyEvent)
	fmt.Printf("Update key: %s, OffsetY: %d/%d\n", keyEvent, m.OffsetY, maxOffset)

	if keyEvent.Ctrl {
		// CTRL-
		switch key {
		case 'c':
			return true
		case 'd':
			m.OffsetY += m.bodyHeight()
		case 'u':
			m.OffsetY -= m.bodyHeight()
		}
	} else {
		switch key {
		case 'q': // 'q' または Ctrl-C
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
		case ' ':
			m.OffsetY += m.bodyHeight()
		case 'r':
			m.Ruler = !m.Ruler
		}
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
