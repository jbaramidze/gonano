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
}

func (c *Display) dump() {
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

func (c *Display) setBlinker(b blinker) {
	c.blinker = b
}

func (c *Display) insert(char rune) {
	c.getCurrentEl().data = insertInSlice(c.getCurrentEl().data, char, c.getCurrentEl().pos)
	c.getCurrentEl().pos++

	c.resyncAll() // No need to call if nothing height-related changes
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
	} else {
		c.getCurrentEl().pos--
		c.getCurrentEl().data = removeFromSlice(c.getCurrentEl().data, c.getCurrentEl().pos)
	}

	c.resyncAll() // No need to call if nothing height-related changes
	c.syncCoords()
}

// Current line should have correct startingY !
func (c *Display) resyncAll() {
	curr := c.currentElement
	startingY := c.getCurrentEl().startingCoordY
	for curr != nil {
		line := curr.Value.(*Line)
		line.startingCoordY = startingY
		line.resync()
		startingY += line.height
		curr = curr.Next()
	}

	// Clean at startingY
	for i := 0; i < c.getWidth(); i++ {
		c.screen.putStr(i, startingY, rune(0))
	}
}

func (c *Display) pollKeyboard(resp chan bool) {
	for {
		ev := c.screen.pollKeyPress()
		if ev.k == tcell.KeyTAB {
			return
		}
		c.monitorChannel <- TypeOperation{rn: ev.rn, key: ev.k, resp: resp}
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
		}
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
