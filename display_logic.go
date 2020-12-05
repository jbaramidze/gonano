package main

func (c *Display) insert(char rune) {
	oldH := c.getCurrentEl().getOnScreenLineEndingY()
	oldCursorY := c.getCurrentEl().getRelativeCursorY()
	c.getCurrentEl().data = insertInSlice(c.getCurrentEl().data, char, c.getCurrentEl().pos)
	c.getCurrentEl().pos++
	newH := c.getCurrentEl().getOnScreenLineEndingY()
	newCursorY := c.getCurrentEl().getRelativeCursorY()

	if oldCursorY != newCursorY {
		c.resyncNewCursorY()
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

		// Remove current line
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
