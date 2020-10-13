package display

import (
	"container/list"
	"fmt"

	"github.com/jbaramidze/term_collab_editor/logger"

	"github.com/gdamore/tcell"
	"github.com/jbaramidze/term_collab_editor/helper"
)

// Display ss
type Display struct {
	data               *list.List // list of LinkedListElement
	currentElement     *list.Element
	ScreenHandler      *ScreenHandler
	currentX, currentY int
	monitorChannel     chan ContentOperation
	blinkIsSet         bool
}

type LinkedListElement struct {
	data           []rune
	startingCoordY int
}

func Initialize() Display {
	d := InitScreen()

	ch1 := make(chan ContentOperation)

	cm := newContentManager(d, ch1)
	go cm.startLoop()

	InitAndStartBlinker(ch1)

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
				switch decoded.key {
				case tcell.KeyLeft:
					{
						if c.currentX > 0 {
							c.currentX--
						}
					}
				case tcell.KeyRight:
					{
						if len(c.currentElement.Value.([]rune)) > c.currentX {
							c.currentX++
						}
					}
				case tcell.KeyUp:
					{
						if c.currentY > 0 {
							c.currentY--
							c.currentX = helper.MinOf(
								c.currentX,
								len(c.currentElement.Value.([]rune)),
								len(c.currentElement.Prev().Value.([]rune)))
							c.currentElement = c.currentElement.Prev()
						}
					}
				case tcell.KeyDown:
					{
						if c.data.Len() > c.currentY+1 {
							c.currentY++
							c.currentX = helper.MinOf(
								c.currentX,
								len(c.currentElement.Value.([]rune)),
								len(c.currentElement.Next().Value.([]rune)))
							c.currentElement = c.currentElement.Next()
						}
					}
				case tcell.KeyEnter:
					{
						c.data.InsertAfter([]rune{}, c.currentElement)
						c.currentElement = c.currentElement.Next()
						c.currentX = 0
						c.currentY++
					}
				default:
					{
						c.newChar(decoded.rn)
					}
				}

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

func (c *Display) refreshBlinkStatus() {
	if c.blinkIsSet {
		c.setBlinkStatus()
	} else {
		c.clearBlinkStatus()
	}
}

func (c *Display) clearBlinkStatus() {
	if len(c.currentElement.Value.([]rune)) > c.currentX {
		c.ScreenHandler.putStr(c.currentX, c.currentY, rune(c.currentElement.Value.([]rune)[c.currentX]))
	} else {
		c.ScreenHandler.putStr(c.currentX, c.currentY, rune(' '))
	}
}
func (c *Display) setBlinkStatus() {
	c.ScreenHandler.putStr(c.currentX, c.currentY, rune('▉'))
}

func (c *Display) newChar(char rune) {
	logger.L.Log(fmt.Sprintf("Writing %v at %v:%v!", string(char), c.currentY, c.currentX))
	text := c.currentElement.Value.([]rune)
	text = append(append(text[:c.currentX], rune(char)), text[c.currentX:]...)
	c.currentElement.Value = text
	c.ScreenHandler.putStr(c.currentX, c.currentY, char)
	c.currentX++
	c.refreshBlinkStatus()
}

func (c *Display) resyncLine(line int) {

}

func newContentManager(s *ScreenHandler, c chan ContentOperation) Display {
	l := list.New()
	l.PushBack([]rune{})
	current := l.Back()
	cm := Display{ScreenHandler: s, data: l, currentElement: current, monitorChannel: c}
	return cm
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
