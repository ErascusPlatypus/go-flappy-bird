package files

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	X, Y float64
	sprite *ebiten.Image
	active bool 

	spriteChangeTimer *Timer
	spritePos int
	dir int
}

const (
	EnemyVelX = 2.0
)

func NewEnemy(x, y float64) *Enemy {
	return &Enemy{
		X: x,
		Y: y,
		sprite: EnemySprites[0],
		spritePos: 0,
		active: true,
		spriteChangeTimer: NewTimer(250 * time.Millisecond),
		dir : 1,
	}
}

func (e *Enemy) Update() error {
	if !e.spriteChangeTimer.IsActive() {
		e.spriteChangeTimer.Start()
	}

	if e.spriteChangeTimer.IsReady() {
		e.spritePos += e.dir

		if e.spritePos >= len(EnemySprites)-1 {
			e.spritePos = len(EnemySprites) - 1
			e.dir = -1
		}

		if e.spritePos <= 0 {
			e.spritePos = 0
			e.dir = 1
		}

		e.sprite = EnemySprites[e.spritePos]
		e.spriteChangeTimer.Reset()
	}

	if e.X < -50 {
		e.active = false
	}

	if e.active {
		e.X -= EnemyVelX
	}

	return nil
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	opts := & ebiten.DrawImageOptions{}

	w := e.sprite.Bounds().Dx()
	opts.GeoM.Scale(3, 3)
	opts.GeoM.Translate(float64(w)*3+e.X, e.Y)
	screen.DrawImage(e.sprite, opts)
}

func (e *Enemy) GetRect() Rect {
	w, h := e.sprite.Bounds().Dx(), e.sprite.Bounds().Dy()
	return NewRect(e.X, e.Y, float64(w) * 0.5, float64(h) * 0.5)
}