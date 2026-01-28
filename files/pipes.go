package files

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// MOVE      = 2.5
	GAP = 160
)

type Pipe struct {
	X, Y   float64
	sprite *ebiten.Image
}

type PipePair struct {
	Top                   *Pipe
	Bottom                *Pipe
	active                bool
	gettingDestroyed      bool
	gettingDestroyedTimer *Timer
	sprites               []*ebiten.Image
	spritePos             int
}

func NewPipePair(SCREEN_H, SCREEN_W int) (*PipePair, float64) {
	frac := rand.Float64()
	gap := float64(SCREEN_H) * frac
	shift := gap / 2
	extra := 25 * frac

	top := &Pipe{
		sprite: PipeSprites[0],
		X:      float64(SCREEN_W),
		Y:      0 - shift,
	}

	bottom := &Pipe{
		sprite: PipeSprites[0],
		X:      float64(SCREEN_W),
		Y:      float64(SCREEN_H) - shift - extra,
	}

	return &PipePair{
		Top:                   top,
		Bottom:                bottom,
		active:                true,
		sprites:               PipeSprites,
		spritePos:             0,
		gettingDestroyed:      false,
		gettingDestroyedTimer: NewTimer(350 * time.Millisecond),
	}, shift
}

func (pp *PipePair) Destroy() {
	if pp.gettingDestroyed && pp.gettingDestroyedTimer.IsReady() {
		pp.spritePos++
		log.Printf("current spritePos : %v", pp.spritePos)
		if pp.spritePos >= len(pp.sprites) {
			pp.gettingDestroyed = false
			pp.spritePos = len(pp.sprites) - 1
			pp.gettingDestroyedTimer.Stop()
		} else {
			pp.gettingDestroyedTimer.Reset()
		}
	}
}

func (pp *PipePair) Update(move float64, abilityActive bool) {
	pp.Top.X -= float64(move)
	pp.Bottom.X -= float64(move)

	if pp.Top.X < -50 {
		pp.active = false
	}

	if abilityActive {
		log.Printf("calling destroy function in pipes.go")
		pp.gettingDestroyed = true

		if !pp.gettingDestroyedTimer.IsActive() {
			pp.gettingDestroyedTimer.Start()
		}

		pp.Destroy()
	}

	pp.Top.sprite = pp.sprites[pp.spritePos]
	// log.Printf("set top/btm sprites to %v pos", pp.spritePos)
	pp.Bottom.sprite = pp.sprites[pp.spritePos]
}

func (pp *PipePair) Draw(screen *ebiten.Image) {
	scale := 1.0
	if pp.spritePos > 0 {
		scale = 0.8
	}

	pp.Top.Draw(screen, true, scale)
	pp.Bottom.Draw(screen, false, scale)
}

func (p *Pipe) Draw(screen *ebiten.Image, inverted bool, scaleX float64) {
	opts := &ebiten.DrawImageOptions{}

	if inverted {
		h := p.sprite.Bounds().Dy()
		opts.GeoM.Scale(1, -1)
		opts.GeoM.Translate(0, float64(h))
	}

	opts.GeoM.Scale(scaleX, 1)

	w := float64(p.sprite.Bounds().Dx())
	opts.GeoM.Translate(-w/2, 0)
	opts.GeoM.Scale(scaleX, 1)
	opts.GeoM.Translate(w/2, 0)

	opts.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.sprite, opts)
}

func (p *Pipe) GetRect() Rect {
	w, h := p.sprite.Bounds().Dx(), p.sprite.Bounds().Dy()
	return NewRect(p.X, p.Y, float64(w), float64(h))
}
