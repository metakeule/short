package main

import (
	"github.com/gdamore/tcell"
)

type ModalWindow struct {
	s *Screen
}

func NewModalWindow(s *Screen) *ModalWindow {
	return &ModalWindow{s: s}
}

func (m *ModalWindow) KeyF1(ev *tcell.EventKey) (quit bool) {
	return
}

func (f *ModalWindow) KeyF2(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyF4(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyCtrlE(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyCtrlP(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyCtrlS(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyCtrlSpace(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyCtrlF(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyCtrlC(ev *tcell.EventKey) (quit bool) {
	return m.s.doQuit(ev)
}

func (m *ModalWindow) KeyCtrlL(ev *tcell.EventKey) (quit bool) {
	return m.s.doSync(ev)
}

func (m *ModalWindow) KeyCtrlR(ev *tcell.EventKey) (quit bool) {
	return m.s.doSync(ev)
}

func (m *ModalWindow) KeyDEL(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyUp(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyRight(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyLeft(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyDown(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyTab(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyBackTab(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyEscape(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyEnter(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) KeyBackspace(ev *tcell.EventKey) (quit bool) {
	return
}

func (m *ModalWindow) Print() {
}

func (m *ModalWindow) KeyOther(ev *tcell.EventKey) (quit bool) {
	return
}
