package main

import (
	"log"

	gc "github.com/rthornton128/goncurses"
)

const aboutWindowHeight = 8
const aboutWindowWidth = 40
const aboutWindowTitle = "About"

var aboutText = []string{
	"This game is created by Alvo",
	"for beloved Tanya :)",
	"",
	"Have fun!"}

// CreateAboutWindow creates and shows window with the info about creator
func CreateAboutWindow(s *gc.Window) {
	log.Println("Creating about window...")

	lines, cols := s.MaxYX()
	height, width := aboutWindowHeight, aboutWindowWidth
	contentOffset := 3

	wnd, windowCreateError := createWindow(height, width, (lines/2)-height/2, (cols/2)-width/2)
	if windowCreateError != nil {
		log.Println("Error creating high score window: ", windowCreateError)
		return
	}

	wnd.ColorOn(1)
	wnd.MovePrint(
		1,
		(width/2)-(len(aboutWindowTitle)/2),
		aboutWindowTitle)
	wnd.ColorOff(1)

	wnd.Box(0, 0)
	wnd.ColorOn(3)
	for idx, aboutLine := range aboutText {
		wnd.MovePrint(idx+contentOffset, contentOffset, aboutLine)
	}
	wnd.ColorOff(3)
	wnd.MoveAddChar(2, 0, gc.ACS_LTEE)
	wnd.HLine(2, 1, gc.ACS_HLINE, width-2)
	wnd.MoveAddChar(2, width-1, gc.ACS_RTEE)
	wnd.Refresh()

	log.Println("About window created")

	awaitClosingAction(wnd)
}
