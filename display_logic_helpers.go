package main

import (
	"container/list"
)

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

func (c *Display) recalcBelow(from *list.Element) {
	startingY := from.Value.(*Line).startingCoordY
	for ; from != nil; from = from.Next() {
		line := from.Value.(*Line)
		line.startingCoordY = startingY
		line.height = line.calculateHeight()
		startingY += line.height
	}
}

func (c *Display) resyncNewCursorY() {
	onScreenCursorY := c.getCurrentEl().getOnScreenCursorY()
	// If cursor jumped below screen
	if onScreenCursorY >= c.getHeight() {
		c.offsetY++
		c.resyncBelow(c.data.Front())
	} else if onScreenCursorY < 0 {
		c.offsetY--
		c.resyncBelow(c.data.Front())
	} else {
		c.resyncBelow(c.currentElement)
	}
}
