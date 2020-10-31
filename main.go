package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Error: Please pass filename as argument")
		os.Exit(1)
	}

	arg := os.Args[1]

	handler := initPhysicalScreenHandler()
	editor := createEditor(handler)
	editor.initData(arg)

	blinkr := initRealBlinker(editor)
	editor.setBlinker(blinkr)

	go editor.startLoop()
	defer editor.display.Close()

	editor.pollKeyboard(nil)
}
