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

func (l *Line) getBlinkerCoords() (int, int) {
	y := l.pos / l.display.getWidth()
	x := l.pos - y*l.display.getWidth()
	return x, y
}

func (l *Line) insertCharInCurrentPosition(char rune) {
	l.data = insertInSlice(l.data, char, l.pos)
	l.pos++

	l.resyncCurrentLine()
	l.display.syncCoords()
}

func (l *Line) resyncCurrentLine() {
	usableWidth := l.display.getWidth()
	var line int
	for i, r := range l.data {
		line = i / usableWidth
		l.display.screen.putStr(i-(line*usableWidth), l.startingCoordY+line, r)
	}
	l.height = line + 1
}
