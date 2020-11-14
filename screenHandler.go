package main

import "github.com/gdamore/tcell"

type screenHandler interface {
	putStr(x, y int, b rune)
	getSize() (int, int)
	pollKeyPress() interface{}
	close()
}

type keyEvent struct {
	rn rune
	// Instead copy-pasting and mapping all constants....
	k tcell.Key
}

type resizeEvent struct {
}
