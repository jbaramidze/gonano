package main

import (
	"reflect"
	"testing"

	"github.com/gdamore/tcell"
)

var emptyRow []rune = []rune{0, 0, 0, 0}

func initEditor(resp chan bool) (*mockScreenHandler, *Editor) {
	handler := initMockScreenHandler()
	editor := createEditor(handler)

	blinkr := initMockBlinker(editor)
	editor.setBlinker(blinkr)

	go editor.startLoop()
	defer editor.display.Close()

	go editor.pollKeyboard(resp)

	return handler.(*mockScreenHandler), editor
}

type context struct {
	h    *mockScreenHandler
	resp chan bool
	t    *testing.T
	e    *Editor
}

func sendChar(ctx context, c rune) {
	ctx.h.keyChan <- keyEvent{rn: c}
	<-ctx.resp
}

func sendKey(ctx context, k tcell.Key) {
	ctx.h.keyChan <- keyEvent{k: k}
	<-ctx.resp
}

func expectScreen(ctx context, data [][]rune) {
	if !reflect.DeepEqual(ctx.h.data, data) {
		ctx.t.Errorf("Display content is wrong: %v", ctx.h.data)
	}
}

func expectData(ctx context, data [][]rune) {
	sz := 0
	for i, j := ctx.e.display.data.Front(), 0; i != nil; i, j = i.Next(), j+1 {
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

func expectPositionOnScreen(ctx context, x int, y int) {
	if ctx.e.display.currentX != x || ctx.e.display.currentY != y {
		ctx.t.Errorf("Incorrect coords (%v, %v) vs (%v, %v)", ctx.e.display.currentX, ctx.e.display.currentY, x, y)
	}
}

func expectLineAndPosition(ctx context, line int, pos int) {
	firstLine := ctx.e.display.data.Front()
	for i := 0; i < line; i++ {
		firstLine = firstLine.Next()
	}
	if firstLine.Value != ctx.e.display.getCurrentEl() {
		ctx.t.Errorf("Incorrect line %v", line)
	}

	if ctx.e.display.getCurrentEl().pos != pos {
		ctx.t.Errorf("Incorrect pos %v vs %v", ctx.e.display.getCurrentEl().pos, pos)
	}
}

func TestBasic1(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp)
	ctx := context{h: h, resp: resp, t: t, e: e}

	expectPositionOnScreen(ctx, 0, 0)

	// Test typing at the end of line, overflowing
	sendChar(ctx, 97)
	expectScreen(ctx, [][]rune{{97, 0, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97}})
	expectPositionOnScreen(ctx, 1, 0)

	sendChar(ctx, 98)
	expectScreen(ctx, [][]rune{{97, 98, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98}})
	expectPositionOnScreen(ctx, 2, 0)
	expectLineAndPosition(ctx, 0, 2)

	sendChar(ctx, 99)
	expectScreen(ctx, [][]rune{{97, 98, 99, 0}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99}})
	expectPositionOnScreen(ctx, 3, 0)

	sendChar(ctx, 100)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100}})
	expectPositionOnScreen(ctx, 0, 1)

	sendChar(ctx, 101)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}})
	expectPositionOnScreen(ctx, 1, 1)

	// Test newline on last line
	sendKey(ctx, tcell.KeyEnter)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}, {}})
	expectPositionOnScreen(ctx, 0, 2)

	sendChar(ctx, 102)
	sendChar(ctx, 103)
	sendChar(ctx, 104)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, {102, 103, 104, 0}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}, {102, 103, 104}})
	expectPositionOnScreen(ctx, 3, 2)
	expectLineAndPosition(ctx, 1, 3)

	/*
	  T E S T   A R R O W S
	*/
	// Up - takes us to previous, longer line
	sendKey(ctx, tcell.KeyUp)
	expectPositionOnScreen(ctx, 3, 0)
	expectLineAndPosition(ctx, 0, 3)
	// Up - first line, cannot go further
	sendKey(ctx, tcell.KeyUp)
	expectPositionOnScreen(ctx, 3, 0)
	expectLineAndPosition(ctx, 0, 3)
	// Right - jump to next line
	sendKey(ctx, tcell.KeyRight)
	expectPositionOnScreen(ctx, 0, 1)
	// Right - regular
	sendKey(ctx, tcell.KeyRight)
	expectPositionOnScreen(ctx, 1, 1)
	// Right - last char, cannot go further
	sendKey(ctx, tcell.KeyRight)
	expectPositionOnScreen(ctx, 1, 1)
	expectLineAndPosition(ctx, 0, 5)
	// Down - takes us to next, shorter line
	sendKey(ctx, tcell.KeyDown)
	expectPositionOnScreen(ctx, 3, 2)
	expectLineAndPosition(ctx, 1, 3)
	// Down - cannot go any more below, last line
	sendKey(ctx, tcell.KeyDown)
	expectPositionOnScreen(ctx, 3, 2)
	expectLineAndPosition(ctx, 1, 3)
	// Left - go to previus location
	sendKey(ctx, tcell.KeyLeft)
	expectPositionOnScreen(ctx, 2, 2)
	expectLineAndPosition(ctx, 1, 2)
	// 2 more times, go to beginning of line
	sendKey(ctx, tcell.KeyLeft)
	sendKey(ctx, tcell.KeyLeft)
	expectPositionOnScreen(ctx, 0, 2)
	expectLineAndPosition(ctx, 1, 0)
	// Cannot go any more left
	sendKey(ctx, tcell.KeyLeft)
	expectPositionOnScreen(ctx, 0, 2)
	expectLineAndPosition(ctx, 1, 0)

	/*
	  T Y P I N G
	*/
	// Type at the beginning of line
	sendChar(ctx, 105)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, {105, 102, 103, 104}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}, {105, 102, 103, 104}})
	expectPositionOnScreen(ctx, 1, 2)
	expectLineAndPosition(ctx, 1, 1)

	sendChar(ctx, 106)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, {105, 106, 102, 103}, {104, 0, 0, 0}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}, {105, 106, 102, 103, 104}})
	expectPositionOnScreen(ctx, 2, 2)
	expectLineAndPosition(ctx, 1, 2)

	// Enter in the middle of line
	sendKey(ctx, tcell.KeyEnter)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, {105, 106, 0, 0}, {102, 103, 104, 0}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}, {105, 106}, {102, 103, 104}})
	expectPositionOnScreen(ctx, 0, 3)
	expectLineAndPosition(ctx, 2, 0)

	// Enter for empty line
	sendKey(ctx, tcell.KeyEnter)
	expectScreen(ctx, [][]rune{{97, 98, 99, 100}, {101, 0, 0, 0}, {105, 106, 0, 0}, emptyRow, {102, 103, 104, 0}, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 100, 101}, {105, 106}, {}, {102, 103, 104}})
	expectPositionOnScreen(ctx, 0, 4)
	expectLineAndPosition(ctx, 3, 0)
}

func TestDeletes(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp)
	ctx := context{h: h, resp: resp, t: t, e: e}

	sendKey(ctx, tcell.KeyDEL)
	expectPositionOnScreen(ctx, 0, 0)

	// Delete in the middle of line
	sendChar(ctx, 97)
	sendChar(ctx, 98)
	sendChar(ctx, 99)

	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{{97, 98, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98}})
	expectPositionOnScreen(ctx, 2, 0)

	// Delete in the middle of 2nd line, both chars
	sendKey(ctx, tcell.KeyEnter)
	sendChar(ctx, 100)
	sendChar(ctx, 101)
	expectScreen(ctx, [][]rune{{97, 98, 0, 0}, {100, 101, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98}, {100, 101}})
	expectPositionOnScreen(ctx, 2, 1)

	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{{97, 98, 0, 0}, {100, 0, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98}, {100}})
	expectPositionOnScreen(ctx, 1, 1)

	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{{97, 98, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98}, {}})
	expectPositionOnScreen(ctx, 0, 1)

	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{{97, 98, 0, 0}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98}})
	expectPositionOnScreen(ctx, 2, 0)

	sendKey(ctx, tcell.KeyDEL)
	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{emptyRow, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{}})
	expectPositionOnScreen(ctx, 0, 0)
}

func TestBasic2(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp)
	ctx := context{h: h, resp: resp, t: t, e: e}

	// Test case when first line starts overflowiing => further lines get shifted

	sendChar(ctx, 97)
	sendChar(ctx, 98)
	sendChar(ctx, 99)
	sendKey(ctx, tcell.KeyEnter)
	sendChar(ctx, 100)
	sendKey(ctx, tcell.KeyEnter)
	sendChar(ctx, 101)
	sendChar(ctx, 102)
	sendChar(ctx, 103)
	sendChar(ctx, 104)
	sendChar(ctx, 105)

	// Make sure we good
	expectScreen(ctx, [][]rune{{97, 98, 99, 0}, {100, 0, 0, 0}, {101, 102, 103, 104}, {105, 0, 0, 0}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99}, {100}, {101, 102, 103, 104, 105}})
	expectPositionOnScreen(ctx, 1, 3)
	expectLineAndPosition(ctx, 2, 5)

	// Jump to the beginning and overflow by typing something
	sendKey(ctx, tcell.KeyUp)
	sendKey(ctx, tcell.KeyUp)
	sendKey(ctx, tcell.KeyRight)
	sendKey(ctx, tcell.KeyRight)
	sendKey(ctx, tcell.KeyRight)
	expectPositionOnScreen(ctx, 3, 0)
	expectLineAndPosition(ctx, 0, 3)

	// Text full, console newline but still no shift
	sendChar(ctx, 106)
	expectScreen(ctx, [][]rune{{97, 98, 99, 106}, {100, 0, 0, 0}, {101, 102, 103, 104}, {105, 0, 0, 0}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 106}, {100}, {101, 102, 103, 104, 105}})

	// Shift happening
	sendChar(ctx, 107)
	expectScreen(ctx, [][]rune{{97, 98, 99, 106}, {107, 0, 0, 0}, {100, 0, 0, 0}, {101, 102, 103, 104}, {105, 0, 0, 0}, emptyRow})
	expectData(ctx, [][]rune{{97, 98, 99, 106, 107}, {100}, {101, 102, 103, 104, 105}})
	sendKey(ctx, tcell.KeyCtrlF)
}
