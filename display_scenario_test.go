package main

import (
	"testing"

	"github.com/gdamore/tcell"
)

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
