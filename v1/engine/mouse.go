package engine

const (
	MouseMoveEventType = "MOUSE_MOVE"
)

type MouseMoveEvent struct {
	x, y int32
}

func newMouseMoveEvent(x, y int32) *MouseMoveEvent {
	return &MouseMoveEvent{x, y}
}

func (e *MouseMoveEvent) GetType() string {
	return MouseMoveEventType
}

func (e *MouseMoveEvent) GetPosition() (int32, int32) {
	return e.x, e.y
}
