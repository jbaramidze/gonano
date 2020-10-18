package main

import (
	"time"
)

func initAndStartBlinker(c chan ContentOperation) {
	go func(c chan ContentOperation) {
		ticker := time.NewTicker(500 * time.Millisecond)
		init := false
		for {
			init = !init
			<-ticker.C
			c <- BlinkOperation{blink: init}
		}
	}(c)
}

func (d *Display) refreshBlinkStatus() {
	if d.blinkIsSet {
		d.setBlinkStatus()
	} else {
		d.clearBlinkStatus()
	}
}

func (d *Display) clearBlinkStatus() {
	if len(d.getCurrentEl().data) > d.getCurrentEl().pos {
		d.screen.putStr(d.currentX, d.currentY, d.getCurrentEl().getCurrentChar())
	} else {
		d.screen.putStr(d.currentX, d.currentY, rune(' '))
	}
}
func (d *Display) setBlinkStatus() {
	d.screen.putStr(d.currentX, d.currentY, rune('▉'))
}
