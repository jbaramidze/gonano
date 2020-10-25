package main

import (
	"time"
)

type blinker interface {
	refresh()
	set()
	clear()
}

type realBlinker struct {
	blinkIsSet bool
	d          *Display
}

func (r *realBlinker) refresh() {
	if r.blinkIsSet {
		r.set()
	} else {
		r.clear()
	}
}

func (r *realBlinker) set() {
	r.d.screen.putStr(r.d.currentX, r.d.currentY, rune('▉'))
}

func (r *realBlinker) clear() {
	if len(r.d.getCurrentEl().data) > r.d.getCurrentEl().pos {
		r.d.screen.putStr(r.d.currentX, r.d.currentY, r.d.getCurrentEl().getCurrentChar())
	} else {
		// FIXME: exception: it might be on beginning of another line. Fix the case.
		if r.d.currentX == 0 && r.d.hasNextEl() && r.d.currentY == r.d.getNextEl().startingCoordY {
			if len(r.d.getNextEl().data) > 0 {
				r.d.screen.putStr(r.d.currentX, r.d.currentY, r.d.getNextEl().data[0])
			} else {
				r.d.screen.putStr(r.d.currentX, r.d.currentY, rune(0))
			}
		} else {
			r.d.screen.putStr(r.d.currentX, r.d.currentY, rune(0))
		}
	}
}

func initRealBlinker(d *Display) blinker {
	b := &realBlinker{d: d, blinkIsSet: false}
	go func(c chan ContentOperation) {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			<-ticker.C
			b.blinkIsSet = !b.blinkIsSet
			c <- BlinkOperation{blink: b.blinkIsSet}
		}
	}(d.monitorChannel)

	return b
}

type mockBlinker struct {
	rb    realBlinker
	ticks chan bool
}

func (b *mockBlinker) refresh() {
	b.rb.refresh()
}
func (b *mockBlinker) set() {
	b.rb.set()
}
func (b *mockBlinker) clear() {
	b.rb.clear()
}

func initMockBlinker(d *Display) blinker {
	ticker := make(chan bool)
	b := &realBlinker{d: d, blinkIsSet: false}
	go func(c chan ContentOperation) {
		for {
			<-ticker
			b.blinkIsSet = !b.blinkIsSet
			c <- BlinkOperation{blink: b.blinkIsSet}
		}
	}(d.monitorChannel)

	return b
}
