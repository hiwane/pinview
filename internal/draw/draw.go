package draw

import (
	"fmt"
)

func Draw(lines []string) {
	fmt.Print("\033[0m") // リセット
	for _, line := range lines {
		fmt.Print(line, "\r\n")
	}
}
