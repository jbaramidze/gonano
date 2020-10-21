package main

import (
	"reflect"
	"testing"

	"github.com/gdamore/tcell"
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

func sendKey(ctx context, k tcell.Key) {
	ctx.h.keyChan <- event{k: k}
	<-ctx.resp
}

func expectScreen(ctx context, data [][]rune) {
	if !reflect.DeepEqual(ctx.h.data, data) {
		ctx.t.Errorf("Display content is wrong: %v", ctx.h.data)
	}
}

func expectData(ctx context, data [][]rune) {
	sz := 0
	for i, j := ctx.d.data.Front(), 0; i != nil; i, j = i.Next(), j+1 {
		l := i.Value.(*Line)
		sz++
		if !reflect.DeepEqual(l.data, data[j]) {
			ctx.t.Errorf("Data content is wrong: %v vs %v", l.data, data[j])
		}
	}
	if sz != len(data) {
		ctx.t.Errorf("Data not of same length: %v vs %v", sz, len(data))
	}
}

func expectPosition(ctx context, x int, y int) {
	if ctx.d.currentX != x || ctx.d.currentY != y {
		ctx.t.Errorf("Incorrect coords (%v, %v) vs (%v, %v)", ctx.d.currentX, ctx.d.currentY, x, y)
	}
}

func TestIntMinBasic(t *testing.T) {
	resp := make(chan bool)
	h, d := initDisplay(resp)
	ctx := context{h: h, resp: resp, t: t, d: d}

	expectPosition(ctx, 0, 0)

	sendChar(ctx, 97)
	expectScreen(ctx, [][]rune{{97, 0, 0, 0}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97}})
	expectPosition(ctx, 1, 0)

	sendChar(ctx, 98)
	expectScreen(ctx, [][]rune{{97, 98, 0, 0}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98}})
	expectPosition(ctx, 2, 0)

	sendChar(ctx, 99)
	expectScreen(ctx, [][]rune{{97, 98, 99, 0}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99}})
	expectPosition(ctx, 3, 0)

	sendChar(ctx, 100)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100}})
	expectPosition(ctx, 0, 1)

	sendChar(ctx, 101)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}})
	expectPosition(ctx, 1, 1)

	sendKey(ctx, tcell.KeyEnter)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}, {}})
	expectPosition(ctx, 0, 2)

	sendChar(ctx, 102)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, {102, 0, 0, 0}, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}, {102}})
	expectPosition(ctx, 1, 2)
}
