package main

func main() {

	// arg := os.Args[1]

	// data, err := ioutil.ReadFile(arg)
	// if err != nil {
	// 	fmt.Printf("Opening file failed: %v", err)
	// 	return
	// }
	//
	//fmt.Printf("%v", data)

	handler := initPhysicalScreenHandler()
	editor := createEditor(handler)

	blinkr := initRealBlinker(editor)
	editor.setBlinker(blinkr)

	go editor.startLoop()
	defer editor.display.Close()

	editor.pollKeyboard(nil)
}
