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
	statusBar := newPhysicalStatusBar(editor.display)
	editor.setBlinker(blinkr)
	editor.setStausBar(statusBar)

	go editor.startLoop()
	defer editor.display.close()

	editor.pollKeyboard(nil)
}
