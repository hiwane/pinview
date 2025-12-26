package pager

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

func (m *Model) Separator() string {
	return strings.Repeat("-", 40)
}

func sliceLineByWidth(line string, startCol, width int) string {
	if width <= 0 {
		return ""
	}

	col := 0
	var out []rune

	for _, r := range line {
		w := runewidth.RuneWidth(r)
		if w == 0 {
			continue
		}

		// 表示開始していない
		if col+w <= startCol {
			col += w
			continue
		}

		// 表示領域を超えた
		if col >= startCol+width || col+w > startCol+width {
			break
		}

		out = append(out, r)
		col += w
	}

	return string(out)
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
			out = append(out, sliceLineByWidth(m.Lines[i], m.OffsetX, m.Width))
		}
		out = append(out, m.Separator())
	}

	/////////////////////////////
	// 本文
	/////////////////////////////
	start := m.header + m.OffsetY
	end := start + m.bodyHeight()
	if m.Ruler {
		end -= 1
	}
	if end > len(m.Lines) {
		end = len(m.Lines)
	}

	for i := start; i < end; i++ {
		out = append(out, sliceLineByWidth(m.Lines[i], m.OffsetX, m.Width))
	}

	/////////////////////////////
	// フッター部
	/////////////////////////////
	if m.footer > 0 {
		out = append(out, m.Separator())
		for i := len(m.Lines) - m.footer; i < len(m.Lines); i++ {
			out = append(out, sliceLineByWidth(m.Lines[i], m.OffsetX, m.Width))
		}
	}

	/////////////////////////////
	// ruler
	/////////////////////////////
	if m.Ruler {
		out = append(out, fmt.Sprintf("    start=(%d,%d,%d), pin=(%d,%d) len=%d, offfset=(%d,%d), key=0x%x, sig=%v",
			start, end, m.Height, m.header, m.footer, len(m.Lines), m.OffsetX, m.OffsetY, m.key, m.sizeUpdate))
	}

	return out
}
