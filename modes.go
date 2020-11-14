package main

import (
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

type mode interface {
	handleKeyPress(e event, resp chan bool) (exit bool)
	init()
}

type normalMode struct {
	e *Editor
}

func newNormalMode(e *Editor) normalMode {
	return normalMode{e}
}

func (m normalMode) init() {
	tmp := m.e.display.currentElement
	m.e.display.currentElement = m.e.display.data.Front()
	m.e.display.resyncBelowCurrent()

	// Leave it at the end
	m.e.display.currentElement = tmp
}

func (m normalMode) handleKeyPress(ev event, resp chan bool) (exit bool) {
	if ev.k == tcell.KeyCtrlQ {
		if m.e.modified == false {
			// Exit the editor
			return true
		}
		m.e.setMode(newQuitWithoutSavingMode(m.e))
	} else if ev.k == tcell.KeyCtrlW {
		m.e.setMode(newSavedMode(m.e))
	} else {
		if ev.k != tcell.KeyLeft && ev.k != tcell.KeyRight && ev.k != tcell.KeyUp && ev.k != tcell.KeyDown {
			m.e.modified = true
		}
		m.e.display.monitorChannel <- TypeOperation{rn: ev.rn, key: ev.k, resp: resp}
	}

	return false
}

type quitWithoutSavingMode struct {
	e *Editor
}

func newQuitWithoutSavingMode(e *Editor) quitWithoutSavingMode {
	return quitWithoutSavingMode{e}
}

func (m quitWithoutSavingMode) init() {
	m.e.display.monitorChannel <- AnnouncementOperation{text: []string{"You will lose your changes!", "Are you sure you want to quit? Y/N"}}
}

func (m quitWithoutSavingMode) handleKeyPress(ev event, resp chan bool) (exit bool) {
	if ev.rn == rune('y') {
		return true
	} else if ev.rn == rune('n') {
		m.e.setMode(normalMode{m.e})
	}

	return false
}

type savedMode struct {
	e    *Editor
	lock *sync.Mutex
}

func newSavedMode(e *Editor) savedMode {
	return savedMode{e: e, lock: new(sync.Mutex)}
}

func (m savedMode) init() {
	if err := m.e.saveData(); err != nil {
		// Show somewhere?
		return
	}
	m.e.modified = false
	m.e.display.monitorChannel <- AnnouncementOperation{text: []string{"Saved!"}}
	go func() {
		time.Sleep(3 * time.Second)
		m.lock.Lock()
		defer m.lock.Unlock()
		if _, ok := m.e.mode.(savedMode); ok == true {
			m.e.setMode(normalMode{m.e})
		}
	}()
}

func (m savedMode) handleKeyPress(ev event, resp chan bool) (exit bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.e.setMode(normalMode{m.e})
	return false
}
