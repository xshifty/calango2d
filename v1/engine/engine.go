package engine

import (
	"fmt"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type drawable interface {
	Draw(rend *sdl.Renderer)
	IsEnabled() bool
}

type sceneAware interface {
	GetName() string
	GetDrawableObjects() []drawable
	Run() error
	Stop() error
	SetFPS(fps uint32)
}

type engine struct {
	title        string
	width        int
	height       int
	window       *sdl.Window
	renderer     *sdl.Renderer
	currentScene string
	scenes       map[string]sceneAware
	running      bool
	initialized  bool
	fullscreen   bool
	timerFPS     uint64
	lastFrame    uint64
	frameRate    uint32
	fps          uint32
	mode         *sdl.DisplayMode
}

// New create a new engine instance
func New(t string, w, h int, frate uint32) *engine {
	e := engine{
		title:       t,
		width:       w,
		height:      h,
		window:      nil,
		renderer:    nil,
		scenes:      make(map[string]sceneAware),
		running:     false,
		initialized: false,
		fullscreen:  false,
		timerFPS:    0,
		lastFrame:   0,
		frameRate:   frate,
		fps:         0,
		mode:        nil,
	}

	e.AddScene(NewScene("__default__", func(s *Scene) error {
		w, h := e.window.GetSize()

		rectWidth := int32(10)
		rectHeight := int32(10)

		maxW := w / rectWidth
		maxH := h / rectHeight

		for x := int32(0); x <= maxW; x++ {
			for y := int32(0); y <= maxH; y++ {
				if (x+y)%2 == 0 {
					s.drawings[fmt.Sprintf("rect_%d_%d", x, y)] = NewRect(int32(x*rectWidth), int32(y*rectHeight), int32(rectWidth), int32(rectHeight), AlphaColor(0x85858520))
				}
			}
		}

		s.drawings["cursor"] = NewRect(0, 0, 40, 40, AlphaColor(0xffcc00ff))

		return nil
	}, func(s *Scene) error {
		max_entropy := 2
		min_entropy := 1

		for k, r := range s.drawings {
			if k == "cursor" {
				continue
			}

			cx, cy := r.(*Rect).GetPosition()

			fx := int32(rand.Intn(max_entropy-min_entropy) + min_entropy)
			fy := int32(rand.Intn(max_entropy-min_entropy) + min_entropy)

			sx := int32(1)
			sy := int32(1)

			if rand.Intn(2000000)%2 == 0 {
				sx = -1
			}

			if rand.Intn(2000000)%2 == 0 {
				sy = -1
			}

			s.drawings[k].(*Rect).SetPosition(fx+(cx*sx), fy+(cy*sy))
		}

		for _, e := range s.GetEvents() {
			if e.GetType() == MouseMoveEventType {
				m := e.(*MouseMoveEvent)

				mx := m.x - 20
				my := m.y - 20

				s.drawings["cursor"].(*Rect).SetPosition(mx, my)
			}
		}

		return nil
	}))

	e.SetScene("__default__")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	mode, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		panic(err)
	}
	e.mode = &mode

	return &e
}

// AddScene will add a new scene to engine
func (e *engine) AddScene(s sceneAware) {
	if s == nil {
		panic("invalid pointer to scene")
	}

	if _, ok := e.scenes[s.GetName()]; !ok {
		fmt.Println("adding new scene: ", s.GetName())
		e.scenes[s.GetName()] = s
	}
}

// SetScene set current scene to be rendered
func (e *engine) SetScene(n string) error {
	if len(e.currentScene) > 0 {
		if _, ok := e.scenes[e.currentScene]; ok {
			e.scenes[e.currentScene].Stop()
		}
	}

	if _, ok := e.scenes[n]; ok {
		e.currentScene = n
		return nil
	}

	return fmt.Errorf("cannot find scene %s", n)
}

// GetFPS will return current FPS
func (e *engine) GetFPS() uint32 {
	return e.fps
}

// Run the engine
func (e *engine) Run() int {
	if !e.initialized {
		e.setup()
		defer sdl.Quit()
		defer e.window.Destroy()
		defer e.renderer.Destroy()
		e.initialized = true
		e.running = true
	}

	sdl.Delay(1)

	for e.running {
		e.clear()
		e.input()
		if err := e.update(); err != nil {
			panic(err)
		}
		e.draw()
	}

	return 0
}

func (e *engine) GetDesktopFrameRate() uint32 {
	return uint32(e.mode.RefreshRate)
}

func (e *engine) SetFrameRate(frate uint32) {
	e.frameRate = frate
}

func (e *engine) clear() {
	e.renderer.SetDrawColor(0, 0, 0, 255)
	e.renderer.Clear()
}

func (e *engine) setup() {
	win, err := sdl.CreateWindow(
		e.title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(e.width),
		int32(e.height),
		sdl.WINDOW_ALLOW_HIGHDPI|sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN)

	if err != nil {
		panic(err)
	}
	e.window = win

	rend, err := sdl.CreateRenderer(e.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	e.renderer = rend

	e.window.Show()
	e.renderer.SetViewport(&sdl.Rect{0, 0, int32(e.width), int32(e.height)})
}

func (e *engine) update() error {
	if s, ok := e.scenes[e.currentScene]; ok {
		return s.Run()
	}

	return fmt.Errorf("scene %s not found", e.currentScene)
}

func (e *engine) input() {
	if _, ok := e.scenes[e.currentScene]; !ok {
		return
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			e.running = false
			break
		case *sdl.KeyboardEvent:
			switch t.Keysym.Sym {
			case sdl.K_F11:
				if t.State == sdl.RELEASED {
					e.fullscreen = !e.fullscreen

					nw := e.mode.W
					nh := e.mode.H
					m := uint32(sdl.WINDOW_FULLSCREEN_DESKTOP)
					if !e.fullscreen {
						nw = int32(e.width)
						nh = int32(e.height)
						m = 0
					}

					e.window.SetSize(nw, nh)
					_ = e.window.SetFullscreen(m)
				}
				break
			}
		case *sdl.MouseMotionEvent:
			e.scenes[e.currentScene].(*Scene).addEvent(newMouseMoveEvent(t.X, t.Y))
		}
	}
}

func (e *engine) draw() {
	currentFrame := sdl.GetTicks64()
	e.timerFPS = currentFrame - e.lastFrame
	e.lastFrame = currentFrame

	delay := int64(1000/e.frameRate) - int64(e.timerFPS)
	if delay < 0 {
		delay = int64(1000 / e.frameRate)
	}

	if s, ok := e.scenes[e.currentScene]; ok {
		if delay > 0 {
			e.fps = uint32(1000 / delay)
			s.SetFPS(e.fps)
		}

		for _, d := range s.GetDrawableObjects() {
			if d.IsEnabled() {
				d.Draw(e.renderer)
			}
		}
	}

	e.renderer.Present()

	sdl.Delay(uint32(delay))
}
