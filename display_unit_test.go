package main

import (
	"testing"

	"github.com/gdamore/tcell"
)

func TestArrow(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp, 3, 3)
	ctx := context{h: h, resp: resp, t: t, e: e}

	/*
		Right arrow test
	*/

	// Can go more right
	setupScenario(ctx, [][]rune{{'a'}, {'b', 'c', 'd', 'e'}}, 1, 1, 2)
	sendKey(ctx, tcell.KeyRight)
	expectParams(ctx, 1, 1, 3)

	// Cannot go any more right
	setupScenario(ctx, [][]rune{{'a'}, {'b', 'c'}}, 1, 1, 2)
	sendKey(ctx, tcell.KeyRight)
	expectParams(ctx, 1, 1, 2)

	// If on last char, can cause offsetY increase
	setupScenario(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'b', 'c', 'd', 'e'}}, 0, 1, 2)
	sendKey(ctx, tcell.KeyRight)
	expectParams(ctx, 1, 1, 3)

	/*
		Left arrow test
	*/

	// Can go more left
	setupScenario(ctx, [][]rune{{'a'}, {'b', 'c', 'd', 'e'}}, 1, 1, 2)
	sendKey(ctx, tcell.KeyLeft)
	expectParams(ctx, 1, 1, 1)

	// Cannot go more left
	setupScenario(ctx, [][]rune{{'a'}, {'b', 'c', 'd', 'e'}}, 1, 1, 0)
	sendKey(ctx, tcell.KeyLeft)
	expectParams(ctx, 1, 1, 0)

	// If on first char on screen, can cause offsetY decrease
	setupScenario(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'b', 'c', 'd', 'e'}}, 3, 1, 3)
	sendKey(ctx, tcell.KeyLeft)
	expectParams(ctx, 2, 1, 2)

	/*
		Down arrow test
	*/

	// Cannot go more down
	setupScenario(ctx, [][]rune{{'a', 'b'}, {'b', 'c', 'd', 'e'}}, 1, 1, 2)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 1, 1, 2)

	// Move to next longer line
	setupScenario(ctx, [][]rune{{'a', 'b'}, {'b', 'c', 'd'}}, 0, 0, 2)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 0, 1, 2)

	// Move to next shorter line
	setupScenario(ctx, [][]rune{{'a', 'b', 'c'}, {'b', 'c'}}, 0, 0, 3)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 0, 1, 2)

	// Move to next line: visible beginning, non-visible end, fits screen
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b', 'c', 'd'}}, 0, 1, 1)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 1, 2, 1)

	// Move to next line: visible beginning, non-visible end, does not fit screen V1
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g'}}, 0, 1, 1)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 2, 2, 1)

	// Move to next line: visible beginning, non-visible end, does not fit screen V2
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}}, 0, 1, 1)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 2, 2, 1)
}
