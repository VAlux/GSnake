package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	gc "github.com/rthornton128/goncurses"
)

//======================= texture :) definitions =======================

const (
	headTexture  = `#`
	tailTexture  = `o`
	emptyTexture = ` `
)

var foodAnimation = NewAnimation([]string{`-`, `\`, `|`, `/`}, 1)

//======================= event definitions =======================

const (
	collisionEvent = "collision"
	exitEvent      = "exit"
	foodEatenEvent = "foodEaten"
	newGameEvent   = "newGame"
	helpEvent      = "help"
	highScoreEvent = "highScore"
	aboutEvent     = "about"
)

//======================= direction definitions =======================

var (
	up      = &point{-1, 0}
	down    = &point{1, 0}
	left    = &point{0, -1}
	right   = &point{0, 1}
	nowhere = &point{0, 0}
)

//======================= object definitions =======================

var (
	objects     = make([]object, 0)
	events      = make(chan string, 1)
	currentFood = &food{}
	playerSnake = &snake{}
)

//======================= window definitions =======================

var (
	maxX   = 0
	maxY   = 0
	statsX = 0
	statsY = 0
	statsH = 0
	statsW = 0
)

//========================= Gameplay definitions =========================

//once it is false - game is over :(
var isRunning = true

//Main menu is shown during isPaused = true
var isPaused = false

var score = 0
var boundFactor = 0

const scorePointValue = 6
const speedFactor = 8
const initialLength = 4

//======================= Main menu definitions =======================

var menu = &MenuWindow{}

const (
	continueMenuItemTitle  = "Continue"
	newmenuItemTitle       = "New Game"
	optionsMenuItemTitle   = "Help"
	highScoreMenuItemTitle = "High Score"
	aboutMenuItemTitle     = "About"
	exitMenuItemTitle      = "Exit"
)

const (
	continueMenuItemDescription  = " -- Resume current game"
	newmenuItemDescription       = " -- Begin new game"
	optionsMenuItemDescription   = " -- See the gameplay help"
	highScoreMenuItemDescription = " -- See the leadership table"
	aboutMenuItemDescription     = " -- Info about creator"
	exitMenuItemDescription      = " -- Save score and close the game"
)

var menuOptionsKeySet = []*MenuItem{
	&MenuItem{
		MenuItemTitle:       continueMenuItemTitle,
		MenuItemDescription: continueMenuItemDescription,
		MenuItemHandler:     continueOptionHandler},

	&MenuItem{
		MenuItemTitle:       newmenuItemTitle,
		MenuItemDescription: newmenuItemDescription,
		MenuItemHandler:     newGameOptionHandler},

	&MenuItem{
		MenuItemTitle:       optionsMenuItemTitle,
		MenuItemDescription: optionsMenuItemDescription,
		MenuItemHandler:     helpOptionHandler},

	&MenuItem{
		MenuItemTitle:       highScoreMenuItemTitle,
		MenuItemDescription: highScoreMenuItemDescription,
		MenuItemHandler:     highScoreOptionHandler},

	&MenuItem{
		MenuItemTitle:       aboutMenuItemTitle,
		MenuItemDescription: aboutMenuItemDescription,
		MenuItemHandler:     aboutOptionHandler},

	&MenuItem{
		MenuItemTitle:       exitMenuItemTitle,
		MenuItemDescription: exitMenuItemDescription,
		MenuItemHandler:     exitOptionHandler}}

//======================= Types =======================
type point struct {
	y, x int
}

type object interface {
	update(*gc.Window)
	draw(*gc.Window)
}

type snake struct {
	head        *Node
	body        *LinkedList
	direction   *point
	headTexture string
	tailTexture string
}

type food struct {
	position  *point
	animation Animation
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
	dy := s.head.Data.(point).y + s.direction.y
	dx := s.head.Data.(point).x + s.direction.x
	newHead := &Node{Data: point{dy, dx}}

	if s.checkCollision(newHead) {
		events <- collisionEvent
	}

	if s.checkFoodCollision(newHead) {
		s.body.Prepend(&Node{Data: point{dy, dx}})
		newHead.Data.(point).offsetP(s.direction)
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
	w.MovePrint(s.head.Data.(point).y, s.head.Data.(point).x, s.headTexture)
	for node := s.head.Next; node.Next != nil; node = node.Next {
		w.MovePrint(node.Data.(point).y, node.Data.(point).x, s.tailTexture)
	}
	w.AttrOff(gc.A_BOLD)
	w.ColorOff(2)
}

func (s *snake) checkCollision(n *Node) bool {
	return n.Data.(point).x <= 0 ||
		n.Data.(point).y <= 0 ||
		n.Data.(point).x >= maxX-3 ||
		n.Data.(point).y >= maxY-statsH-1 ||
		s.body.Contains(n)
}

func (s *snake) checkFoodCollision(n *Node) bool {
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
	f.animation.MoveFrameIndex()
}

func (f *food) draw(w *gc.Window) {
	w.ColorOn(1)
	w.AttrOn(gc.A_BOLD)
	w.MovePrint(f.position.y, f.position.x, f.animation.CurrentFrame())
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
	key := byte(w.GetChar())
	if key == 0 {
		return
	}

	switch key {
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
	case 'p', 27: // p or escape key
		pause(w)
		break
	default:
		break
	}
}

func pause(w *gc.Window) {
	isPaused = !isPaused
	if isPaused {
		menu = createMenu(w).(*MenuWindow)
	}
}

func createMenu(w *gc.Window) Menu {
	return NewMenu(w, menuOptionsKeySet)
}

func gameOver(s *gc.Window) {
	lines, cols := s.MaxYX()
	msg := "Game Over"

	wnd, err := createWindow(5, len(msg)+4, (lines/2)-2, (cols-len(msg))/2)
	if err != nil {
		log.Panic("Error creating game over window", err)
		return
	}

	wnd.MovePrint(2, 2, msg)
	wnd.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	wnd.Refresh()
	gc.Nap(2000)
}

func drawStats(sn *snake) {
	snakeLength := "length: " + strconv.Itoa(sn.body.Size())
	scoredPoints := "score: " + strconv.Itoa(score)

	wnd, err := createWindow(statsH, statsW-2, statsY, statsX)
	if err != nil {
		log.Panic("Error creating stats window", err)
		return
	}

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
	head := &Node{Data: point{y, x}}
	body := NewList()

	for i := 1; i <= initialLength; i++ {
		body.Append(&Node{Data: point{y, x + i}})
	}

	body.Prepend(head)

	newSnake := &snake{head, body, left, headTexture, tailTexture}
	return newSnake
}

func generateFood(sn *snake) *food {
	randX := 1 + rand.Intn(maxX-4)
	randY := 1 + rand.Intn(maxY-statsH-2)
	foodPos := &point{y: randY, x: randX}
	if sn.containsNodeWithPoint(foodPos) {
		generateFood(sn)
	}
	return &food{foodPos, foodAnimation}
}

func createWindow(height, width, y, x int) (*gc.Window, error) {
	wnd, err := gc.NewWindow(height, width, y, x)
	if err != nil {
		message := "Error during creating the window: " + err.Error()
		return nil, errors.New(message)
	}
	return wnd, nil
}

func removeWindow(wnd *gc.Window) {
	wnd.Erase()
	wnd.Refresh()
	wnd.Delete()
}

func createGameWindow(y, x, height, width int) (*gc.Window, error) {
	wnd, err := createWindow(height, width, y, x)
	if err != nil {
		log.Panic("Error creating game window:", err)
		return nil, err
	}
	wnd.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	wnd.Refresh()
	return wnd, nil
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
		switch event {
		case foodEatenEvent:
			incrementScore(s)
			log.Printf("Score increased. Current score: %d", score)
			break
		case collisionEvent:
			isRunning = false
			break
		case exitEvent:
			isRunning = false
			break
		case newGameEvent:
			newGame(w, maxY/2, maxX/2)
			break
		case highScoreEvent:
			createHighScoreWindow(w)
			break
		case aboutEvent:
			createAboutWindow(w)
			break
		case helpEvent:
			createHelpWindow(w)
			break
		default:
			break
		}
	default:
		break
	}
}

func createAboutWindow(w *gc.Window) {
	const aboutWindowHeight = 8
	const aboutWindowWidth = 40
	const aboutWindowTitle = "About"

	var aboutText = []string{
		"This game is created by Alvo",
		"for beloved Tanya :)",
		"",
		"Have fun!"}

	showMessageBox(aboutWindowHeight, aboutWindowWidth, aboutWindowTitle, aboutText, w)
}

func createHelpWindow(w *gc.Window) {
	const helpWindowHeight = 8
	const helpWindowWidth = 42
	const helpWindowTitle = "Help"

	var helpText = []string{
		"Controls:",
		"'W' 'S' 'A' 'D' for direction change",
		"'P' or 'Esc' for pause/menu"}

	showMessageBox(helpWindowHeight, helpWindowWidth, helpWindowTitle, helpText, w)
}

func createHighScoreWindow(w *gc.Window) {
	const highscoreWindowTitle = "High scores"
	const highScoreWindowWidth = 70
	const highScoreWindowHeight = 14

	scores, scoreLoadError := LoadHighScore()
	if scoreLoadError != nil {
		log.Println("Error loading high scores: ", scoreLoadError)
		scores = HighScores{}
	}

	sort.Sort(scores)

	scoreContent := []string{}
	for idx, score := range scores {
		if idx >= maxAmountOfTopHighScores {
			break
		}
		scoreContent = append(scoreContent, score.String())
	}

	showMessageBox(highScoreWindowHeight, highScoreWindowWidth, highscoreWindowTitle, scoreContent, w)
}

func showMessageBox(height int, width int, title string, text []string, w *gc.Window) {
	mBox := MessageBox{
		Height:      height,
		Width:       width,
		Title:       title,
		MessageText: text}

	mBox.Show(w)
}

func incrementScore(s *snake) {
	score += (scorePointValue*speedFactor + s.body.Size()) - boundFactor
}

// saveHighScore Enter player name and save the high score if it is greater than 0
func saveHighScore(w *gc.Window) {
	if score > 0 {
		playerName := GetPlayerName(w)
		SaveHighScore(
			&HighScore{
				Timestamp:  time.Now(),
				Score:      score,
				PlayerName: playerName})
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

func helpOptionHandler() bool {
	log.Print("Help menu option selected")
	events <- helpEvent
	return false
}

func highScoreOptionHandler() bool {
	log.Print("High Score menu option selected")
	events <- highScoreEvent
	return false
}

func aboutOptionHandler() bool {
	log.Print("About menu option selected")
	events <- aboutEvent
	return false
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
	gc.CBreak(true)
	gc.HalfDelay(1)
}

func initScreenDimensions(stdscr *gc.Window) error {
	maxY, maxX = stdscr.MaxYX()
	statsX, statsY, statsH, statsW = 1, 0, 3, maxX
	log.Printf("Resolution: %d x %d", maxX, maxY)
	// Check the resolution and exit if the terminal window is too small
	if maxY < MenuWindowHeight+5 || maxX < MenuWindowWidth+5 {
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

func getScoreBoundFactor() int {
	factor := int(math.Min(float64(maxX), float64(maxY)))
	maxFactor := scorePointValue * speedFactor
	if factor >= maxFactor {
		log.Println("Screen dimensions are too big -> Adjusting score calculation formula")
		factor = maxFactor - 1
	}
	return factor
}

func initScoreBoundFactor() {
	boundFactor = getScoreBoundFactor()
}

// ==================================================================

func main() {
	stdscr, err := gc.Init()

	if err != nil {
		log.Panicln("Error during ncurses Init: ", err)
	}

	stdscr.Keypad(true)
	rand.Seed(int64(time.Now().Second()))
	logFile := initLogging()

	// Finalization
	defer logFile.Close()
	defer gc.End()
	defer gameOver(stdscr)
	defer log.Println(" <==== Game session ended\n ")
	//

	log.Println("====> Game session started")
	initNcurses()

	dimensionsInitError := initScreenDimensions(stdscr)
	if dimensionsInitError != nil {
		log.Panicln("Error initializing the screen dimensions:", dimensionsInitError)
		return
	}

	initScoreBoundFactor()

	ticker := time.NewTicker(time.Second / speedFactor)

	// Create in-game windows
	gameWindow, err := createGameWindow(statsY+statsH, statsX, maxY-statsH, statsW-2)
	if err != nil {
		log.Panic("Error initializing game window:", err)
		return
	}
	menu = createMenu(gameWindow).(*MenuWindow)
	//

	newGame(gameWindow, maxY/2, maxX/2)

	// Game Loop:
	for isRunning {
		select {
		case <-ticker.C:
			if !isPaused {
				handleInput(gameWindow, playerSnake)
				tick(gameWindow)
				drawStats(playerSnake)
				handleEvents(playerSnake, gameWindow)
			} else {
				if !menu.HandleInput() {
					isPaused = false
					menu.Free()
				}
			}
		}
	}

	saveHighScore(stdscr)
}
