package files

import (
	"math/rand"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
)

type Magnet struct {
	sprite *ebiten.Image
	X, Y float64
	active bool 
}

func NewMagnet() *Magnet {
	yPos := rand.Float64() * float64(ScreenH)

	return &Magnet{
		sprite: MagnetSprite,
		X: float64(ScreenW) - 150,
		Y: yPos,
		active: true,
	} 
}
 
func (m *Magnet) Update(speed float64, magnetActive bool) error {
	if magnetActive {
		targetX := 25.0
		targetY := 425.0

		dx := targetX - m.X
		dy := targetY - m.Y

		dist := math.Hypot(dx, dy)

		if dist > 0 {
			vx := (dx / dist) * speed * 2 
			vy := (dy / dist) * speed * 2

			if dist < speed {
				m.X = targetX
				m.Y = targetY
			} else {
				m.X += vx
				m.Y += vy
			}
		}

		return nil
	}

	if m.X < -50 {
		m.active = false
	}

	if m.active {
		m.X -= speed
	}

	return nil
}

func (m *Magnet) Draw(screen *ebiten.Image,) {
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Scale(0.2, 0.2)
	opts.GeoM.Translate(m.X, m.Y)
	screen.DrawImage(m.sprite, opts)
}

func (m *Magnet) GetRect() Rect {
    w, h := m.sprite.Bounds().Dx(), m.sprite.Bounds().Dy()
    return NewRect(m.X, m.Y, float64(w)*0.2, float64(h)*0.2)
}
