package main

import (
	"container/list"
	"reflect"
	"testing"

	"github.com/gdamore/tcell"
)

var emptyRow []rune = []rune{'@', '@', '@', '@'}

func initEditor(resp chan bool, w, h int) (*mockScreenHandler, *Editor) {
	handler := initMockScreenHandler(w, h)
	editor := createEditor(handler)

	blinkr := initMockBlinker(editor)
	editor.setBlinker(blinkr)

	go editor.startLoop()
	defer editor.display.close()

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
		ctx.t.Errorf("Display content is wrong! \n actual %v \n expected %v", ctx.h.data, data)
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
	if ctx.e.display.getBlinkerX() != x || ctx.e.display.getBlinkerY() != y {
		ctx.t.Errorf("Incorrect coords (%v, %v) vs (%v, %v)", ctx.e.display.getBlinkerX(), ctx.e.display.getBlinkerY(), x, y)
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
	h, e := initEditor(resp, 4, 6)
	ctx := context{h: h, resp: resp, t: t, e: e}

	expectPositionOnScreen(ctx, 0, 0)

	// Single line overflows by typing, enough space on screen
	sendChar(ctx, 'a')
	expectScreen(ctx, [][]rune{{'a', '@', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a'}})
	expectPositionOnScreen(ctx, 1, 0)

	sendChar(ctx, 'b')
	expectScreen(ctx, [][]rune{{'a', 'b', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b'}})
	expectPositionOnScreen(ctx, 2, 0)
	expectLineAndPosition(ctx, 0, 2)

	sendChar(ctx, 'c')
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', '@'}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c'}})
	expectPositionOnScreen(ctx, 3, 0)

	sendChar(ctx, 'd')
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'd'}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'd'}})
	expectPositionOnScreen(ctx, 0, 1)

	sendChar(ctx, 'e')
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'e', '@', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'd', 'e'}})
	expectPositionOnScreen(ctx, 1, 1)

	// Newline at last line, last character, enough space on screen
	sendKey(ctx, tcell.KeyEnter)
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'e', '@', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'd', 'e'}, {}})
	expectPositionOnScreen(ctx, 0, 2)

	sendChar(ctx, 'f')
	sendChar(ctx, 'g')
	sendChar(ctx, 'h')
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'e', '@', '@', '@'}, {'f', 'g', 'h', '@'}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'd', 'e'}, {'f', 'g', 'h'}})
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
	sendChar(ctx, 'i')
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'e', '@', '@', '@'}, {'i', 'f', 'g', 'h'}, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'd', 'e'}, {'i', 'f', 'g', 'h'}})
	expectPositionOnScreen(ctx, 1, 2)
	expectLineAndPosition(ctx, 1, 1)

	sendChar(ctx, 'j')
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'e', '@', '@', '@'}, {'i', 'j', 'f', 'g'}, {'h', '@', '@', '@'}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'd', 'e'}, {'i', 'j', 'f', 'g', 'h'}})
	expectPositionOnScreen(ctx, 2, 2)
	expectLineAndPosition(ctx, 1, 2)

	// Enter in the middle of line
	sendKey(ctx, tcell.KeyEnter)
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'e', '@', '@', '@'}, {'i', 'j', '@', '@'}, {'f', 'g', 'h', '@'}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'd', 'e'}, {'i', 'j'}, {'f', 'g', 'h'}})
	expectPositionOnScreen(ctx, 0, 3)
	expectLineAndPosition(ctx, 2, 0)

	// Enter for empty line
	sendKey(ctx, tcell.KeyEnter)
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'e', '@', '@', '@'}, {'i', 'j', '@', '@'}, emptyRow, {'f', 'g', 'h', '@'}, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'd', 'e'}, {'i', 'j'}, {}, {'f', 'g', 'h'}})
	expectPositionOnScreen(ctx, 0, 4)
	expectLineAndPosition(ctx, 3, 0)
}

func TestDeletes(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp, 4, 6)
	ctx := context{h: h, resp: resp, t: t, e: e}

	sendKey(ctx, tcell.KeyDEL)
	expectPositionOnScreen(ctx, 0, 0)

	// Delete in the middle of line
	sendChar(ctx, 'a')
	sendChar(ctx, 'b')
	sendChar(ctx, 'c')

	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{{'a', 'b', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b'}})
	expectPositionOnScreen(ctx, 2, 0)

	// Delete in the middle of 2nd line, both chars
	sendKey(ctx, tcell.KeyEnter)
	sendChar(ctx, 'd')
	sendChar(ctx, 'e')
	expectScreen(ctx, [][]rune{{'a', 'b', '@', '@'}, {'d', 'e', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b'}, {'d', 'e'}})
	expectPositionOnScreen(ctx, 2, 1)

	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{{'a', 'b', '@', '@'}, {'d', '@', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b'}, {'d'}})
	expectPositionOnScreen(ctx, 1, 1)

	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{{'a', 'b', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b'}, {}})
	expectPositionOnScreen(ctx, 0, 1)

	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{{'a', 'b', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b'}})
	expectPositionOnScreen(ctx, 2, 0)

	sendKey(ctx, tcell.KeyDEL)
	sendKey(ctx, tcell.KeyDEL)
	expectScreen(ctx, [][]rune{emptyRow, emptyRow, emptyRow, emptyRow, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{}})
	expectPositionOnScreen(ctx, 0, 0)
}

func TestScreenShift(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp, 4, 6)
	ctx := context{h: h, resp: resp, t: t, e: e}

	sendChar(ctx, 's')
	sendKey(ctx, tcell.KeyEnter)
	sendKey(ctx, tcell.KeyEnter)
	sendKey(ctx, tcell.KeyEnter)
	sendKey(ctx, tcell.KeyEnter)
	sendKey(ctx, tcell.KeyEnter)
	sendChar(ctx, 'a')

	expectScreen(ctx, [][]rune{{'s', '@', '@', '@'}, emptyRow, emptyRow, emptyRow, emptyRow, {'a', '@', '@', '@'}})

	// Screen shifts below by 1
	sendKey(ctx, tcell.KeyEnter)
	expectScreen(ctx, [][]rune{emptyRow, emptyRow, emptyRow, emptyRow, {'a', '@', '@', '@'}, emptyRow})
	sendChar(ctx, 'b')
	expectScreen(ctx, [][]rune{emptyRow, emptyRow, emptyRow, emptyRow, {'a', '@', '@', '@'}, {'b', '@', '@', '@'}})

	sendChar(ctx, 'c')
	sendChar(ctx, 'd')

	sendChar(ctx, 'e')
	expectScreen(ctx, [][]rune{emptyRow, emptyRow, emptyRow, {'a', '@', '@', '@'}, {'b', 'c', 'd', 'e'}, emptyRow})
}

func TestBasic2(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp, 4, 6)
	ctx := context{h: h, resp: resp, t: t, e: e}

	// Test case when first line starts overflowiing => further lines get shifted

	sendChar(ctx, 'a')
	sendChar(ctx, 'b')
	sendChar(ctx, 'c')
	sendKey(ctx, tcell.KeyEnter)
	sendChar(ctx, 'd')
	sendKey(ctx, tcell.KeyEnter)
	sendChar(ctx, 'e')
	sendChar(ctx, 'f')
	sendChar(ctx, 'g')
	sendChar(ctx, 'h')
	sendChar(ctx, 'i')

	// Make sure we good
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', '@'}, {'d', '@', '@', '@'}, {'e', 'f', 'g', 'h'}, {'i', '@', '@', '@'}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c'}, {'d'}, {'e', 'f', 'g', 'h', 'i'}})
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
	sendChar(ctx, 'j')
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'j'}, {'d', '@', '@', '@'}, {'e', 'f', 'g', 'h'}, {'i', '@', '@', '@'}, emptyRow, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'j'}, {'d'}, {'e', 'f', 'g', 'h', 'i'}})

	// Shift happening
	sendChar(ctx, 'k')
	expectScreen(ctx, [][]rune{{'a', 'b', 'c', 'j'}, {'k', '@', '@', '@'}, {'d', '@', '@', '@'}, {'e', 'f', 'g', 'h'}, {'i', '@', '@', '@'}, emptyRow})
	expectData(ctx, [][]rune{{'a', 'b', 'c', 'j', 'k'}, {'d'}, {'e', 'f', 'g', 'h', 'i'}})
	sendKey(ctx, tcell.KeyCtrlF)
}

func testOperation(ctx context, data [][]rune, offsetY, curLine, pos int, key tcell.Key) (int, int, int) {
	d := ctx.e.display
	d.data = list.New()
	for _, item := range data {
		d.data.PushBack(&Line{data: item, startingCoordY: 0, height: -1, pos: 0, display: d})
	}
	ctx.e.offsetY = offsetY
	d.currentElement = d.data.Front()
	for i := 0; i < curLine; i++ {
		d.currentElement = d.currentElement.Next()
	}
	d.getCurrentEl().pos = pos

	d.recalcBelow(d.data.Front())

	// Do the operation
	sendKey(ctx, tcell.KeyRight)

	currentLine := 0
	p := d.getCurrentEl().pos
	for d.currentElement != d.data.Front() {
		currentLine++
		d.currentElement = d.currentElement.Prev()
	}
	return ctx.e.offsetY, currentLine, p
}

func expectParams(ctx context, a, b, c, d, e, f int) {
	if a != d || b != e || c != f {
		ctx.t.Errorf("FAILED: %v!=%v or %v!=%v or %v!=%v", a, d, b, e, c, f)
	}

}

func TestArrows(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp, 2, 2)
	ctx := context{h: h, resp: resp, t: t, e: e}

	// /*
	// 	Right arrow test
	// */

	// // Cannot go any more right
	// // a, b, c := testOperation(ctx, [][]rune{{'a'}, {'b', 'c'}}, 0, 1, 2, tcell.KeyRight)
	// // expectParams(ctx, a, b, c, 0, 1, 2)
	a, b, c := testOperation(ctx, [][]rune{{'a'}, {'b', 'c', 'd', 'e'}}, 0, 1, 2, tcell.KeyRight)
	expectParams(ctx, a, b, c, 0, 1, 2)
}
