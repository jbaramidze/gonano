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
	statusBar          statusBar
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
	channel := make(chan ContentOperation)
	lst := list.New()
	d := Display{screen: handler, data: lst, monitorChannel: channel}
	lst.PushBack(d.newLine())
	d.currentElement = lst.Back()
	return &d
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
				c.statusBar.draw(decoded.text)
				if decoded.resp != nil {
					decoded.resp <- true
				}
			}
		}
	}
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
