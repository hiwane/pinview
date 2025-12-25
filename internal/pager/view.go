package pager

import (
	"fmt"
	"strings"
)

func (m *Model) Separator() string {
	return strings.Repeat("-", 40)
}

// VisibleLines は現在の Model 状態に基づいて
// 「画面に表示すべき行」を返す。
// 戻り値の各要素は改行を含まない論理行。
func (m *Model) View() []string {
	if len(m.Lines) == 0 || m.Height <= 0 {
		return nil
	}

	out := make([]string, 0, m.Height)

	// ヘッダ部
	for i := 0; i < m.Pin; i++ {
		out = append(out, m.Lines[i])
	}

	if m.Pin > 0 {
		out = append(out, m.Separator())
	}

	start := m.Pin + m.Offset
	end := start + (m.Height - m.Pin - 1)
	if m.Ruler {
		end -= 1
	}
	if end > len(m.Lines) {
		end = len(m.Lines)
	}

	for i := start; i < end; i++ {
		out = append(out, m.Lines[i])
	}

	if m.Ruler {
		out = append(out, fmt.Sprintf("    start=(%d,%d,%d), pin=%d, len=%d, scroll=%d, key=0x%x",
			start, end, m.Height, m.Pin, len(m.Lines), m.Offset, m.key))
	}

	return out
}
