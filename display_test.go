package main

import (
	"reflect"
	"testing"
)

var emptyRow []rune = []rune{0, 0, 0, 0}

func initDisplay(resp chan bool) (*mockScreenHandler, *Display) {
	handler := initMockScreenHandler()
	display := createDisplay(handler)

	blinkr := initMockBlinker(display)
	display.setBlinker(blinkr)

	go display.startLoop()
	defer display.Close()

	go display.pollKeyboard(resp)

	return handler.(*mockScreenHandler), display
}

type context struct {
	h    *mockScreenHandler
	resp chan bool
	t    *testing.T
	d    *Display
}

func sendChar(ctx context, c rune) {
	ctx.h.keyChan <- event{rn: c}
	<-ctx.resp
}

func expectScreen(ctx context, data [][]rune) {
	if !reflect.DeepEqual(ctx.h.data, data) {
		ctx.t.Errorf("Display content is wrong in [A]: %v", ctx.h.data)
	}
}

func expectData(ctx context, data [][]rune) {
	sz := 0
	for i, j := ctx.d.data.Front(), 0; i != nil; i, j = i.Next(), j+1 {
		l := i.Value.(*Line)
		sz++
		if !reflect.DeepEqual(l.data, data[j]) {
			ctx.t.Errorf("Data content is wrong in [A]: %v vs %v", l.data, data[j])
		}
	}
	if sz != len(data) {
		ctx.t.Errorf("Data not of same length: %v vs %v", sz, len(data))
	}
}

func TestIntMinBasic(t *testing.T) {
	resp := make(chan bool)
	h, d := initDisplay(resp)
	ctx := context{h: h, resp: resp, t: t, d: d}

	sendChar(ctx, 97)
	expectScreen(ctx, [][]rune{{97, 32, 0, 0}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97}})

	sendChar(ctx, 98)
	expectScreen(ctx, [][]rune{{97, 98, 32, 0}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98}})
}
