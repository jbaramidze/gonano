package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Editor is main editor structure.
type Editor struct {
	display  *Display
	filename string
	modified bool
	mode     mode
}

func createEditor(handler screenHandler) *Editor {
	display := createDisplay(handler)
	editor := Editor{display: display, modified: false}
	editor.mode = normalMode{&editor}
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

	fields := strings.Split(string(data), "\n")
	for i, field := range fields {
		if i == 0 {
			e.display.getCurrentEl().data = []rune(field)
			e.display.getCurrentEl().pos = 0
		} else {
			newItem := Line{data: []rune(field), startingCoordY: -1, pos: 0, display: e.display}
			e.display.data.InsertAfter(&newItem, e.display.currentElement)
			e.display.currentElement = e.display.currentElement.Next()
		}
	}

	e.display.currentElement = e.display.data.Front()
	e.display.resyncBelow(e.display.currentElement)
}

func (e *Editor) saveData() error {
	data := []rune{}
	for it := e.display.data.Front(); it != nil; it = it.Next() {
		data = append(data, it.Value.(*Line).data...)
		data = append(data, rune(10))
	}

	err := ioutil.WriteFile(e.filename, []byte(string(data)), 0644)

	if err != nil {
		fmt.Printf("Error: failed writing to file: %v", err)
		return err
	}

	return nil
}

func (e *Editor) pollKeyboard(resp chan bool) {
	for {
		ev := e.display.screen.pollKeyPress()
		switch t := ev.(type) {
		case keyEvent:
			exit := e.mode.handleKeyPress(t, resp)
			if exit == true {
				return
			}
		case resizeEvent:
			e.display.resyncBelow(e.display.data.Front())
		}

	}
}

func (e *Editor) setMode(mode mode) {
	e.mode = mode
	e.mode.init()
}

func (e *Editor) setBlinker(b blinker) {
	e.display.blinker = b
}

func (e *Editor) setStausBar(s statusBar) {
	e.display.statusBar = s
}
