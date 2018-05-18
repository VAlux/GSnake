package mainmenu

import (
	gc "github.com/rthornton128/goncurses"
)

type MainMenuWindow interface {
	SetVisible(isVisible bool)
	Free()
	eventLoop()
}

type MainMenu struct {
	window    *gc.Window
	menu      *gc.Menu
	menuItems []*gc.MenuItem
	isVisible bool
}

func (m *MainMenu) SetVisible(isVisible bool) {
	m.isVisible = isVisible
	if isVisible {
		go m.eventLoop()
	}
}

func (m *MainMenu) Free() {
	for _, item := range m.menuItems {
		item.Free()
	}
	m.menu.Free()
}

func (m *MainMenu) eventLoop() {
	m.menu.Post()
	defer m.menu.UnPost()
	m.window.Refresh()

	for m.isVisible {
		gc.Update()
		ch := m.window.GetChar()

		switch ch {
		case 'q':
			return
		case gc.KEY_DOWN:
			m.menu.Driver(gc.REQ_DOWN)
		case gc.KEY_UP:
			m.menu.Driver(gc.REQ_UP)
		}
	}
}

func New(stdscr *gc.Window, options []string) *MainMenu {
	menu := new(MainMenu)
	menu.init(stdscr, options)
	return menu
}

// Init creates main menu window
func (m *MainMenu) init(stdscr *gc.Window, options []string) {
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)

	m.menuItems = make([]*gc.MenuItem, len(options))

	for index, item := range options {
		m.menuItems[index], _ = gc.NewItem(item, "")
	}

	menu, _ := gc.NewMenu(m.menuItems)

	menuWindow, _ := gc.NewWindow(10, 40, 4, 14)
	menuWindow.Keypad(true)

	menu.SetWindow(menuWindow)
	derWin := menuWindow.Derived(6, 38, 3, 1)
	menu.SubWindow(derWin)
	menu.Mark(" => ")

	m.menu = menu

	_, x := menuWindow.MaxYX()
	title := "Main Menu"

	menuWindow.Box(0, 0)
	menuWindow.ColorOn(1)
	menuWindow.MovePrint(1, (x/2)-(len(title)/2), title)
	menuWindow.ColorOff(1)
	menuWindow.MoveAddChar(2, 0, gc.ACS_LTEE)
	menuWindow.HLine(2, 1, gc.ACS_HLINE, x-2)
	menuWindow.MoveAddChar(2, x-1, gc.ACS_RTEE)

	m.window = menuWindow
}
