package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

type physicalScreenHandler struct {
	screen tcell.Screen
}

func (s *physicalScreenHandler) close() {
	s.screen.Fini()
}

func (s *physicalScreenHandler) putStr(x, y int, b rune) {
	s.screen.SetContent(x, y, b, []rune{}, tcell.StyleDefault)
	s.screen.Show()
}

func (s *physicalScreenHandler) getSize() (int, int) {
	return s.screen.Size()
}

func (s *physicalScreenHandler) pollKeyPress() event {
	for {
		switch ev := s.screen.PollEvent().(type) {
		case *tcell.EventKey:
			return event{rn: ev.Rune(), k: ev.Key()}
		}
	}
}

func initPhysicalScreenHandler() screenHandler {
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

	return &physicalScreenHandler{screen: s}
}
