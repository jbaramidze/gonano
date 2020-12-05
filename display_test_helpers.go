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

func setupScenario(ctx context, data [][]rune, offsetY, curLine, pos int) {
	d := ctx.e.display
	d.data = list.New()
	for _, item := range data {
		d.data.PushBack(&Line{data: item, startingCoordY: 0, height: -1, pos: 0, display: d})
	}
	d.offsetY = offsetY
	d.currentElement = d.data.Front()
	for i := 0; i < curLine; i++ {
		d.currentElement = d.currentElement.Next()
	}
	d.getCurrentEl().pos = pos

	d.recalcBelow(d.data.Front())
}

func expectParams(ctx context, expectY, expectLine, expectPos int) {
	d := ctx.e.display

	if expectY != d.offsetY {
		ctx.t.Errorf("Incorrect OffsetY! Expected %v actual %v", expectY, d.offsetY)
		return
	}

	currentLine := 0
	for t := d.data.Front(); t != d.currentElement; t = t.Next() {
		currentLine++
	}
	if expectLine != currentLine {
		ctx.t.Errorf("Incorrect currentLine! Expected %v actual %v", expectLine, currentLine)
		return
	}

	if expectPos != d.getCurrentEl().pos {
		ctx.t.Errorf("Incorrect pos! Expected %v actual %v", expectPos, d.getCurrentEl().pos)
		return
	}
}
