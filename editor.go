package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gdamore/tcell"
)

// Editor is main editor structure.
type Editor struct {
	display  *Display
	filename string
	modified bool
	mode     editoMode
}

type editoMode int

const (
	normal editoMode = iota
	quitWithoutSaving
)

func createEditor(handler screenHandler) *Editor {
	display := createDisplay(handler)
	editor := Editor{display: display, modified: false}
	return &editor
}

func (e *Editor) startLoop() {
	e.display.startLoop()
}

func (e *Editor) initData(filename string) {
	e.filename = filename

	data, err := ioutil.ReadFile(e.filename)
	if err != nil {
		fmt.Printf("Error: Opening file failed: %v", err)
		return
	}

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
	e.display.resyncBelowCurrent()

	// Leave it at the end
	e.display.currentElement = e.display.data.Back()
	e.display.getCurrentEl().pos = len(e.display.getCurrentEl().data)
	e.display.syncCoords()

}

func (e *Editor) saveData() {
	data := []rune{}
	for it := e.display.data.Front(); it != nil; it = it.Next() {
		data = append(data, it.Value.(*Line).data...)
		data = append(data, rune(10))
	}

	err := ioutil.WriteFile(e.filename, []byte(string(data)), 0644)

	if err != nil {
		fmt.Printf("Error: failed writing to file: %v", err)
		return
	}
}

func (e *Editor) pollKeyboard(resp chan bool) {
	for {
		ev := e.display.screen.pollKeyPress()
		if ev.k == tcell.KeyCtrlQ {
			if e.modified == false {
				return
			}
			e.display.monitorChannel <- AnnouncementOperation{text: []string{"You will lose your changes!", "Are you sure you want to quit? Y/N"}}
		} else if ev.k == tcell.KeyCtrlW {
			e.saveData()
		} else {
			e.modified = true
			e.display.monitorChannel <- TypeOperation{rn: ev.rn, key: ev.k, resp: resp}
		}
	}
}

func (e *Editor) setBlinker(b blinker) {
	e.display.blinker = b
}
