package main

import (
	gc "github.com/rthornton128/goncurses"
)

// MenuWindowWidth defined the width of the menu window in characters
const MenuWindowWidth = 55

// MenuWindowHeight defined the height of the menu window in characters
const MenuWindowHeight = 10

const menuTitle = "Main Menu"

// MenuWindow interface for interaction with GameMenu type
type MenuWindow interface {
	Free()
	HandleInput() bool
	init()
}

type menuItemHandlerMap = map[string]MenuItemHandlerFunction

// GameMenu contains all of the ncurses main-menu realted stuff
type GameMenu struct {
	window             *gc.Window
	menu               *gc.Menu
	menuItems          []*gc.MenuItem
	optionsHandlersMap menuItemHandlerMap
}

// MenuItemContent describes the title and description of the menu item
type MenuItemContent struct {
	MenuItemTitle       string
	MenuItemDescription string
}

// MenuItemHandlerFunction represents an action point on the particular menu item
type MenuItemHandlerFunction func() bool

// Free removes the menu, clear it from the screen and free the resources
func (m *GameMenu) Free() {
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
func (m *GameMenu) HandleInput() bool {
	m.menu.Post()
	gc.Update()
	ch := m.window.GetChar()

	switch ch {
	case gc.KEY_DOWN:
		m.menu.Driver(gc.REQ_DOWN)
	case gc.KEY_UP:
		m.menu.Driver(gc.REQ_UP)
	case gc.KEY_RETURN:
		current := m.menu.Current(nil).Name()
		return m.optionsHandlersMap[current]()
	default:
		break
	}
	m.window.Refresh()
	return true
}

// NewMenu creates new instance of main menu nested in specified Window with specified option items
func NewMenu(stdscr *gc.Window, items *[]MenuItemContent, handlers *menuItemHandlerMap) *GameMenu {
	menu := new(GameMenu)
	menu.init(stdscr, *items)
	menu.optionsHandlersMap = *handlers
	return menu
}

func (m *GameMenu) init(stdscr *gc.Window, options []MenuItemContent) {
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)

	m.menuItems = make([]*gc.MenuItem, len(options))

	maxY, maxX := stdscr.MaxYX()

	for index, item := range options {
		m.menuItems[index], _ = gc.NewItem(item.MenuItemTitle, item.MenuItemDescription)
	}

	menu, _ := gc.NewMenu(m.menuItems)

	// Centered relative to game-window
	menuWindow, _ := gc.NewWindow(MenuWindowHeight, MenuWindowWidth, maxY/2-5, maxX/2-30)
	menuWindow.Keypad(true)

	menu.SetWindow(menuWindow)
	derWin := menuWindow.Derived(6, 52, 3, 1)
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
