package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	initLog()
}

func initLog() {
	f, err := os.Create("log")
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(-1)
	}

	log.SetOutput(f)
}
