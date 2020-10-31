package main

import (
	"strings"

	"github.com/gdamore/tcell"
)

// Editor is main editor structure.
type Editor struct {
	display *Display
}

func createEditor(handler screenHandler) *Editor {
	display := createDisplay(handler)
	return &Editor{display: display}
}

func (e *Editor) startLoop() {
	e.display.startLoop()
}

func (e *Editor) initData(data []byte) {
	fields := strings.Fields(string(data))
	for i, field := range fields {
		if i == 0 {
			e.display.getCurrentEl().data = []rune(field)
			e.display.getCurrentEl().pos = len(field)
		} else {
			newItem := Line{data: []rune(field), startingCoordY: -1, height: -1, pos: 0, display: e.display}
			e.display.data.InsertAfter(&newItem, e.display.currentElement)
			e.display.currentElement = e.display.currentElement.Next()
		}
	}

	// Move current to the beginning to resync
	e.display.currentElement = e.display.data.Front()
	e.display.resyncAll()

	// Leave it at the end
	e.display.currentElement = e.display.data.Back()
	e.display.getCurrentEl().pos = len(e.display.getCurrentEl().data)
	e.display.syncCoords()

}

func (e *Editor) pollKeyboard(resp chan bool) {
	for {
		ev := e.display.screen.pollKeyPress()
		if ev.k == tcell.KeyTAB {
			return
		}
		e.display.monitorChannel <- TypeOperation{rn: ev.rn, key: ev.k, resp: resp}
	}
}

func (e *Editor) setBlinker(b blinker) {
	e.display.blinker = b
}
