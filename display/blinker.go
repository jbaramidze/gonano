package display

import "time"

// InitAndStartBlinker does
func InitAndStartBlinker(c chan ContentOperation) {
	go func(c chan ContentOperation) {
		ticker := time.NewTicker(500 * time.Millisecond)
		init := false
		for {
			init = !init
			<-ticker.C
			c <- BlinkOperation{blink: init}
		}
	}(c)
}
