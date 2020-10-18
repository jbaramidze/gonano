package main

import (
	"reflect"
	"testing"
)

var emptyRow []rune = []rune{0, 0, 0, 0}

func initDisplay() *mockScreenHandler {
	handler := initMockScreenHandler()
	display := createDisplay(handler)
	defer display.Close()

	go display.poll()

	return handler.(*mockScreenHandler)
}

func TestIntMinBasic(t *testing.T) {
	h := initDisplay()
	h.keyChan <- event{rn: 97}
	if !reflect.DeepEqual(h, [][]rune{emptyRow, emptyRow, emptyRow, emptyRow, emptyRow}) {
		t.Errorf("IntMin(2, -2); want -2")
	}
}
