package main

import (
	"log"

	"github.com/gdamore/tcell"
)

func (c *Display) getBlinkerX() int {
	blinkerX, _ := c.getCurrentEl().getRelativeBlinkerCoordsByPos()
	return blinkerX
}

func (c *Display) getBlinkerY() int {
	_, blinkerY := c.getCurrentEl().getRelativeBlinkerCoordsByPos()
	return blinkerY + c.getCurrentEl().startingCoordY - c.offsetY
}

func (c *Display) handleKeyPress(op typeOperation) {
	switch op.key {
	case tcell.KeyLeft:
		{
			oldCursorY := c.getCurrentEl().getRelativeCursorY()
			c.getCurrentEl().moveLeft()
			newCursorY := c.getCurrentEl().getRelativeCursorY()

			if oldCursorY != newCursorY {
				c.resyncNewCursorY()
			}
		}
	case tcell.KeyRight:
		{
			oldCursorY := c.getCurrentEl().getRelativeCursorY()
			c.getCurrentEl().moveRight()
			newCursorY := c.getCurrentEl().getRelativeCursorY()

			if oldCursorY != newCursorY {
				c.resyncNewCursorY()
			}
		}
	case tcell.KeyUp:
		{
			if c.hasPrevEl() {
				pos := c.getCurrentEl().pos
				c.currentElement = c.currentElement.Prev()
				c.getCurrentEl().pos = minOf(len(c.getCurrentEl().data), pos)
				if c.getCurrentEl().getOnScreenLineStartingY() < 0 {
					c.offsetY = c.getCurrentEl().startingCoordY
					c.resyncBelow(c.data.Front())
				}
			}
		}
	case tcell.KeyDown:
		{
			if c.hasNextEl() {
				pos := c.getCurrentEl().pos
				c.currentElement = c.currentElement.Next()
				c.getCurrentEl().pos = minOf(len(c.getCurrentEl().data), pos)
				if c.getCurrentEl().getOnScreenLineEndingY() >= c.getHeight() {
					// Try to fit next line.
					c.getCurrentEl().makeSmallestOffsetToFitLineOnDisplay()
					c.resyncBelow(c.data.Front())
				}
			}
		}
	case tcell.KeyEnter:
		{
			cur := c.getCurrentEl()

			// Create new line
			newData := make([]rune, len(cur.data)-cur.pos)
			copy(newData, cur.data[cur.pos:])
			newItem := Line{data: newData, startingCoordY: cur.startingCoordY + cur.getRelativeCharBeforeCursorY() + 1, pos: 0, display: c}

			c.data.InsertAfter(&newItem, c.currentElement)

			cur.data = cur.data[:cur.pos] // we can optimize memory here, by duplicating it.
			c.currentElement = c.currentElement.Next()

			if c.getCurrentEl().getOnScreenLineEndingY() >= c.getHeight() {
				// Try to fit next line.
				c.getCurrentEl().makeSmallestOffsetToFitLineOnDisplay()
				c.resyncBelow(c.data.Front())
			} else {
				c.resyncBelow(c.currentElement.Prev())
			}
		}
	case tcell.KeyDEL:
		{
			c.remove()
		}
	case tcell.KeyCtrlF:
		{
			c.dump()
		}
	default:
		{
			log.Printf("Key pressed %v (%v)", op.rn, op.key)
			c.insert(op.rn)
		}
	}
}
