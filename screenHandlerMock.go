package main

type mockScreenHandler struct {
	data    [][]rune
	keyChan chan event
}

func (s *mockScreenHandler) close() {
}

func (s *mockScreenHandler) putStr(x, y int, b rune) {
	s.data[y][x] = b
}

func (s *mockScreenHandler) getSize() (int, int) {
	return 4, 4
}

func (s *mockScreenHandler) pollKeyPress() event {
	return <-s.keyChan
}

func initMockScreenHandler() screenHandler {
	data := make([][]rune, 6)
	c := make(chan event)

	for i := range data {
		data[i] = make([]rune, 4)
	}
	return &mockScreenHandler{data: data, keyChan: c}
}
