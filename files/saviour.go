package files

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Saviour struct {
	X, Y float64
	sprite *ebiten.Image
	active bool 

	spriteChangeTimer *Timer
	spritePos int
	dir int
}

const (
	SaviourVelX = 3.0 
)

func NewSaviour(x, y float64) *Saviour {
	log.Printf("num of saviour sprites : %v", len(SaviourSprites))
	return &Saviour{
		X: x,
		Y: y,
		sprite: SaviourSprites[0],
		spritePos: 0,
		active: true,
		spriteChangeTimer: NewTimer(250 * time.Millisecond),
		dir : 1,
	}
}

func (s *Saviour) Update() error {
	if !s.spriteChangeTimer.IsActive() {
		s.spriteChangeTimer.Start()
	}

	if s.spriteChangeTimer.IsReady() {
		s.spritePos += s.dir

		if s.spritePos >= len(SaviourSprites)-1 {
			s.spritePos = len(SaviourSprites) - 1
			s.dir = -1
		}

		if s.spritePos <= 0 {
			s.spritePos = 0
			s.dir = 1
		}

		s.sprite = SaviourSprites[s.spritePos]
		s.spriteChangeTimer.Reset()
	}

	if s.X > float64(ScreenW)+50 {
		s.active = false
	}

	if s.active {
		s.X += SaviourVelX
	}

	return nil
}


func (s *Saviour) Draw(screen *ebiten.Image) {
	opts := & ebiten.DrawImageOptions{}

	opts.GeoM.Scale(-0.35, 0.35)
	opts.GeoM.Translate(s.X, s.Y)
	screen.DrawImage(s.sprite, opts)
}

func (s *Saviour) GetRect() Rect {
	w, h := s.sprite.Bounds().Dx(), s.sprite.Bounds().Dy()
	return NewRect(s.X, s.Y, float64(w)*0.6, float64(h)*0.6)
}