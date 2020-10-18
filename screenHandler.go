package main

import "github.com/gdamore/tcell"

type screenHandler interface {
	putStr(x, y int, b rune)
	getSize() (int, int)
	pollKeyPress() event
	close()
}

type event struct {
	rn rune
	// Instead copy-pasting and mapping all constants....
	k tcell.Key
}
