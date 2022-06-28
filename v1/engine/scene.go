package engine

type Event interface {
	GetType() string
}

type Scene struct {
	name     string
	setup    func(s *Scene) error
	update   func(s *Scene) error
	drawings map[string]drawable
	fps      uint32
	events   []Event
	running  bool
}

// NewScene will create a new scene
func NewScene(n string, s, u func(s *Scene) error) *Scene {
	return &Scene{
		name:     n,
		setup:    s,
		update:   u,
		drawings: make(map[string]drawable),
		running:  false,
	}
}

// GetName will return a string representing scene name
func (s *Scene) GetName() string {
	return s.name
}

// GetDrawableObjects return a list containing all scene drawings
func (s *Scene) GetDrawableObjects() []drawable {
	l := []drawable{}
	for _, d := range s.drawings {
		l = append(l, d)
	}
	return l
}

// Run will start the scene
func (s *Scene) Run() error {
	if !s.running {
		s.setup(s)
		s.running = true
	}

	s.update(s)

	return nil
}

// Stop will set scene running state to false
func (s *Scene) Stop() error {
	if s.running {
		s.running = false
	}

	return nil
}

// SetFPS will set current fps value for scene
func (s *Scene) SetFPS(fps uint32) {
	s.fps = fps
}

// GetFPS return current fps value
func (s *Scene) GetFPS() uint32 {
	return s.fps
}

func (s *Scene) GetEvents() []Event {
	ev := s.events
	s.clearEvents()
	return ev
}

func (s *Scene) clearEvents() {
	s.events = []Event{}
}

func (s *Scene) addEvent(e Event) {
	s.events = append(s.events, e)
}
