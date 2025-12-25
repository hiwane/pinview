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
	fmt.Printf("%s\n", lines[len(lines)-1])
}
