package main

type blinker interface {
	refresh()
	set()
	clear()
}
