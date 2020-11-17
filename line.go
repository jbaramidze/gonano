package main

// Line structure
type Line struct {
	data           []rune
	startingCoordY int
	height         int
	display        *Display
	pos            int
}

func (c *Display) newLine() *Line {
	return &Line{
		data:           []rune{},
		startingCoordY: 0,
		pos:            0,
		height:         1,
		display:        c,
	}
}

func (l *Line) moveLeft() {
	if l.pos > 0 {
		l.pos--
	}
}

func (l *Line) moveRight() {
	if l.pos < len(l.data) {
		l.pos++
	}
}

func (l *Line) getCurrentChar() rune {
	return l.data[l.pos]
}

func (l *Line) getRelativeBlinkerCoordsByPos() (int, int) {
	y := l.pos / l.display.getWidth()
	x := l.pos - y*l.display.getWidth()
	return x, y
}

func (l *Line) getSmallestOffsetToFitLineOnDisplay() int {
	h := l.calculateHeight()
	// last chars are on last display line: lineH + startingY - offsetY - 1 == displayH - 1
	offset := h + l.startingCoordY - l.display.getHeight()
	if offset < 0 {
		return 0
	}

	return offset
}

func (l *Line) getAbsoluteStartingY() int {
	return l.startingCoordY - l.display.offsetY
}
func (l *Line) getAbsoluteEndingY() int {
	return l.startingCoordY - l.display.offsetY + l.calculateHeight() - 1
}

func (l *Line) getCurrentY() int {
	return l.pos / l.display.getWidth()
}

func (l *Line) calculateHeight() int {
	return 1 + (len(l.data)-1)/l.display.getWidth()
}

func (l *Line) resync() {
	usableWidth := l.display.getWidth()
	var line int
	for i, r := range l.data {
		line = i / usableWidth
		y := l.startingCoordY + line - l.display.offsetY
		if y >= 0 {
			l.display.screen.putStr(i-(line*usableWidth), y, r)
		}
	}
	l.height = line + 1
	// Clear the rest
	idx := 0
	for i := 0; i < l.height; i++ {
		for j := 0; j < usableWidth; j++ {
			idx++
			if idx <= len(l.data) {
				continue
			}
			y := l.startingCoordY + i - l.display.offsetY
			if y >= 0 {
				l.display.screen.putStr(j, y, 0)
			}
		}
	}
}
