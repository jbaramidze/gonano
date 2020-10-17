package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

// ScreenHandler Bla bla
type ScreenHandler struct {
	screen tcell.Screen
}

// InitScreenHandler s
func InitScreenHandler() (*ScreenHandler, chan ContentOperation) {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error 1: %v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "Error 2:%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	channel := make(chan ContentOperation)

	return &ScreenHandler{screen: s}, channel
}

func (d ScreenHandler) putStr(x, y int, b rune) {
	d.screen.SetContent(x, y, b, []rune{}, tcell.StyleDefault)
	d.screen.Show()
}
