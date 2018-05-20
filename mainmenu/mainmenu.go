package mainmenu

import (
	gc "github.com/rthornton128/goncurses"
)

const menuTitle = "Main Menu"
const menuWindowWidth = 40
const menuWindowHeight = 10

// MenuWindow interface for interaction with MainMenu type
type MenuWindow interface {
	Free()
	HandleInput() bool
	init()
}

// MainMenu contains all of the ncurses main-menu realted stuff
type MainMenu struct {
	window    *gc.Window
	menu      *gc.Menu
	menuItems []*gc.MenuItem
}

// Free removes the menu, clear it from the screen and free the resources
func (m *MainMenu) Free() {
	m.menu.UnPost()
	for _, item := range m.menuItems {
		item.Free()
	}
	m.menu.Free()
	m.window.Erase()
	m.window.Refresh()
	m.window.Delete()
}

// HandleInput contains all of the menu window input action handling
func (m *MainMenu) HandleInput() bool {
	m.menu.Post()
	gc.Update()
	ch := m.window.GetChar()

	switch ch {
	case 'q':
		return false
	case gc.KEY_DOWN:
		m.menu.Driver(gc.REQ_DOWN)
	case gc.KEY_UP:
		m.menu.Driver(gc.REQ_UP)
	}
	m.window.Refresh()
	return true
}

// New creates new instance of main menu nested in specified Window with specified option items
func New(stdscr *gc.Window, options *[]string) *MainMenu {
	menu := new(MainMenu)
	menu.init(stdscr, *options)
	return menu
}

func (m *MainMenu) init(stdscr *gc.Window, options []string) {
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)

	m.menuItems = make([]*gc.MenuItem, len(options))

	maxY, maxX := stdscr.MaxYX()

	for index, item := range options {
		m.menuItems[index], _ = gc.NewItem(item, "")
	}

	menu, _ := gc.NewMenu(m.menuItems)

	// Centrized relative to game-window
	menuWindow, _ := gc.NewWindow(menuWindowHeight, menuWindowWidth, maxY/2-5, maxX/2-20)
	menuWindow.Keypad(true)

	menu.SetWindow(menuWindow)
	derWin := menuWindow.Derived(6, 38, 3, 1)
	menu.SubWindow(derWin)
	menu.Mark(" => ")

	m.menu = menu

	_, x := menuWindow.MaxYX()

	menuWindow.Box(0, 0)
	menuWindow.ColorOn(1)
	menuWindow.MovePrint(1, (x/2)-(len(menuTitle)/2), menuTitle)
	menuWindow.ColorOff(1)
	menuWindow.MoveAddChar(2, 0, gc.ACS_LTEE)
	menuWindow.HLine(2, 1, gc.ACS_HLINE, x-2)
	menuWindow.MoveAddChar(2, x-1, gc.ACS_RTEE)

	m.window = menuWindow
}
