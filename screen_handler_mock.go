package main

type mockScreenHandler struct {
	data     [][]rune
	unsynced [][]rune
	keyChan  chan interface{}
	w        int
	h        int
}

func (s *mockScreenHandler) close() {
}

func (s *mockScreenHandler) putStr(x, y int, b rune) {
	s.unsynced[y][x] = b
}

func (s *mockScreenHandler) sync() {
	for i := range s.unsynced {
		copy(s.data[i], s.unsynced[i])
	}
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
	unsynced := make([][]rune, h)
	c := make(chan interface{})

	for i := range data {
		data[i] = make([]rune, w)
		unsynced[i] = make([]rune, w)
		for j := range data[i] {
			data[i][j] = '@'
			unsynced[i][j] = '@'
		}
	}
	return &mockScreenHandler{data: data, unsynced: unsynced, keyChan: c, w: w, h: h}
}
