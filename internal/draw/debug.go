package draw

import (
	"fmt"
)

type DebugDrawer struct {
}

func NewDebugDrawer() *DebugDrawer {
	return &DebugDrawer{}
}

func (d *DebugDrawer) Render(lines []string) {
	if len(lines) > 0 {
		fmt.Printf("%s\n", lines[len(lines)-1])
	} else {
		fmt.Printf("empty line\n")
	}
}
