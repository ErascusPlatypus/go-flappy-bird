package files

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	GRAVITY  = 0.4
	JUMP_VEL = -9
)

type Player struct {
	X, Y       float64
	VelY       float64
	sprite     *ebiten.Image
	spriteIdle *ebiten.Image
	spriteUp   *ebiten.Image
	spriteDown *ebiten.Image
	jumpTimer  *Timer
}

func NewPlayer() *Player {
	var spriteIdle = loadAsset("assets/yellowbird-midflap.png")
	var spriteUp = loadAsset("assets/yellowbird-upflap.png")
	var spriteDown = loadAsset("assets/yellowbird-downflap.png")

	return &Player{
		X:          100,
		Y:          20,
		sprite:     spriteIdle,
		spriteIdle: spriteIdle,
		spriteUp:   spriteUp,
		spriteDown: spriteDown,
		jumpTimer:  NewTimer(200 * time.Millisecond),
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
		log.Printf("jump ")
		p.jumpTimer.Reset()
	}

	p.VelY += GRAVITY
	p.Y += p.VelY

	if p.VelY < 0 {
		p.sprite = p.spriteUp
	} else if p.VelY == 0 {
		p.sprite = p.spriteIdle
	} else {
		p.sprite = p.spriteDown
	}

	// log.Printf("VelY : %v \n", p.VelY)
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
