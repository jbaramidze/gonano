package main

import (
	"container/list"
	"log"

	"github.com/gdamore/tcell"
)

// Display ss
type Display struct {
	data               *list.List // list of *Line
	currentElement     *list.Element
	screen             screenHandler
	currentX, currentY int
	monitorChannel     chan ContentOperation
	blinker            blinker
	offsetY            int
}

func (c *Display) dump() {
	log.Printf("Current: x:%v, y:%v", c.currentX, c.currentY)
	log.Println("Dumping lines:")
	for i, e := 0, c.data.Front(); e != nil; i, e = i+1, e.Next() {
		l := e.Value.(*Line)
		log.Printf("Line %v: data %v startY %v height %v pos %v", i, string(l.data), l.startingCoordY, l.height, l.pos)
	}
}

func (c *Display) getWidth() int {
	w, _ := c.screen.getSize()
	return w
}

func (c *Display) getHeight() int {
	_, h := c.screen.getSize()
	return h
}

func (c *Display) getCurrentEl() *Line {
	return c.currentElement.Value.(*Line)
}

func (c *Display) hasPrevEl() bool {
	return c.currentElement.Prev() != nil
}

func (c *Display) getPrevEl() *Line {
	return c.currentElement.Prev().Value.(*Line)
}

func (c *Display) hasNextEl() bool {
	return c.currentElement.Next() != nil
}

func (c *Display) getNextEl() *Line {
	return c.currentElement.Next().Value.(*Line)
}

func createDisplay(handler screenHandler) *Display {
	d := initializeDisplay(handler)

	return d
}

func (c *Display) insert(char rune) {
	oldCursorY := c.getCurrentEl().getRelativeCursorY()
	c.getCurrentEl().data = insertInSlice(c.getCurrentEl().data, char, c.getCurrentEl().pos)
	c.getCurrentEl().pos++
	newCursorY := c.getCurrentEl().getRelativeCursorY()

	if oldCursorY != newCursorY {
		onScreenCursorY := c.getCurrentEl().getOnScreenCursorY()
		if onScreenCursorY >= c.getHeight() {
			c.offsetY++
			c.resyncBelow(c.data.Front())
		} else {
			c.resyncBelow(c.currentElement)
		}
	} else {
		c.getCurrentEl().resync()
	}
	c.syncCoords()
}

func (c *Display) remove() {
	if c.getCurrentEl().pos == 0 {
		if !c.hasPrevEl() {
			return
		}

		p := c.currentElement.Prev()
		pl := p.Value.(*Line)
		pl.data = append(pl.data, c.getCurrentEl().data...)
		c.data.Remove(c.currentElement)
		c.currentElement = p

		c.resyncBelow(c.currentElement)
	} else {
		oldHeight := c.getCurrentEl().calculateHeight()
		c.getCurrentEl().pos--
		c.getCurrentEl().data = removeFromSlice(c.getCurrentEl().data, c.getCurrentEl().pos)
		newHeight := c.getCurrentEl().calculateHeight()

		if newHeight != oldHeight {
			c.resyncBelow(c.currentElement)
		} else {
			c.getCurrentEl().resync()
		}
	}

	c.syncCoords()
}

func (c *Display) recalcBelow(from *list.Element) {
	startingY := from.Value.(*Line).startingCoordY
	for ; from != nil; from = from.Next() {
		line := from.Value.(*Line)
		line.startingCoordY = startingY
		line.height = line.calculateHeight()
		startingY += line.height
	}
}

// Current line should have correct startingY !
func (c *Display) resyncBelow(from *list.Element) {
	c.recalcBelow(from)
	for ; from != nil && from.Value.(*Line).startingCoordY-c.offsetY < c.getHeight(); from = from.Next() {
		from.Value.(*Line).resync()
	}

	// Clean at startingY
	startingY := 0
	if from != nil {
		startingY = from.Value.(*Line).startingCoordY + from.Value.(*Line).height
	} else {
		startingY = c.data.Back().Value.(*Line).startingCoordY + c.data.Back().Value.(*Line).height
	}

	for ; startingY-c.offsetY < c.getHeight(); startingY++ {
		for i := 0; i < c.getWidth(); i++ {
			c.screen.clearStr(i, startingY-c.offsetY)
		}
	}
}

// Close display
func (c *Display) Close() {
	c.screen.close()
}

func (c *Display) startLoop() {
	for {
		op := <-c.monitorChannel
		c.blinker.clear()
		switch decoded := op.(type) {
		case TypeOperation:
			{
				c.handleKeyPress(decoded)
				c.blinker.refresh()
				if decoded.resp != nil {
					decoded.resp <- true
				}
			}
		case BlinkOperation:
			{
				c.blinker.refresh()
			}
		case AnnouncementOperation:
			{
				c.drawText(decoded.text)
				if decoded.resp != nil {
					decoded.resp <- true
				}
			}
		}
	}
}

/*
 *****************************
 * Text will be in something *
 *        like this          *
 *****************************
 */

func (c *Display) drawText(text []string) {
	w, h := c.screen.getSize()

	maxW := 0
	for j := 0; j < len(text); j++ {
		maxW = maxOf(maxW, len(text[j]))
	}

	startX := w/2 - maxW/2 - 2
	endX := w/2 - maxW/2 - 2 + maxW + 4

	for i := startX; i < endX; i++ {
		c.screen.putStr(i, h/2-1, rune('*'))
	}
	for j := 0; j < len(text); j++ {
		internalStartX := w/2 - len(text[j])/2
		internalEndX := w/2 - len(text[j])/2 + len(text[j])
		// first *
		c.screen.putStr(startX, h/2+j, rune('*'))
		// spaces
		for i := startX + 1; i < internalStartX; i++ {
			c.screen.clearStr(i, h/2+j)
		}
		// text
		for i, ch := range text[j] {
			c.screen.putStr(internalStartX+i, h/2+j, ch)
		}
		// spaces after text
		for i := internalEndX; i < endX-1; i++ {
			c.screen.clearStr(i, h/2+j)
		}
		// last *
		c.screen.putStr(endX-1, h/2+j, rune('*'))
	}

	for i := startX; i < endX; i++ {
		c.screen.putStr(i, h/2+len(text), rune('*'))
	}
}

func initializeDisplay(s screenHandler) *Display {
	channel := make(chan ContentOperation)
	lst := list.New()
	d := Display{screen: s, data: lst, monitorChannel: channel}
	lst.PushBack(d.newLine())
	d.currentElement = lst.Back()
	return &d
}

// ContentOperation s
type ContentOperation interface{}

// TypeOperation ss
type TypeOperation struct {
	rn   rune
	key  tcell.Key
	resp chan bool
}

// BlinkOperation ss
type BlinkOperation struct {
	blink bool
}

// AnnouncementOperation ss
type AnnouncementOperation struct {
	text []string
	resp chan bool
}
