package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Error: Please pass filename as argument")
		os.Exit(1)
	}

	arg := os.Args[1]

	data, err := ioutil.ReadFile(arg)
	if err != nil {
		fmt.Printf("Error: Opening file failed: %v", err)
		return
	}

	handler := initPhysicalScreenHandler()
	editor := createEditor(handler)
	editor.initData(data)

	blinkr := initRealBlinker(editor)
	editor.setBlinker(blinkr)

	go editor.startLoop()
	defer editor.display.Close()

	editor.pollKeyboard(nil)
}
