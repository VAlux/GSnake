package main

// Animation declares a way of how to manipulate the internal animation frames
type Animation interface {
	NextFrame() string
	CurrentFrame() string

	framesAmount() int
	hasNextFrame() bool
	moveFrameIndex()
}

// NewAnimation creates new animation object with specified frames array
func NewAnimation(frames []string) Animation {
	return &animation{frames, 0}
}

type animation struct {
	frames            []string
	currentFrameIndex int
}

func (a *animation) framesAmount() int {
	return len(a.frames)
}

func (a *animation) hasNextFrame() bool {
	return a.currentFrameIndex < a.framesAmount()-1
}

func (a *animation) moveFrameIndex() {
	if a.hasNextFrame() {
		a.currentFrameIndex++
	} else {
		a.currentFrameIndex = 0
	}
}

// CurrentFrame obtains current frame from the animation sequence
func (a *animation) CurrentFrame() string {
	return a.frames[a.currentFrameIndex]
}

// NextFrame get the next animation frame.
// If there are no frames in the animation sequence - sequence will reset from 0
func (a *animation) NextFrame() string {
	a.moveFrameIndex()
	return a.CurrentFrame()
}
