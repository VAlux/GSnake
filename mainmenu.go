package main

import gc "github.com/rthornton128/goncurses"

// CreateMenuWindow creates main menu window
func CreateMenuWindow(stdscr *gc.Window) {
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)

	menuOptions := []string{"New Game", "High Scores", "About", "Exit"}
	menuItems := make([]*gc.MenuItem, len(menuOptions))

	for index, item := range menuOptions {
		menuItems[index], _ = gc.NewItem(item, "")
		defer menuItems[index].Free()
	}

	menu, _ := gc.NewMenu(menuItems)
	defer menu.Free()

	menuWindow, _ := gc.NewWindow(10, 40, 4, 14)
	menuWindow.Keypad(true)

	menu.SetWindow(menuWindow)
	derWin := menuWindow.Derived(6, 38, 3, 1)
	menu.SubWindow(derWin)
	menu.Mark(" => ")

	_, x := menuWindow.MaxYX()
	title := "Main Menu"

	menuWindow.Box(0, 0)
	menuWindow.ColorOn(1)
	menuWindow.MovePrint(1, (x/2)-(len(title)/2), title)
	menuWindow.ColorOff(1)
	menuWindow.MoveAddChar(2, 0, gc.ACS_LTEE)
	menuWindow.HLine(2, 1, gc.ACS_HLINE, x-2)
	menuWindow.MoveAddChar(2, x-1, gc.ACS_RTEE)

	menu.Post()
	defer menu.UnPost()
	menuWindow.Refresh()

	for {
		gc.Update()
		ch := menuWindow.GetChar()

		switch ch {
		case 'q':
			return
		case gc.KEY_DOWN:
			menu.Driver(gc.REQ_DOWN)
		case gc.KEY_UP:
			menu.Driver(gc.REQ_UP)
		}
	}
}
