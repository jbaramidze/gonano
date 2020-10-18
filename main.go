package main

func main() {
	handler := initPhysicalScreenHandler()
	display := createDisplay(handler)
	defer display.Close()

	display.poll()
}
