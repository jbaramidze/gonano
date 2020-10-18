package main

func main() {
	handler := initPhysicalScreenHandler()
	display := createDisplay(handler)

	blinkr := initRealBlinker(display)
	display.setBlinker(blinkr)

	go display.startLoop()
	defer display.Close()

	display.pollKeyboard()
}
