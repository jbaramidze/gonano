package main

import "container/list"

func (c *Display) insert(char rune) {
	oldH := c.getCurrentEl().getOnScreenLineEndingY()
	oldCursorY := c.getCurrentEl().getRelativeCursorY()
	c.getCurrentEl().data = insertInSlice(c.getCurrentEl().data, char, c.getCurrentEl().pos)
	c.getCurrentEl().pos++
	newH := c.getCurrentEl().display.getHeight()
	newCursorY := c.getCurrentEl().getOnScreenLineEndingY()

	if oldCursorY != newCursorY {
		onScreenCursorY := c.getCurrentEl().getOnScreenCursorY()
		if onScreenCursorY >= c.getHeight() {
			c.offsetY++
			c.resyncBelow(c.data.Front())
		} else {
			c.resyncBelow(c.currentElement)
		}
	} else {
		if oldH != newH {
			c.resyncBelow(c.currentElement)
		} else {
			c.getCurrentEl().resync()
		}
	}
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

		c.resyncBelow(c.currentElement)
	} else {
		oldHeight := c.getCurrentEl().calculateHeight()
		c.getCurrentEl().pos--
		c.getCurrentEl().data = removeFromSlice(c.getCurrentEl().data, c.getCurrentEl().pos)
		newHeight := c.getCurrentEl().calculateHeight()

		if newHeight != oldHeight {
			c.resyncBelow(c.currentElement)
		} else {
			c.getCurrentEl().resync()
		}
	}

	c.syncCoords()
}

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
