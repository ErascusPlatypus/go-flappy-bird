package files

import "github.com/hajimehoshi/ebiten/v2"

type Ability struct {
	sprite *ebiten.Image
	X, Y float64
}

func NewAbility() *Ability {
	return &Ability{
		sprite: AbilitySprite,
		X: 65.0,
		Y: 425.0,
	}
}

func (a *Ability) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Scale(0.6, 0.6)
	opts.GeoM.Translate(a.X, a.Y)
	screen.DrawImage(a.sprite, opts)
}