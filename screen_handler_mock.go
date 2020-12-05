package main

type mockScreenHandler struct {
	data    [][]rune
	keyChan chan interface{}
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
	return 4, 6
}

func (s *mockScreenHandler) pollKeyPress() interface{} {
	return <-s.keyChan
}

func initMockScreenHandler() screenHandler {
	data := make([][]rune, 6)
	c := make(chan interface{})

	for i := range data {
		data[i] = make([]rune, 4)
		for j := range data[i] {
			data[i][j] = '@'
		}
	}
	return &mockScreenHandler{data: data, keyChan: c}
}