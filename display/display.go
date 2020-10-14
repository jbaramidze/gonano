package display

import (
	"container/list"
	"fmt"

	"github.com/jbaramidze/term_collab_editor/helper"

	"github.com/jbaramidze/term_collab_editor/logger"

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

func (c *Display) getPrevEl() *Line {
	return c.currentElement.Prev().Value.(*Line)
}

func (c *Display) getNextEl() *Line {
	return c.currentElement.Next().Value.(*Line)
}

func Initialize() Display {
	handler, channel := InitScreenHandler()

	cm := newContentManager(handler, channel)
	go cm.startLoop()

	initAndStartBlinker(channel)
	return cm
}

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

func (c *Display) newChar(char rune) {
	logger.L.Log(fmt.Sprintf("Writing %v at %v:%v!", string(char), c.currentY, c.currentX))
	c.getCurrentEl().data = helper.InsertInSlice(c.getCurrentEl().data, char, c.currentX)
	logger.L.Log("Current string " + string(c.getCurrentEl().data))
	c.resyncCurrentLine()
	c.currentX++
	if c.currentX == c.getWidth()-1 {
		c.currentX = 0
		c.currentY++
	}
}

func (c *Display) resyncCurrentLine() {
	usableWidth := c.getWidth() - 1
	var line int
	for i, r := range c.getCurrentEl().data {
		line = i / usableWidth
		c.ScreenHandler.putStr(i-(line*usableWidth), c.getCurrentEl().startingCoordY+line, r)
	}
	c.getCurrentEl().height = line + 1
}

func newContentManager(s *ScreenHandler, c chan ContentOperation) Display {
	lst := list.New()
	d := Display{ScreenHandler: s, data: lst, monitorChannel: c}
	lst.PushBack(d.newLine())
	d.currentElement = lst.Back()
	return d
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
