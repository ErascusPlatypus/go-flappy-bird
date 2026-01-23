package files

import (
	"math"
	"github.com/hajimehoshi/ebiten/v2"
)

type Coin struct {
	sprite *ebiten.Image
	X, Y   float64
	active bool 
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
		spacingX = 60
		amplitude = 60
	)

	startX := float64(ScreenW) + 40

	for i := 0; i < count; i++ {
		x := startX + float64(i)*spacingX
		y := gapCenterY + amplitude*math.Sin(float64(i)*0.8)

		coins.coins = append(coins.coins, &Coin{
			sprite: CoinSprites[0],
			X:      x - 150,
			Y:      y - 75,
			active: true,
		})
	}

	return coins
}

func (c *Coin) Update(speed float64, magnetActive bool, targetX, targetY float64) error {
	if c.Y < -50 {
		c.active = false 
	}
	
	if magnetActive && c.active {
		dx := targetX - c.X 
		dy := targetY - c.Y 

		dist := math.Hypot(dx, dy)
		if dist > 0 {
			c.X += (dx/dist) * (2 * speed) 
			c.Y += (dy/dist) * (2 * speed)
		}

		return nil 
	}

	if !magnetActive && c.active {
		c.X -= speed
	}

	return nil
}

func (cc *Coins) Update(speed float64, magnetActive bool, X, Y float64) error {
	for _, co := range cc.coins {
		co.Update(speed, magnetActive, X, Y)
	}

	return nil
}

func (cc *Coins) Draw(screen *ebiten.Image) {
	for _, c := range cc.coins {
		if c.active {
			c.Draw(screen)
		} 
	}
}

func (c *Coin) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Scale(0.25, 0.25)
	opts.GeoM.Translate(c.X, c.Y)

	screen.DrawImage(c.sprite, opts)
}

func (c *Coin) GetRect() Rect {
	w := float64(c.sprite.Bounds().Dx())
	h := float64(c.sprite.Bounds().Dy())

	shrink := 0.50
	x := c.X + w*shrink/2
	y := c.Y + h*shrink/2

	return NewRect(
		x,
		y,
		w*(1-shrink),
		h*(1-shrink),
	)
}

