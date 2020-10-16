package main

import "github.com/jbaramidze/term_collab_editor/display"

func main() {
	cm := display.Initialize()
	defer cm.Close()

	cm.Poll()
}
