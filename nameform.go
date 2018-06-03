package main

import (
	"log"

	gc "github.com/rthornton128/goncurses"
)

const defaultPlayerName = "Anon"
const playerNameWindowTitle = "Player name"
const playerNameWindowHeight = 20
const playerNameWindowWidth = 80

// CreatePlayerNameInputFormWindow create and show the window with player name input form
func CreatePlayerNameInputFormWindow(s *gc.Window) string {
	lines, cols := s.MaxYX()
	height, width := playerNameWindowHeight, playerNameWindowWidth

	wnd, windowCreateError := createWindow(
		height,
		width,
		(lines/2)-height/2,
		(cols/2)-width/2)

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

	playerName := createPlayerNameForm(wnd)
	removeWindow(wnd)
	return playerName
}

// createPlayerNameForm create input form for entering the player name
func createPlayerNameForm(w *gc.Window) string {
	nameField, _ := gc.NewField(1, 10, 4, 18, 0, 0)
	defer nameField.Free()
	nameField.SetForeground(gc.ColorPair(1))
	nameField.SetBackground(gc.ColorPair(2) | gc.A_UNDERLINE | gc.A_BOLD)
	nameField.SetOptionsOff(gc.FO_AUTOSKIP)

	fields := make([]*gc.Field, 1)
	fields[0] = nameField

	form, _ := gc.NewForm(fields)
	form.Post()
	defer form.UnPost()
	defer form.Free()
	w.Refresh()

	w.AttrOn(gc.ColorPair(2) | gc.A_BOLD)
	w.MovePrint(3, 5, "Player Name:")
	w.AttrOff(gc.ColorPair(2) | gc.A_BOLD)
	w.Refresh()

	form.Driver(gc.REQ_FIRST_FIELD)

	ch := w.GetChar()
	for ch != 'q' {
		switch ch {
		case gc.KEY_ENTER, gc.KEY_RETURN:
			return nameField.Buffer()
		case gc.KEY_BACKSPACE:
			form.Driver(gc.REQ_CLR_FIELD)
		default:
			form.Driver(ch)
		}
		ch = w.GetChar()
	}

	return defaultPlayerName
}
