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

	/////////////////////////////
	// ヘッダー部
	/////////////////////////////
	if m.header > 0 {
		for i := 0; i < m.header; i++ {
			out = append(out, m.Lines[i])
		}
		out = append(out, m.Separator())
	}

	/////////////////////////////
	// 本文
	/////////////////////////////
	start := m.header + m.Offset
	end := start + m.bodyHeight()
	if m.Ruler {
		end -= 1
	}
	if end > len(m.Lines) {
		end = len(m.Lines)
	}

	for i := start; i < end; i++ {
		out = append(out, m.Lines[i])
	}

	/////////////////////////////
	// フッター部
	/////////////////////////////
	if m.footer > 0 {
		out = append(out, m.Separator())
		for i := len(m.Lines) - m.footer; i < len(m.Lines); i++ {
			out = append(out, m.Lines[i])
		}
	}

	/////////////////////////////
	// ruler
	/////////////////////////////
	if m.Ruler {
		out = append(out, fmt.Sprintf("    start=(%d,%d,%d), pin=(%d,%d) len=%d, scroll=%d, key=0x%x, sig=%v",
			start, end, m.Height, m.header, m.footer, len(m.Lines), m.Offset, m.key, m.sizeUpdate))
	}

	return out
}
