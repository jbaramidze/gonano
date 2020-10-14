package display

type Line struct {
	data           []rune
	startingCoordY int
	height         int
	currentY       int
	display        *Display
}

func (c *Display) newLine() *Line {
	return &Line{
		data:           []rune{},
		startingCoordY: 0,
		height:         1,
		currentY:       0,
		display:        c,
	}
}

func (l *Line) hasCharsOnTheRight() bool {
	return len(l.data) > l.display.currentX
}
