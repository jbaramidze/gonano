package main

import (
	"log"

	"github.com/gdamore/tcell"
)

func (c *Display) syncCoords() {
	blinkerX, blinkerY := c.getCurrentEl().getBlinkerCoords()
	c.currentX = blinkerX
	c.currentY = blinkerY + c.getCurrentEl().startingCoordY
}

func (c *Display) handleKeyPress(op TypeOperation) {
	switch op.key {
	case tcell.KeyLeft:
		{
			c.getCurrentEl().moveLeft()
			c.syncCoords()
		}
	case tcell.KeyRight:
		{
			c.getCurrentEl().moveRight()
			c.syncCoords()
		}
	case tcell.KeyUp:
		{
			if c.hasPrevEl() {
				pos := c.getCurrentEl().pos
				c.currentElement = c.currentElement.Prev()
				c.getCurrentEl().pos = minOf(len(c.getCurrentEl().data), pos)
				c.syncCoords()
			}
		}
	case tcell.KeyDown:
		{
			if c.hasNextEl() {
				pos := c.getCurrentEl().pos
				c.currentElement = c.currentElement.Next()
				c.getCurrentEl().pos = minOf(len(c.getCurrentEl().data), pos)
				c.syncCoords()
			}
		}
	case tcell.KeyEnter:
		{
			cur := c.getCurrentEl()
			newData := make([]rune, len(cur.data)-cur.pos)
			copy(newData, cur.data[cur.pos:])
			newItem := Line{data: newData, startingCoordY: cur.startingCoordY + cur.getCurrentY(), height: -1, pos: 0, display: c}
			c.data.InsertAfter(&newItem, c.currentElement)
			cur.data = cur.data[:cur.pos] // we can optimize memory here, by duplicating it.
			c.resyncAll()
			c.currentElement = c.currentElement.Next()
			c.syncCoords()
		}
	case tcell.KeyCtrlF:
		{
			c.dump()
		}
	default:
		{
			log.Print("Key pressed", op.rn)
			c.getCurrentEl().insertCharInCurrentPosition(op.rn)
		}
	}
}
