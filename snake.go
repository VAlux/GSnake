package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	ll "Snake/linkedlist"

	gc "github.com/rthornton128/goncurses"

	mm "Snake/mainmenu"
)

//======================= texture :) definitions =======================

const headTexture = `#`
const tailTexture = `o`
const foodTexture = `*`
const emptyTexture = ` `

//======================= event definitions =======================

const collisionEvent = "collision"
const exitEvent = "exit"
const foodEatenEvent = "foodEaten"
const newGameEvent = "newGame"
const optionsEvent = "options"
const highScoreEvent = "highScore"
const aboutEvent = "about"

//======================= direction definitions =======================

var up = &point{-1, 0}
var down = &point{1, 0}
var left = &point{0, -1}
var right = &point{0, 1}
var nowhere = &point{0, 0}

//======================= object definitions =======================

var objects = make([]object, 0)
var events = make(chan string, 1)
var currentFood = &food{}
var playerSnake = &snake{}

//======================= window definitions =======================

var maxX = 0
var maxY = 0
var statsX = 0
var statsY = 0
var statsH = 0
var statsW = 0

//========================= Gameplay definitions =========================

//once it is false - game is over :(
var isRunning = true

//Main menu is shown during isPaused = true
var isPaused = true

var score = 0

const scorePointValue = 6
const speedFactor = 8
const initialLength = 4

//======================= Main menu definitions =======================

var mainMenu = &mm.MainMenu{}

const continueMenuItemTitle = "Continue"
const newGameMenuItemTitle = "New Game"
const optionsMenuItemTitle = "Options"
const highScoreMenuItemTitle = "High Score"
const aboutMenuItemTitle = "About"
const exitMenuItemTitle = "Exit"

const continueMenuItemDescription = " -- Resume current game"
const newGameMenuItemDescription = " -- Begin new game"
const optionsMenuItemDescription = " -- Review or change game settings"
const highScoreMenuItemDescription = " -- See the leadership table"
const aboutMenuItemDescription = " -- Info about creator"
const exitMenuItemDescription = " -- Save score and close the game"

var menuOptionsKeySet = &[]mm.MenuItemContent{
	mm.MenuItemContent{MenuItemTitle: continueMenuItemTitle, MenuItemDescription: continueMenuItemDescription},
	mm.MenuItemContent{MenuItemTitle: newGameMenuItemTitle, MenuItemDescription: newGameMenuItemDescription},
	mm.MenuItemContent{MenuItemTitle: optionsMenuItemTitle, MenuItemDescription: optionsMenuItemDescription},
	mm.MenuItemContent{MenuItemTitle: highScoreMenuItemTitle, MenuItemDescription: highScoreMenuItemDescription},
	mm.MenuItemContent{MenuItemTitle: aboutMenuItemTitle, MenuItemDescription: aboutMenuItemDescription},
	mm.MenuItemContent{MenuItemTitle: exitMenuItemTitle, MenuItemDescription: exitMenuItemDescription}}

var menuOptionsHandlerMap = &map[string]mm.MenuItemHandlerFunction{
	continueMenuItemTitle:  continueOptionHandler,
	newGameMenuItemTitle:   newGameOptionHandler,
	optionsMenuItemTitle:   optionsOptionHandler,
	highScoreMenuItemTitle: highScoreOptionHandler,
	aboutMenuItemTitle:     aboutOptionHandler,
	exitMenuItemTitle:      exitOptionHandler}

//======================= Types =======================
type point struct {
	y, x int
}

type object interface {
	update(*gc.Window)
	draw(*gc.Window)
}

type snake struct {
	head      *ll.Node
	body      *ll.LinkedList
	direction *point
}

type food struct {
	position *point
	color    int
}

//=====================================================

func (p point) String() string {
	return fmt.Sprintf("y: %d, x: %d", p.y, p.x)
}

func (p point) offset(dy, dx int) {
	p.y += dy
	p.x += dx
}

func (p point) offsetP(off *point) {
	p.y += off.y
	p.x += off.x
}

func (s *snake) update(w *gc.Window) {
	offset := nowhere
	switch s.direction {
	case up:
		offset = up
		break
	case down:
		offset = down
		break
	case left:
		offset = left
		break
	case right:
		offset = right
		break
	default:
		return
	}

	dy := s.head.Data.(point).y + offset.y
	dx := s.head.Data.(point).x + offset.x
	newHead := &ll.Node{Data: point{dy, dx}}

	if s.checkCollision(newHead) {
		events <- collisionEvent
	}

	if s.checkFoodCollision(newHead) {
		s.body.Prepend(&ll.Node{Data: point{dy, dx}})
		newHead.Data.(point).offsetP(offset)
		*currentFood = *generateFood(s)
		events <- foodEatenEvent
	}

	last := s.body.Back()
	w.MovePrint(last.Data.(point).y, last.Data.(point).x, emptyTexture)
	s.body.RemoveLast()
	s.body.Prepend(newHead)
	s.head = newHead
}

func (s *snake) draw(w *gc.Window) {
	w.ColorOn(2)
	w.AttrOn(gc.A_BOLD)
	w.MovePrint(s.head.Data.(point).y, s.head.Data.(point).x, headTexture)
	for node := s.head.Next; node.Next != nil; node = node.Next {
		w.MovePrint(node.Data.(point).y, node.Data.(point).x, tailTexture)
	}
	w.AttrOff(gc.A_BOLD)
	w.ColorOff(2)
}

func (s *snake) checkCollision(n *ll.Node) bool {
	return n.Data.(point).x <= 0 ||
		n.Data.(point).y <= 0 ||
		n.Data.(point).x >= maxX-1 ||
		n.Data.(point).y >= maxY-statsH-1 ||
		s.body.Contains(n)
}

func (s *snake) checkFoodCollision(n *ll.Node) bool {
	return n.Data.(point) == *currentFood.position
}

func (s *snake) containsNodeWithPoint(pt *point) bool {
	for node := s.head; node != nil; node = node.Next {
		if node.Data.(point) == *pt {
			return true
		}
	}
	return false
}

func (f *food) update(w *gc.Window) {
	// TODO: update the food color and/or animation
}

func (f *food) draw(w *gc.Window) {
	w.ColorOn(1)
	w.AttrOn(gc.A_BOLD)
	w.MovePrint(f.position.y, f.position.x, foodTexture)
	w.AttrOff(gc.A_BOLD)
	w.ColorOff(1)
}

func drawObjects(s *gc.Window) {
	for _, obj := range objects {
		obj.draw(s)
	}
}

func updateObjects(w *gc.Window) {
	for _, obj := range objects {
		obj.update(w)
	}
}

func tick(w *gc.Window) {
	updateObjects(w)
	drawObjects(w)
	w.Refresh()
}

func handleInput(w *gc.Window, s *snake) {
	key := w.GetChar()
	switch byte(key) {
	case 'w':
		if s.direction != down {
			s.direction = up
		}
		break
	case 's':
		if s.direction != up {
			s.direction = down
		}
		break
	case 'a':
		if s.direction != right {
			s.direction = left
		}
		break
	case 'd':
		if s.direction != left {
			s.direction = right
		}
		break
	case 'p':
		isPaused = !isPaused
		if isPaused {
			mainMenu = createMainMenu(w)
		}
		break
	case 'q':
		events <- exitEvent
		break
	default:
		break
	}
}

func createMainMenu(w *gc.Window) *mm.MainMenu {
	return mm.New(w, menuOptionsKeySet, menuOptionsHandlerMap)
}

func gameOver(s *gc.Window) {
	lines, cols := s.MaxYX()
	msg := "Game Over"

	wnd := createWindow(5, len(msg)+4, (lines/2)-2, (cols-len(msg))/2)
	wnd.MovePrint(2, 2, msg)
	wnd.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	wnd.Refresh()
	gc.Nap(2000)
}

func drawStats(sn *snake) {
	snakeLength := "length: " + strconv.Itoa(sn.body.Size())
	scoredPoints := "score: " + strconv.Itoa(score)

	wnd := createWindow(statsH, statsW-2, statsY, statsX)
	wnd.ColorOn(3)
	wnd.AttrOn(gc.A_BOLD)
	wnd.MovePrint(1, 1, snakeLength)
	wnd.MovePrint(1, len(snakeLength)+3, scoredPoints)
	wnd.ColorOff(3)
	wnd.AttrOff(gc.A_BOLD)
	wnd.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	wnd.Refresh()
}

func createSnake(y, x int) *snake {
	head := &ll.Node{Data: point{y, x}}
	body := ll.New()

	for i := 1; i <= initialLength; i++ {
		body.Append(&ll.Node{Data: point{y, x + i}})
	}

	body.Prepend(head)

	newSnake := &snake{head, body, left}
	return newSnake
}

func generateFood(sn *snake) *food {
	randX := 1 + rand.Intn(maxX-3)
	randY := 1 + rand.Intn(maxY-statsH-3)
	foodPos := &point{y: randY, x: randX}
	if sn.containsNodeWithPoint(foodPos) {
		generateFood(sn)
	}
	return &food{position: foodPos}
}

func createWindow(height, width, y, x int) *gc.Window {
	wnd, err := gc.NewWindow(height, width, y, x)
	if err != nil {
		log.Fatal("Error during creating the window...")
	}
	return wnd
}

func createGameWindow(y, x, height, width int) *gc.Window {
	wnd := createWindow(height, width, y, x)
	wnd.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	wnd.Refresh()
	return wnd
}

func openLogFile() *os.File {
	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	return logFile
}

func newGame(w *gc.Window, headY int, headX int) {
	log.Print("Starting new game...")
	playerSnake = createSnake(headY, headX)
	currentFood = generateFood(playerSnake)
	objects = make([]object, 0)
	objects = append(objects, playerSnake, currentFood)
	score = 0
	w.Erase()
	w.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	w.Refresh()
}

func handleEvents(s *snake, w *gc.Window) {
	select {
	case event := <-events:
		log.Printf("Event occurred: %s", event)
		if event == foodEatenEvent {
			score += scorePointValue*speedFactor + s.body.Size()
			log.Printf("Score increased. Current score: %d", score)
		}
		if event == collisionEvent {
			isRunning = false // exit
			break
		}
		if event == exitEvent {
			isRunning = false // exit
			break
		}
		if event == newGameEvent {
			newGame(w, maxY/2, maxX/2)
			break
		}
	default:
		break
	}
}

//======================= Main Menu Handlers =======================

func continueOptionHandler() bool {
	log.Print("Continue menu option selected")
	return false
}

func newGameOptionHandler() bool {
	log.Print("New Game menu option selected")
	events <- newGameEvent
	return false
}

func optionsOptionHandler() bool {
	log.Print("Options menu option selected")
	events <- optionsEvent
	return true
}

func highScoreOptionHandler() bool {
	log.Print("High Score menu option selected")
	events <- highScoreEvent
	return true
}

func aboutOptionHandler() bool {
	log.Print("About menu option selected")
	events <- aboutEvent
	return true
}

func exitOptionHandler() bool {
	log.Print("Exit menu option selected")
	events <- exitEvent
	return false
}

//======================= Initialization =======================

func initNcurses() {
	// Coloring setup
	gc.StartColor()
	gc.InitPair(1, gc.C_RED, gc.C_BLACK)
	gc.InitPair(2, gc.C_GREEN, gc.C_BLACK)
	gc.InitPair(3, gc.C_YELLOW, gc.C_BLACK)

	gc.Cursor(0)
	gc.Echo(false)
	gc.Raw(true)
	gc.HalfDelay(1)
}

func initScreenDimensions(stdscr *gc.Window) error {
	maxY, maxX = stdscr.MaxYX()
	statsX, statsY, statsH, statsW = 1, 0, 3, maxX
	log.Printf("Resolution: %d x %d", maxX, maxY)
	// Check the resolution and exit if the terminal window is too small
	if maxY < mm.MenuWindowHeight+5 || maxX < mm.MenuWindowWidth+5 {
		log.Print("Recommended resolution is 60x25")
		return errors.New("Too small game window. Program will exit")
	}
	return nil
}

func initLogging() *os.File {
	logFile := openLogFile()
	log.SetOutput(logFile)
	return logFile
}

// ==================================================================

func main() {
	stdscr, err := gc.Init()

	if err != nil {
		log.Panicln("Error during ncurses Init: ", err)
	}

	rand.Seed(int64(time.Now().Second()))
	logFile := initLogging()

	// Finalization
	defer logFile.Close()
	defer gc.End()
	defer gameOver(stdscr)
	defer log.Println(" <==== Game session ended")
	//

	log.Println(" ====> Game session started")
	initNcurses()

	dimensionsInitError := initScreenDimensions(stdscr)
	if dimensionsInitError != nil {
		log.Panicln("Error initializing the screen dimensions:", dimensionsInitError)
		return
	}

	ticker := time.NewTicker(time.Second / speedFactor)

	// Create in-game windows
	gameWindow := createGameWindow(statsY+statsH, statsX, maxY-statsH, statsW-2)
	mainMenu = createMainMenu(gameWindow)
	//

	newGame(gameWindow, maxY/2, maxX/2)

	//Game Loop:
	for isRunning {
		select {
		case <-ticker.C:
			if !isPaused {
				handleInput(gameWindow, playerSnake)
				tick(gameWindow)
				drawStats(playerSnake)
				handleEvents(playerSnake, gameWindow)
			} else {
				if !mainMenu.HandleInput() {
					isPaused = false
					mainMenu.Free()
				}
			}
		}
	}
}
