package main

// Animation declares a way of how to manipulate the internal animation frames
type Animation interface {
	NextFrame() string
	CurrentFrame() string
	MoveFrameIndex()

	framesAmount() int
	hasNextFrame() bool
}

// NewAnimation creates new animation object with specified frames array
func NewAnimation(frames []string, duration int) Animation {
	return &animation{
		frames:            frames,
		currentFrameIndex: 0,
		frameDuration:     duration,
		currentFrameTime:  0}
}

type animation struct {
	frames            []string
	currentFrameIndex int
	frameDuration     int
	currentFrameTime  int
}

func (a *animation) framesAmount() int {
	return len(a.frames)
}

func (a *animation) hasNextFrame() bool {
	return a.currentFrameIndex < a.framesAmount()-1
}

// MoveFrameIndex moves the current frame caret to the next frame.
// If there are no frames left in the sequence - caret will be reset and point to the 0 frame
func (a *animation) MoveFrameIndex() {
	if a.currentFrameTime < a.frameDuration-1 {
		a.currentFrameTime++
	} else {
		a.currentFrameTime = 0
		if a.hasNextFrame() {
			a.currentFrameIndex++
		} else {
			a.currentFrameIndex = 0
		}
	}
}

// CurrentFrame obtains current frame from the animation sequence
func (a *animation) CurrentFrame() string {
	return a.frames[a.currentFrameIndex]
}

// NextFrame get the next animation frame.
// If there are no frames in the animation sequence - sequence will reset from 0
func (a *animation) NextFrame() string {
	a.MoveFrameIndex()
	return a.CurrentFrame()
}
