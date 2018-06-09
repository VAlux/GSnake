package main

import (
	"log"

	gc "github.com/rthornton128/goncurses"
)

const defaultPlayerName = "Anon"
const playerNameWindowTitle = "Player name"
const playerNameWindowHeight = 10
const playerNameWindowWidth = 50

// GetPlayerName create and show the window with player name input form
func GetPlayerName(s *gc.Window) string {

	lines, cols := s.MaxYX()
	height, width := playerNameWindowHeight, playerNameWindowWidth

	wnd, windowCreateError := createWindow(
		height,
		width,
		(lines/2)-height/2,
		(cols/2)-width/2)

	// we need to enable echo and cursor to be able to input something in the terminal
	gc.Echo(true)
	gc.Cursor(1)
	defer gc.Echo(false)
	defer gc.Cursor(0)
	//

	if windowCreateError != nil {
		log.Println("Error creating player name input form window: ", windowCreateError)
		return defaultPlayerName
	}

	wnd.Box(0, 0)
	wnd.ColorOn(1)
	wnd.MovePrint(
		1,
		(width/2)-(len(playerNameWindowTitle)/2),
		playerNameWindowTitle)
	wnd.ColorOff(1)

	wnd.MoveAddChar(2, 0, gc.ACS_LTEE)
	wnd.HLine(2, 1, gc.ACS_HLINE, width-2)
	wnd.MoveAddChar(2, width-1, gc.ACS_RTEE)
	wnd.Refresh()

	log.Println("High score window created")

	playerName := promptPlayerName(wnd)
	log.Println("player name is: ", playerName)
	removeWindow(wnd)
	return playerName
}

func promptPlayerName(w *gc.Window) string {
	msg := "Enter your name: "
	row, col := w.MaxYX()
	row, col = (row/2)-1, 4
	w.MovePrint(row, col, msg)

	var str string
	str, err := w.GetString(12)
	if err != nil {
		log.Panic("Error getting player name string: ", err)
		return defaultPlayerName
	}
	w.Refresh()
	return str
}
