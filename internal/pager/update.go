package pager

import (
	"github.com/hiwane/pinview/internal/action"
)

// Update は action.Action に基づいて Model の状態を更新する。
// 画面描画や入出力は行わない。
func (m *Model) Update(a *action.Action) bool {

	maxOffset := m.maxScroll()

	switch a.Type {
	case action.ActQuit:
		return true
	case action.ActDown:
		m.OffsetY += a.Count
	case action.ActUp:
		m.OffsetY -= a.Count
	case action.ActPageDown:
		m.OffsetY += m.bodyHeight() * a.Count
	case action.ActPageUp:
		m.OffsetY -= m.bodyHeight() * a.Count
	case action.ActTop:
		if a.Count > 1 {
			m.OffsetY = a.Count - 1
		} else {
			m.OffsetY = 0
		}
	case action.ActBottom:
		if a.Count > 1 {
			m.OffsetY = maxOffset - (a.Count - 1)
		} else {
			m.OffsetY = maxOffset
		}
	case action.ActToggleRuler:
		m.Ruler = !m.Ruler
	case action.ActIncreaseHeader:
		m.SetHeader(m.header + a.Count)
	case action.ActDecreaseHeader:
		m.SetHeader(m.header - a.Count)
	case action.ActIncreaseFooter:
		m.SetFooter(m.footer + a.Count)
	case action.ActDecreaseFooter:
		m.SetFooter(m.footer - a.Count)
	case action.ActIncreaseHeight:
		m.Height += a.Count
	case action.ActDecreaseHeight:
		m.Height -= a.Count
	case action.ActScrollLeft:
		m.OffsetX -= a.Count
	case action.ActScrollRight:
		m.OffsetX += a.Count
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
