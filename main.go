package main

func main() {
	display := createDisplay()
	defer display.Close()

	display.poll()
}
