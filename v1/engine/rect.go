package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Rect struct {
	x, y          int32
	width, height int32
	color         AlphaColor
	enabled       bool
}

func NewRect(x, y, w, h int32, c AlphaColor) *Rect {
	return &Rect{x, y, w, h, c, true}
}

func (r *Rect) Draw(rend *sdl.Renderer) {
	rend.GetDrawColor()
	rend.SetDrawColor(r.color.GetRed(), r.color.GetGreen(), r.color.GetBlue(), r.color.GetAlpha())
	rend.FillRect(&sdl.Rect{r.x, r.y, r.width, r.height})
}

func (r *Rect) IsEnabled() bool {
	return r.enabled
}

func (r *Rect) SetPosition(x, y int32) {
	r.x = x
	r.y = y
}

func (r *Rect) GetPosition() (int32, int32) {
	return r.x, r.y
}
