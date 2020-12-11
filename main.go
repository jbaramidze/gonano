package main

import (
	"fmt"
	"os"
)

func printHelp() {
	fmt.Println("Usage: gonano filename.txt")
	fmt.Println("Commands:")
	fmt.Println("Save:\t\tctrl + w")
	fmt.Println("Exit:\t\tctrl + q")
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Error: Please pass filename as argument")
		os.Exit(1)
	}

	arg := os.Args[1]

	if arg == "--help" {
		printHelp()
		os.Exit(0)
	}

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
