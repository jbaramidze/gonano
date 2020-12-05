package main

import (
	"testing"

	"github.com/gdamore/tcell"
)

func TestArrow(t *testing.T) {
	resp := make(chan bool)
	h, e := initEditor(resp, 2, 2)
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
	setupScenario(ctx, [][]rune{{'a'}, {'b', 'c', 'd', 'e'}}, 0, 1, 1)
	sendKey(ctx, tcell.KeyRight)
	expectParams(ctx, 1, 1, 2)

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
	setupScenario(ctx, [][]rune{{'a'}, {'b', 'c', 'd', 'e'}}, 2, 1, 2)
	sendKey(ctx, tcell.KeyLeft)
	expectParams(ctx, 1, 1, 1)

}
