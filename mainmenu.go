package main

import (
	"fmt"
	"log"

	gc "github.com/rthornton128/goncurses"
)

const (
	menuWindowWidth  = 55
	menuWindowHeight = 10
	menuTitle        = "Main Menu"
)

// Menu is an interface for interaction with Menu type
type Menu interface {
	HandleInput() bool
	init(stdscr *gc.Window, items []*MenuItem)
}

// MenuItemHandlerFunction represents an action point on the particular menu item
type MenuItemHandlerFunction func() bool

// MenuWindow  contains all of the ncurses main-menu realted stuff
type MenuWindow struct {
	window           *gc.Window
	items            []*MenuItem
	currentItemIndex int
}

// MenuItem describes the title description and functionality of the menu item
type MenuItem struct {
	MenuItemTitle       string
	MenuItemDescription string
	MenuItemHandler     MenuItemHandlerFunction
}

// NewMenuItem creates new menu item with specified title, description and handler
func NewMenuItem(title string, description string, handler MenuItemHandlerFunction) *MenuItem {
	return &MenuItem{
		MenuItemTitle:       title,
		MenuItemDescription: description,
		MenuItemHandler:     handler,
	}
}

// HandleInput obtains the user input and executes actions based on it.
func (m *MenuWindow) HandleInput() bool {
	gc.Update()
	ch := m.window.GetChar()

	switch ch {
	case gc.KEY_DOWN:
		m.moveCaretDown()
	case gc.KEY_UP:
		m.moveCaretUp()
	case gc.KEY_RETURN:
		return m.executeCurrentHandler()
	default:
		break
	}

	m.window.Refresh()
	return true
}

func (m *MenuWindow) moveCaretDown() {
	if m.currentItemIndex == len(m.items)-1 {
		m.currentItemIndex = 0
	} else {
		m.currentItemIndex++
	}
}

func (m *MenuWindow) moveCaretUp() {
	if m.currentItemIndex == 0 {
		m.currentItemIndex = len(m.items) - 1
	} else {
		m.currentItemIndex--
	}
}

func (m *MenuWindow) getCurrentItem() *MenuItem {
	return m.items[m.currentItemIndex]
}

func (m *MenuWindow) executeCurrentHandler() bool {
	return m.getCurrentItem().MenuItemHandler()
}

func (m *MenuWindow) init(stdscr *gc.Window, items []*MenuItem) {
	maxY, maxX := stdscr.MaxYX()
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	m.currentItemIndex = 0
	m.items = items
	m.window = createMenuWindow(stdscr, items, maxX)
	m.window.Refresh()
}

func createMenuWindow(stdscr *gc.Window, items []*MenuItem, x int) *gc.Window {
	wnd, windowCreateError := gc.NewWindow(menuWindowHeight, menuWindowWidth, maxY/2-5, maxX/2-30)
	if windowCreateError != nil {
		log.Panic(fmt.Sprintf("Error creating main menu window: %s", windowCreateError))
	}

	wnd.Keypad(true)
	wnd.Box(0, 0)
	wnd.ColorOn(1)
	wnd.MovePrint(1, (x/2)-(len(menuTitle)/2), menuTitle)
	wnd.ColorOff(1)
	for idx, item := range items {
		wnd.MovePrint(idx+2, 1, item.MenuItemTitle+item.MenuItemDescription)
	}
	wnd.MoveAddChar(2, 0, gc.ACS_LTEE)
	wnd.HLine(2, 1, gc.ACS_HLINE, x-2)
	wnd.MoveAddChar(2, x-1, gc.ACS_RTEE)
	return wnd
}

// NewMenu creates new instance of main menu nested in specified Window with specified option items
func NewMenu(stdscr *gc.Window, items []*MenuItem) Menu {
	menu := new(MenuWindow)
	menu.init(stdscr, items)
	return menu
}
