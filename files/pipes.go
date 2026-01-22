package files

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// MOVE      = 2.5
	GAP       = 160
)

type Pipe struct {
	X, Y   float64
	sprite *ebiten.Image
}

type PipePair struct {
	Top    *Pipe
	Bottom *Pipe
	active bool 
}

func NewPipePair(SCREEN_H, SCREEN_W int) (*PipePair, float64){
	frac := rand.Float64()
	gap := float64(SCREEN_H) * frac
	shift := gap / 2 
	extra := 25 * frac

	top := &Pipe{
		sprite: PipeSprite,
		X:      float64(SCREEN_W),
		Y:      0 - shift,
	}

	bottom := &Pipe{
		sprite: PipeSprite,
		X:      float64(SCREEN_W),
		Y:      float64(SCREEN_H) - shift - extra,
	}

	return &PipePair{
		Top:    top,
		Bottom: bottom,
		active: true,
	} , shift
}

func (pp *PipePair) Update(move float64) {
	pp.Top.X -= float64(move)
	pp.Bottom.X -= float64(move)

	if pp.Top.X < -50 {
		pp.active = false 
	}
}

func (pp *PipePair) Draw(screen *ebiten.Image) {
	pp.Top.Draw(screen, true)
	pp.Bottom.Draw(screen, false)
}

func (p *Pipe) Draw(screen *ebiten.Image, inverted bool) {
	opts := &ebiten.DrawImageOptions{}

	if inverted {
		h := p.sprite.Bounds().Dy()
		opts.GeoM.Scale(1, -1)
		opts.GeoM.Translate(0, float64(h))
	}

	opts.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.sprite, opts)
}

func (p *Pipe) GetRect() Rect {
    w, h := p.sprite.Bounds().Dx(), p.sprite.Bounds().Dy()
    return NewRect(p.X, p.Y, float64(w), float64(h))
}   
