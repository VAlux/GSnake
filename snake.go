package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	ll "Snake/linkedlist"

	gc "github.com/rthornton128/goncurses"
)

const headTexture = `#`
const tailTexture = `o`
const foodTexture = `*`
const emptyTexture = ` `

const collisionEvent = "collision"
const exitEvent = "exit"

const initialLength = 10

var up = &point{-1, 0}
var down = &point{1, 0}
var left = &point{0, -1}
var right = &point{0, 1}
var nowhere = &point{0, 0}

var objects = make([]object, 0)
var events = make(chan string, 1)

var maxX = 0
var maxY = 0

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

func (p point) String() string {
	return fmt.Sprintf("y: %d, x: %d", p.y, p.x)
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

	last := s.body.Back()
	w.MovePrint(last.Data.(point).y, last.Data.(point).x, emptyTexture)
	s.body.RemoveLast()
	s.body.Prepend(newHead)
	s.head = newHead
}

func (s *snake) draw(w *gc.Window) {
	w.MovePrint(s.head.Data.(point).y, s.head.Data.(point).x, headTexture)
	for node := s.head.Next(); node.Next() != nil; node = node.Next() {
		w.MovePrint(node.Data.(point).y, node.Data.(point).x, tailTexture)
	}
}

func (s *snake) checkCollision(n *ll.Node) bool {
	return n.Data.(point).x <= 0 ||
		n.Data.(point).y <= 0 ||
		n.Data.(point).x >= maxX-1 ||
		n.Data.(point).y >= maxY-1 ||
		s.body.Contains(n)
}

func (s *snake) containsNodeWithPoint(pt *point) bool {
	for node := s.head; node != nil; node = node.Next() {
		if node.Data.(point) == *pt {
			return true
		}
	}
	return false
}

func (f *food) update(w *gc.Window) {
	// TODO: update the food color
}

func (f *food) draw(w *gc.Window) {
	w.MovePrint(f.position.y, f.position.x, foodTexture)
}

func drawObjects(s *gc.Window) {
	for _, obj := range objects {
		obj.draw(s)
	}
}

func udpateObjects(w *gc.Window) {
	for _, obj := range objects {
		obj.update(w)
	}
}

func tick(w *gc.Window) {
	udpateObjects(w)
	drawObjects(w)
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
	case 'q':
		events <- exitEvent
		break
	default:
		break
	}
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

func drawDebugStats(height int, width int, y int, x int, sn *snake) {
	snakeLength := "length: " + strconv.Itoa(sn.body.Size())
	dir := "direction: " + sn.direction.String()
	objectsAmount := "objects: " + strconv.Itoa(len(objects))
	rem := "head: " + sn.head.Data.(point).String()

	wnd := createWindow(height, width-2, y, x)
	wnd.MovePrint(1, 1, snakeLength)
	wnd.MovePrint(2, 1, dir)
	wnd.MovePrint(3, 1, objectsAmount)
	wnd.MovePrint(4, 1, rem)
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
	randX := 1 + rand.Intn(maxX-1)
	randY := 1 + rand.Intn(maxY-1)
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

func main() {
	stdscr, err := gc.Init()

	if err != nil {
		log.Println("Error during ncurses Init:", err)
	}

	defer gc.End()
	defer gameOver(stdscr)

	gc.Cursor(0)
	gc.Echo(false)
	gc.HalfDelay(2)

	rand.Seed(int64(time.Now().Second()))
	maxY, maxX = stdscr.MaxYX()
	//statsX, statsY, statsH, statsW := 1, 0, 6, maxX

	ticker := time.NewTicker(time.Second / 6)
	snake := createSnake(maxY/2, maxX/2)
	food := generateFood(snake)

	objects = append(objects, snake, food)
	//gameWindow := createGameWindow(statsY+statsH, statsX, maxY-statsH, statsW-2)
	gameWindow := createGameWindow(0, 0, maxY, maxX)

	//Game Loop:
	for {
		handleInput(gameWindow, snake)

		select {
		case <-ticker.C:
			tick(gameWindow)
			gameWindow.Refresh()
			//drawDebugStats(statsH, statsW, statsY, statsX, snake)
		case event := <-events:
			if event == collisionEvent {
				return
			}
			if event == exitEvent {
				return
			}
		}
	}
}
