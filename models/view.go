package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const spacebar = " "

type view struct {
	Width   int
	Height  int
	move    bool
	YOffset int

	YPosition       int
	initialized     bool
	Style           lipgloss.Style
	MouseWheelDelta int
	lines           []string
}

func (v *view) Update(lines []string) {
	v.lines = lines
}

func NewView(width, height int) (v view) {
	v.Width = width
	v.Height = height
	return v
}

func (v *view) setInitialValues() {
	v.initialized = true
}

func (v view) AtTop() bool {
	return v.YOffset <= 0
}

func (v view) AtBottom() bool {
	return v.YOffset >= v.maxYOffset()
}

func (v view) PastBottom() bool {
	return v.YOffset > v.maxYOffset()
}

func (v view) maxYOffset() int {
	return max(0, len(v.lines)-v.Height)
}

func (v view) visibleLines() (lines []string) {
	if len(v.lines) > 0 && v.move {
		top := max(0, v.YOffset)
		bottom := clamp(v.YOffset+v.Height, top, len(v.lines))
		lines = v.lines[top:bottom]
	} else if len(v.lines) >= v.Height && !v.move {
		lines = v.lines[len(v.lines)-v.Height:]
	} else if len(v.lines) > 0 {
		top := max(0, v.YOffset)
		bottom := clamp(v.YOffset+v.Height, top, len(v.lines))
		lines = v.lines[top:bottom]
	}
	return lines
}

func (v *view) ClearLine() {
	v.lines = []string{}
}

func (v view) printVisbleLines() {
	lines := v.visibleLines()
	add := 0
	if len(lines) < v.Height {
		add = v.Height - len(lines)
	}
	for _, line := range lines {
		fmt.Print("\r")
		fmt.Print(strings.Repeat(" ", v.Width))
		fmt.Print("\r")
		fmt.Print(line)
	}
	for add > 0 {
		fmt.Println("")
		add--
	}
	send := fmt.Sprintf("\033[%dA", v.Height)
	fmt.Print(send)
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}

func (v *view) ViewDown() []string {
	if v.AtBottom() {
		return nil
	}

	return v.LineDown(v.Height)
}

func (v *view) SetYOffset(n int) {
	v.YOffset = clamp(n, 0, v.maxYOffset())
}

func (v *view) LineDown(n int) (lines []string) {
	v.move = true
	if v.AtBottom() || n == 0 || len(v.lines) == 0 {
		v.move = false
		return nil
	}
	v.SetYOffset(v.YOffset + n)
	bottom := clamp(v.YOffset+v.Height, 0, len(v.lines))
	top := clamp(v.YOffset+v.Height-n, 0, bottom)
	return v.lines[top:bottom]
}

func (v *view) LineUp(n int) (lines []string) {
	v.move = true
	if v.AtTop() || n == 0 || len(v.lines) == 0 {
		return nil
	}

	v.SetYOffset(v.YOffset - n)
	top := max(0, v.YOffset)
	bottom := clamp(v.YOffset+n, 0, v.maxYOffset())
	return v.lines[top:bottom]
}

func (v *view) ViewUp() []string {
	if v.AtTop() {
		return nil
	}

	return v.LineUp(v.Height)
}
