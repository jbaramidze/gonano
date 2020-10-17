package main

import (
	"fmt"
	"log"
	"os"
)

func initLog() {
	f, err := os.Create("log")
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(-1)
	}

	log.SetOutput(f)
}

func main() {
	initLog()

	log.Printf("ABCD")

	display := createDisplay()
	defer display.Close()

	display.poll()
}
