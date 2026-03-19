// Package ui contains general use case ui printing methods.
// Additionally includes input collectors using bufio.Reader with validation
package ui

import (
	"fmt"
	"strings"
)

// PrintTable assumes it accepts a balanced 2d string slice.
func PrintTable(table [][]string) {
	if len(table) < 1 {
		return
	}
	widths := make([]int, len(table[0]))
	for _, row := range table {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	separator := "+"
	for _, width := range widths {
		separator += strings.Repeat("-", width+2) + "+"
	}
	fmt.Printf("%s\n", separator)
	for _, row := range table {
		fmt.Printf("|")
		for i, cell := range row {
			padding := widths[i] - len(cell)
			fmt.Printf(" %s%s |", cell, strings.Repeat(" ", padding))
		}
		fmt.Printf("\n%s\n", separator)
	}
}
