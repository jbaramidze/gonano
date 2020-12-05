package main

type mockScreenHandler struct {
	data    [][]rune
	keyChan chan interface{}
	w       int
	h       int
}

func (s *mockScreenHandler) close() {
}

func (s *mockScreenHandler) putStr(x, y int, b rune) {
	s.data[y][x] = b
}

func (s *mockScreenHandler) clearStr(x, y int) {
	s.putStr(x, y, '@')
}

func (s *mockScreenHandler) getSize() (int, int) {
	return s.w, s.h
}

func (s *mockScreenHandler) pollKeyPress() interface{} {
	return <-s.keyChan
}

func initMockScreenHandler(w, h int) screenHandler {
	data := make([][]rune, h)
	c := make(chan interface{})

	for i := range data {
		data[i] = make([]rune, w)
		for j := range data[i] {
			data[i][j] = '@'
		}
	}
	return &mockScreenHandler{data: data, keyChan: c, w: w, h: h}
}
