package files

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GRAVITY  = 0.4
	JUMP_VEL = -9
)

type Player struct {
	X, Y      float64
	VelY      float64
	sprite    *ebiten.Image
	jumpTimer *Timer
}

func NewPlayer() *Player {
	var sprite = loadAsset("assets/yellowbird-upflap.png")

	return &Player{
		X:         100,
		Y:         20,
		sprite:    sprite,
		jumpTimer: NewTimer(200 * time.Millisecond),
	}
}

func (p *Player) Update(active bool) error {
	if !active {
		return nil 
	}
	
	if !p.jumpTimer.IsActive() {
		p.jumpTimer.Start()
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && p.jumpTimer.IsReady() {
		p.VelY = JUMP_VEL
		p.jumpTimer.Reset()
	}

	p.VelY += GRAVITY
	p.Y += p.VelY

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(1, 1)

	opts.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.sprite, opts)
}

func (p *Player) GetRect() Rect {
	w, h := p.sprite.Bounds().Dx(), p.sprite.Bounds().Dy()
	return NewRect(p.X, p.Y, float64(w), float64(h))
}
