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
	blinkIsSet         bool
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
	channel := make(chan ContentOperation)
	d := initializeDisplay(handler, channel)
	go d.startLoop()

	initAndStartBlinker(channel)
	return d
}

func (c *Display) poll() {
	for {
		ev := c.screen.pollKeyPress()
		if ev.k == tcell.KeyTAB {
			return
		}
		c.monitorChannel <- TypeOperation{rn: ev.rn, key: ev.k}
	}
}

// Close display
func (c *Display) Close() {
	c.screen.close()
}

func (c *Display) startLoop() {
	for {
		op := <-c.monitorChannel
		c.clearBlinkStatus()
		switch decoded := op.(type) {
		case TypeOperation:
			{
				c.handleKeyPress(decoded)
				c.refreshBlinkStatus()
			}
		case BlinkOperation:
			{
				c.blinkIsSet = decoded.blink
				c.refreshBlinkStatus()
			}
		}
	}
}

func initializeDisplay(s screenHandler, c chan ContentOperation) *Display {
	lst := list.New()
	d := Display{screen: s, data: lst, monitorChannel: c}
	lst.PushBack(d.newLine())
	d.currentElement = lst.Back()
	return &d
}

// ContentOperation s
type ContentOperation interface{}

// TypeOperation ss
type TypeOperation struct {
	rn  rune
	key tcell.Key
}

// BlinkOperation ss
type BlinkOperation struct {
	blink bool
}
