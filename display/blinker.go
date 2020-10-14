package display

import "time"

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
	if len(d.getCurrentEl().data) > d.currentX {
		d.ScreenHandler.putStr(d.currentX, d.currentY, rune(d.getCurrentEl().data[d.currentX]))
	} else {
		d.ScreenHandler.putStr(d.currentX, d.currentY, rune(' '))
	}
}
func (d *Display) setBlinkStatus() {
	d.ScreenHandler.putStr(d.currentX, d.currentY, rune('▉'))
}
