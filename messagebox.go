package main

import (
	"fmt"
	"log"

	gc "github.com/rthornton128/goncurses"
)

// MessageBox representing the window with title and content
type MessageBox struct {
	Height      int
	Width       int
	Title       string
	MessageText []string
}

// Show creates the window and shows it as a child of specified window
func (mBox *MessageBox) Show(s *gc.Window) {
	log.Println(fmt.Sprintf("Creating %s window...", mBox.Title))

	lines, cols := s.MaxYX()
	height, width := mBox.Height, mBox.Width
	contentOffset := 3

	wnd, windowCreateError := createWindow(height, width, (lines/2)-height/2, (cols/2)-width/2)
	if windowCreateError != nil {
		log.Println(fmt.Sprintf("Error creating %s window: %s", mBox.Title, windowCreateError))
		return
	}

	wnd.ColorOn(1)
	wnd.MovePrint(
		1,
		(width/2)-(len(mBox.Title)/2),
		mBox.Title)
	wnd.ColorOff(1)

	wnd.Box(0, 0)
	wnd.ColorOn(3)
	for idx, line := range mBox.MessageText {
		wnd.MovePrint(idx+contentOffset, contentOffset, line)
	}
	wnd.ColorOff(3)
	wnd.MoveAddChar(2, 0, gc.ACS_LTEE)
	wnd.HLine(2, 1, gc.ACS_HLINE, width-2)
	wnd.MoveAddChar(2, width-1, gc.ACS_RTEE)
	wnd.Refresh()

	log.Println(mBox.Title + " window created")

	awaitClosingAction(wnd)
}
