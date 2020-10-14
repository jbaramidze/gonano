package display

import (
	"github.com/gdamore/tcell"
	"github.com/jbaramidze/term_collab_editor/helper"
)

func (c *Display) handleKeyPress(op TypeOperation) {
	switch op.key {
	case tcell.KeyLeft:
		{
			if c.currentX > 0 {
				c.currentX--
			}
		}
	case tcell.KeyRight:
		{
			if len(c.getCurrentEl().data) > c.currentX {
				c.currentX++
			}
		}
	case tcell.KeyUp:
		{
			if c.currentY > 0 {
				c.currentY--
				if c.getCurrentEl().currentY > 0 {
					c.getCurrentEl().currentY--
				} else {
					c.currentX = helper.MinOf(
						c.currentX,
						len(c.getCurrentEl().data),
						len(c.getPrevEl().data))
					c.currentElement = c.currentElement.Prev()
				}

			}
		}
	case tcell.KeyDown:
		{
			if c.data.Len() > c.currentY+1 {
				c.currentY++
				c.currentX = helper.MinOf(
					c.currentX,
					len(c.getCurrentEl().data),
					len(c.getNextEl().data))
				c.currentElement = c.currentElement.Next()
			}
		}
	case tcell.KeyEnter:
		{
			cur := c.getCurrentEl()
			newItem := Line{data: []rune{}, startingCoordY: cur.startingCoordY + cur.height, height: 1}
			c.data.InsertAfter(&newItem, c.currentElement)
			c.currentElement = c.currentElement.Next()
			c.currentX = 0
			c.currentY++
		}
	default:
		{
			c.newChar(op.rn)
		}
	}
}
