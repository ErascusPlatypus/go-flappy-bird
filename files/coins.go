package files

import (
	"math"
	"github.com/hajimehoshi/ebiten/v2"
)

type Coin struct {
	sprite *ebiten.Image
	X, Y   float64
	VelX   float64
}

type Coins struct {
	coins []*Coin
}

func NewCoins(shift float64) *Coins {
	h := PipeSprite.Bounds().Dy()
	gapCenterY := float64(h-int(shift)) + 100

	coins := &Coins{}

	const (
		count  = 5
		spacingX = 50
		amplitude = 50
	)

	startX := float64(ScreenW) + 40

	for i := 0; i < count; i++ {
		x := startX + float64(i)*spacingX
		y := gapCenterY + amplitude*math.Sin(float64(i)*0.8)

		coins.coins = append(coins.coins, &Coin{
			sprite: CoinSprite,
			X:      x - 125,
			Y:      y - 50,
			VelX:   2.5,
		})
	}

	return coins
}


func (c *Coin) Update() error {
	c.X -= c.VelX

	return nil
}

func (cc *Coins) Update() error {
	for _, co := range cc.coins {
		co.Update()
	}

	return nil
}

func (cc *Coins) Draw(screen *ebiten.Image) {
	for _, c := range cc.coins {
		c.Draw(screen)
	}
}

func (c *Coin) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Scale(0.25, 0.25)
	opts.GeoM.Translate(c.X, c.Y)

	screen.DrawImage(c.sprite, opts)
}
