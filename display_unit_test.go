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

	// Move to next line: visible beginning, fits screen, visible part
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b', 'c', 'd'}}, 0, 1, 1)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 1, 2, 1)

	// Move to next line: visible beginning, fits screen, visible part, longer
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g'}}, 0, 1, 1)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 2, 2, 1)

	// Move to next line: visible beginning, fits screen, invisible part
	setupScenario(ctx, [][]rune{{'a'}, {'a', 'b', 'c', 'd'}, {'a', 'b', 'c', 'd'}}, 0, 1, 4)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 2, 2, 4)

	// Move to next line: visible beginning, does not fit screen, visible part
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}}, 0, 1, 1)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 2, 2, 1)

	// Move to next line: visible beginning, does not fit screen, invisible part
	setupScenario(ctx, [][]rune{{'a', 'b', 'c', 'd'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}}, 0, 0, 4)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 2, 1, 4)

	// Move to next line from line that's end is not visible, newline - 1 lines
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}, {'a'}}, 2, 2, 1)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 5, 3, 1)

	// Move to next line from line that's end is not visible, newline - 2 lines, jump to beginning
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}, {'a', 'b', 'c', 'd'}}, 2, 2, 4)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 6, 3, 4)

	// Move to next line: non-visible beginning, non-visible end, fits screen v1
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a'}, {'a', 'b'}}, 0, 2, 0)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 1, 3, 0)

	// Move to next line: non-visible beginning, non-visible end, fits screen v2
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a'}, {'a', 'b', 'c', 'd'}}, 0, 2, 0)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 2, 3, 0)

	// Move to next line: non-visible beginning, non-visible end, does not fit screen
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}}, 0, 2, 0)
	sendKey(ctx, tcell.KeyDown)
	expectParams(ctx, 3, 3, 0)

	/*
		Up arrow test
	*/

	// Cannot go more up
	setupScenario(ctx, [][]rune{{'a', 'b'}, {'b', 'c', 'd', 'e'}}, 0, 0, 2)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 0, 0, 2)

	// Move to previous longer line
	setupScenario(ctx, [][]rune{{'a', 'b', 'c'}, {'b', 'c'}}, 0, 1, 1)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 0, 0, 1)

	// Move to previous shorter line
	setupScenario(ctx, [][]rune{{'a', 'b'}, {'b', 'c', 'd'}}, 0, 1, 3)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 0, 0, 2)

	// Move to previous line: visible end, non-visible beginning, fits screen v1
	setupScenario(ctx, [][]rune{{'a'}, {'a'}, {'a', 'b'}, {'a', 'b', 'c', 'd'}}, 3, 3, 1)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 2, 2, 1)

	// Move to previous line: visible end, non-visible beginning, fits screen v2
	setupScenario(ctx, [][]rune{{'a'}, {'a', 'b', 'c', 'd'}, {'a', 'b', 'c', 'd'}}, 3, 2, 1)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 1, 1, 1)

	// Move to previous line: visible end, non-visible beginning, fits screen v3
	setupScenario(ctx, [][]rune{{'a'}, {'a', 'b', 'c', 'd'}, {'a', 'b', 'c', 'd'}}, 3, 2, 4)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 1, 1, 4)

	// Move to previous line: visible end, non-visible beginning, doesn't fit screen v1
	setupScenario(ctx, [][]rune{{'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}, {'a'}}, 6, 2, 0)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 1, 1, 0)

	// Move to previous line: visible end, non-visible beginning, doesn't fit screen v2
	setupScenario(ctx, [][]rune{{'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}, {'a', 'b', 'c', 'd'}}, 6, 2, 4)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 1, 1, 4)

	// Move to previous line: visible end, non-visible beginning, doesn't fit screen v2
	setupScenario(ctx, [][]rune{{'a'}, {'a', 'b', 'c', 'd', 'e', 'f', 'g', 'i', 'j', 'k', 'l', 'm', 'n'}, {'a', 'b', 'c', 'd'}}, 7, 2, 4)
	sendKey(ctx, tcell.KeyUp)
	expectParams(ctx, 1, 1, 4)
}
