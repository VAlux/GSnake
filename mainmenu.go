package main

import (
	"fmt"
	"log"

	gc "github.com/rthornton128/goncurses"
)

const (
	// MenuWindowWidth represents the width of the menu window in characters
	MenuWindowWidth = 55
	// MenuWindowHeight represents the height of the menu window in characters
	MenuWindowHeight = 10

	menuTitle            = "Main Menu"
	menuMark             = " => "
	menuMarkEmpty        = "    "
	menuItemOffset       = 5
	menuContentTopOffset = 3
)

// Menu is an interface for interaction with Menu type
type Menu interface {
	HandleInput() bool
	Free()
	Refresh()
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

func (item *MenuItem) String() string {
	return item.MenuItemTitle + "\t" + item.MenuItemDescription
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

	m.Refresh()
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
	m.window = createMenuWindow(items, maxX, maxY)
	m.window.Refresh()
}

// Refresh performs redrawing of the menu window contents
func (m *MenuWindow) Refresh() {
	for idx, item := range m.items {
		if idx == m.currentItemIndex {
			m.window.MovePrint(idx+menuContentTopOffset, 1, menuMark)
		} else {
			m.window.MovePrint(idx+menuContentTopOffset, 1, menuMarkEmpty)
		}
		m.window.MovePrint(idx+menuContentTopOffset, menuItemOffset, item.String())
	}
	m.window.Refresh()
}

// Free erase the content of the window from the screen and frees the memory, allocated for it.
func (m *MenuWindow) Free() {
	m.window.Erase()
	m.window.Delete()
}

func createMenuWindow(items []*MenuItem, maxX int, maxY int) *gc.Window {
	wnd, windowCreateError := gc.NewWindow(MenuWindowHeight, MenuWindowWidth, maxY/2-5, maxX/2-30)
	if windowCreateError != nil {
		log.Panic(fmt.Sprintf("Error creating main menu window: %s", windowCreateError))
	}

	wnd.Keypad(true)
	wnd.Box(0, 0)
	wnd.ColorOn(1)
	wnd.MovePrint(1, (MenuWindowWidth/2)-(len(menuTitle)/2), menuTitle)
	wnd.ColorOff(1)
	for idx, item := range items {
		wnd.MovePrint(idx+menuContentTopOffset, menuItemOffset, item.String())
	}
	wnd.MoveAddChar(2, 0, gc.ACS_LTEE)
	wnd.HLine(2, 1, gc.ACS_HLINE, MenuWindowWidth-2)
	wnd.MoveAddChar(2, MenuWindowWidth-1, gc.ACS_RTEE)
	return wnd
}

// NewMenu creates new instance of main menu nested in specified Window with specified option items
func NewMenu(stdscr *gc.Window, items []*MenuItem) Menu {
	menu := new(MenuWindow)
	menu.init(stdscr, items)
	return menu
}
