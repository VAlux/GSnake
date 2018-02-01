package main

import (
	ll "Snake/linkedlist"
	"log"
	"time"

	gc "github.com/rthornton128/goncurses"
)

const headTexture = `#`
const tailTexture = `o`
const foodTexture = `*`

var up = &point{0, -1}
var down = &point{0, 1}
var left = &point{-1, 0}
var right = &point{1, 0}
var nowhere = &point{0, 0}

var objects = make([]object, 0)

type point struct {
	x, y int
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
	position point
	color    int
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

	dx := s.head.Data.(point).x + offset.x
	dy := s.head.Data.(point).y + offset.y
	newHead := &ll.Node{Data: point{dx, dy}}
	s.head = newHead
	last := s.body.Back()
	w.MoveDelChar(last.Data.(point).y, last.Data.(point).x)
	s.body.Prepend(newHead)
	s.body.RemoveLast()
}

func (s *snake) draw(w *gc.Window) {
	for node := s.head; node.Next() != nil; node = node.Next() {
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

func tick(t *time.Ticker, w *gc.Window) {
	udpateObjects(w)
	drawObjects(w)
}

func handleInput(w *gc.Window, s *snake) bool {
	key := w.GetChar()
	switch byte(key) {
	case 'w':
		s.direction = up
		return true
	case 's':
		s.direction = down
		return true
	case 'a':
		s.direction = left
		return true
	case 'd':
		s.direction = right
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
	end, err := gc.NewWindow(5, len(msg)+4, (lines/2)-2, (cols-len(msg))/2)
	if err != nil {
		log.Fatal("game over:", err)
	}
	end.MovePrint(2, 2, msg)
	end.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	end.Refresh()
	gc.Nap(2000)
}

func createSnake(x, y int) *snake {
	head := &ll.Node{Data: point{x, y}}
	body := ll.New()

	body.Append(head)
	body.Append(&ll.Node{Data: point{x + 1, y}})
	body.Append(&ll.Node{Data: point{x + 2, y}})
	body.Append(&ll.Node{Data: point{x + 3, y}})
	body.Append(&ll.Node{Data: point{x + 4, y}})

	newSnake := &snake{head, body, right}
	objects = append(objects, newSnake)
	return newSnake
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
	gc.HalfDelay(1)

	ticker := time.NewTicker(time.Second / 6)
	snake := createSnake(5, 5)

	for {
		stdscr.Refresh()
		select {
		case <-ticker.C:
			tick(ticker, stdscr)
		default:
			if !handleInput(stdscr, snake) {
				return
			}
		}
	}
}
