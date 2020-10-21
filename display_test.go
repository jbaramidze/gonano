package main

import (
	"reflect"
	"testing"
)

var emptyRow []rune = []rune{0, 0, 0, 0}

func initDisplay(resp chan bool) *mockScreenHandler {
	handler := initMockScreenHandler()
	display := createDisplay(handler)

	blinkr := initMockBlinker(display)
	display.setBlinker(blinkr)

	go display.startLoop()
	defer display.Close()

	go display.pollKeyboard(resp)

	return handler.(*mockScreenHandler)
}

type context struct {
	h    *mockScreenHandler
	resp chan bool
	t    *testing.T
}

func sendChar(ctx context, c rune) {
	ctx.h.keyChan <- event{rn: c}
	<-ctx.resp
}

func expect(ctx context, data [][]rune) {
	if !reflect.DeepEqual(ctx.h.data, data) {
		ctx.t.Errorf("Display content is wrong in [A]: %v", ctx.h.data)
	}
}

func TestIntMinBasic(t *testing.T) {
	resp := make(chan bool)
	h := initDisplay(resp)
	ctx := context{h: h, resp: resp, t: t}

	sendChar(ctx, 97)
	expect(ctx, [][]rune{{97, 32, 0, 0}, emptyRow, emptyRow, emptyRow})

	sendChar(ctx, 98)
	expect(ctx, [][]rune{{97, 98, 32, 0}, emptyRow, emptyRow, emptyRow})
}
