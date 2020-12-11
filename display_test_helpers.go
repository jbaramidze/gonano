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
		d.data.PushBack(&Line{data: item, startingCoordY: 0, pos: 0, display: d})
	}
	d.offsetY = offsetY
	d.currentElement = d.data.Front()
	for i := 0; i < curLine; i++ {
		d.currentElement = d.currentElement.Next()
	}
	d.getCurrentEl().pos = pos

	d.resyncBelow(d.data.Front())
}

func validateScreenContent(ctx context) {
	// First make sure Ys are calculated correctly
	y := 0
	for it := ctx.e.display.data.Front(); it != ctx.e.display.data.Back(); it = it.Next() {
		val := it.Value.(*Line)
		if y != val.startingCoordY {
			ctx.t.Errorf("Incorrect OffsetY when validating! %v != %v", y, val.startingCoordY)
			return
		}
		y += val.calculateHeight()
	}

	startingY := ctx.e.display.offsetY
	line := ctx.e.display.data.Front()

	// Skip lines that end before screen begins
	for line != nil && line.Value.(*Line).startingCoordY+line.Value.(*Line).calculateHeight()-1 < startingY {
		line = line.Next()
	}

	if line == nil {
		return
	}

	// Skip wrapped lines insude line before screen begins
	pos := 0
	for line.Value.(*Line).startingCoordY+pos/ctx.h.w < startingY {
		pos += ctx.h.w
	}

	i := 0
	for i < ctx.h.w*ctx.h.w {
		if line == nil {
			// make sure it's '@'
			return
		}
		for i < ctx.h.w*ctx.h.w && pos < len(line.Value.(*Line).data) {
			w := i % ctx.h.w
			h := i / ctx.h.w
			if line.Value.(*Line).data[pos] != ctx.h.data[h][w] {
				ctx.t.Errorf("Byte mismatch when validating! %c != %c w=%v h=%v pos=%v screen: %v line: %v",
					ctx.h.data[h][w], line.Value.(*Line).data[pos], w, h, pos, ctx.h.data, line.Value.(*Line).data)
				return
			}
			i++
			pos++
		}

		if i == ctx.h.w*ctx.h.w {
			return
		}

		// If line was empty
		if len(line.Value.(*Line).data) == 0 {
			for j := 0; j < ctx.h.w; j, i = j+1, i+1 {
				if '@' != ctx.h.data[i/ctx.h.w][j] {
					ctx.t.Errorf("Byte mismatch when validating! %c != @", ctx.h.data[i/ctx.h.w][j])
					return
				}
			}
			line = line.Next()
			pos = 0
			continue
		}

		// If finished exactly
		if i%ctx.h.w == 0 {
			line = line.Next()
			pos = 0
			continue
		}

		// Finalize this line with '@'s
		remaining := ctx.h.w - i%ctx.h.w
		for j := 0; j < remaining; j, i = j+1, i+1 {
			w := i % ctx.h.w
			h := i / ctx.h.w
			if '@' != ctx.h.data[h][w] {
				ctx.t.Errorf("Byte mismatch when validating! %c != @", ctx.h.data[h][w])
				return
			}
		}
		line = line.Next()
		pos = 0
		continue
	}

}

func expectScenario(ctx context, data [][]rune, offsetY, curLine, pos int) {
	expectData(ctx, data)
	expectParams(ctx, offsetY, curLine, pos)
	validateScreenContent(ctx)
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
