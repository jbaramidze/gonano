package main

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

func initMockBlinker(e *Editor) blinker {
	ticker := make(chan bool)
	b := &realBlinker{d: e.display, blinkIsSet: false}
	go func(c chan contentOperation) {
		for {
			<-ticker
			b.blinkIsSet = !b.blinkIsSet
			c <- blinkOperation{blink: b.blinkIsSet}
		}
	}(e.display.monitorChannel)

	return b
}
