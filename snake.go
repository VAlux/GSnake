package main

import (
	ll "Snake/linkedlist"
	"fmt"
	"log"
	"strconv"
	"time"

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
	position point
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
	case 'e':
		udpateObjects(w)
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

func drawDebugStats(maxY int, maxX int, sn *snake, s *gc.Window) {
	snakeLength := "length: " + strconv.Itoa(sn.body.Size())
	dir := "direction: " + sn.direction.String()
	objectsAmount := "objects: " + strconv.Itoa(len(objects))
	rem := "removed: " + removed.String()

	end, err := gc.NewWindow(6, maxX-2, 0, 1)

	if err != nil {
		log.Fatal("Error creating debug window", err)
	}

	end.MovePrint(1, 1, snakeLength)
	end.MovePrint(2, 1, dir)
	end.MovePrint(3, 1, objectsAmount)
	end.MovePrint(4, 1, rem)
	end.Box(gc.ACS_VLINE, gc.ACS_HLINE)
	end.Refresh()
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

	ticker := time.NewTicker(time.Second / 6)
	maxY, maxX := stdscr.MaxYX()
	snake := createSnake(maxY/2, maxX/2)
	for {
		stdscr.Refresh()
		drawDebugStats(maxY, maxX, snake, stdscr)
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
