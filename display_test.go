package main

import (
	"reflect"
	"testing"
)

var emptyRow []rune = []rune{0, 0, 0, 0}

func initDisplay() *mockScreenHandler {
	handler := initMockScreenHandler()
	display := createDisplay(handler)

	blinkr := initMockBlinker(display)
	display.setBlinker(blinkr)

	go display.startLoop()
	defer display.Close()

	go display.pollKeyboard()

	return handler.(*mockScreenHandler)
}

// func maybeLate(arg func(a ...interface{}) bool) {
// 	for i := 0; i < 5; i++ {
// 		if arg() {
// 			return true
// 		}
// 	}
// }

func TestIntMinBasic(t *testing.T) {
	h := initDisplay()
	h.keyChan <- event{rn: 97}
	if !reflect.DeepEqual(h.data, [][]rune{{97, 32, 0, 0}, emptyRow, emptyRow, emptyRow}) {
		t.Errorf("Display content is wrong in [A]: %v", h.data)
	}
}
