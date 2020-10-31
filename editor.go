package main

import "github.com/gdamore/tcell"

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
