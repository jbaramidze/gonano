package display

import (
	"container/list"

	"github.com/gdamore/tcell"
)

// Display ss
type Display struct {
	data               *list.List // list of *Line
	currentElement     *list.Element
	ScreenHandler      *ScreenHandler
	currentX, currentY int
	monitorChannel     chan ContentOperation
	blinkIsSet         bool
}

func (c *Display) getWidth() int {
	w, _ := c.ScreenHandler.screen.Size()
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

// Initialize display
func Initialize() *Display {
	handler, channel := InitScreenHandler()

	cm := newContentManager(handler, channel)
	go cm.startLoop()

	initAndStartBlinker(channel)
	return cm
}

// Poll display
func (c *Display) Poll() {
	for {
		switch ev := c.ScreenHandler.screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyTAB {
				return
			}
			c.monitorChannel <- TypeOperation{rn: ev.Rune(), key: ev.Key()}
		}
	}
}

// Close display
func (c *Display) Close() {
	c.ScreenHandler.screen.Fini()
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

func newContentManager(s *ScreenHandler, c chan ContentOperation) *Display {
	lst := list.New()
	d := Display{ScreenHandler: s, data: lst, monitorChannel: c}
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
