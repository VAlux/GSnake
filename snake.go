package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	ll "Snake/linkedlist"

	gc "github.com/rthornton128/goncurses"
)

const headTexture = `#`
const tailTexture = `o`
const foodTexture = `*`
const emptyTexture = ` `

var up = &point{-1, 0}
var down = &point{1, 0}
var left = &point{0, -1}
var right = &point{0, 1}
var nowhere = &point{0, 0}

var objects = make([]object, 0)

var removed = &ll.Node{Data: &point{0, 0}}

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

func (p *point) String() string {
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

	last := s.body.Back()
	w.MovePrint(last.Data.(point).y, last.Data.(point).x, emptyTexture)
	removed = s.body.RemoveLast()
	s.body.Prepend(newHead)
	s.head = newHead
}

func (s *snake) checkCollision(n *ll.Node) bool {
	return s.body.Contains(n)
}

func (s *snake) draw(w *gc.Window) {
	w.MovePrint(s.head.Data.(point).y, s.head.Data.(point).x, headTexture)
	for node := s.head.Next(); node.Next() != nil; node = node.Next() {
		w.MovePrint(node.Data.(point).y, node.Data.(point).x, tailTexture)
	}
}

func (f *food) update() {
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

func handleInput(w *gc.Window, s *snake) bool {
	key := w.GetChar()
	switch byte(key) {
	case 'w':
		if s.direction != down {
			s.direction = up
		}
		return true
	case 's':
		if s.direction != up {
			s.direction = down
		}
		return true
	case 'a':
		if s.direction != right {
			s.direction = left
		}
		return true
	case 'd':
		if s.direction != left {
			s.direction = right
		}
		return true
	case 'q':
		return false
	default:
		break
	}
	return true
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
	rem := "removed: " + removed.String()

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

	body.Append(&ll.Node{Data: point{y, x + 1}})
	body.Append(&ll.Node{Data: point{y, x + 2}})
	body.Append(&ll.Node{Data: point{y, x + 3}})
	body.Append(&ll.Node{Data: point{y, x + 4}})
	body.Append(&ll.Node{Data: point{y, x + 5}})
	body.Prepend(head)

	newSnake := &snake{head, body, left}
	objects = append(objects, newSnake)
	return newSnake
}

func createWindow(height, width, y, x int) *gc.Window {
	wnd, err := gc.NewWindow(height, width, y, x)
	if err != nil {
		log.Fatal("Error duirng creating window...")
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

	maxY, maxX := stdscr.MaxYX()
	statsX, statsY, statsH, statsW := 1, 0, 6, maxX

	ticker := time.NewTicker(time.Second / 6)
	snake := createSnake(maxY/2, maxX/2)
	gameWindow := createGameWindow(statsY+statsH, statsX, maxY-statsH, statsW-2)

	for {
		select {
		case <-ticker.C:
			tick(gameWindow)
			gameWindow.Refresh()
			drawDebugStats(statsH, statsW, statsY, statsX, snake)
		default:
			if !handleInput(gameWindow, snake) {
				return
			}
		}
	}
}
