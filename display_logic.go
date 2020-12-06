package main

func oneCharOperation(c *Display, f func()) {
	oldH := c.getCurrentEl().getOnScreenLineEndingY()
	oldCursorY := c.getCurrentEl().getRelativeCursorY()
	f()
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
}

func (c *Display) insert(char rune) {
	oneCharOperation(c, func() {
		c.getCurrentEl().data = insertInSlice(c.getCurrentEl().data, char, c.getCurrentEl().pos)
		c.getCurrentEl().pos++
	})
}

func (c *Display) remove() {
	if c.getCurrentEl().pos == 0 {
		if !c.hasPrevEl() {
			return
		}

		// Remove current line
		p := c.currentElement.Prev()
		p.Value.(*Line).pos = len(p.Value.(*Line).data)
		p.Value.(*Line).data = append(p.Value.(*Line).data, c.getCurrentEl().data...)
		c.data.Remove(c.currentElement)
		c.currentElement = p
		c.recalcBelow(c.currentElement)

		// Fix Y!
		if c.getCurrentEl().getOnScreenCursorY() < 0 {
			c.offsetY--
		}
		c.resyncBelow(c.data.Front())
	} else {
		oneCharOperation(c, func() {
			c.getCurrentEl().pos--
			c.getCurrentEl().data = removeFromSlice(c.getCurrentEl().data, c.getCurrentEl().pos)
		})
	}
}
