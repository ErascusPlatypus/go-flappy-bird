package files

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

)

func (g *Game) drawStartScreen(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	startScreenText := "assets/message.png"
	overlay := loadAsset(startScreenText)

	local := *opts
	local.GeoM.Translate(250, 100)
	screen.DrawImage(overlay, &local)
}

func (g *Game) drawEndScreen(screen *ebiten.Image, opts *ebiten.DrawImageOptions, drawScores bool) {
	cx := float64(ScreenW) / 2
	cy := float64(ScreenH) / 2

	overlay := loadAsset("assets/gameover.png")

	local := *opts
	local.GeoM.Translate(cx-100, cy-140)
	screen.DrawImage(overlay, &local)

	if drawScores {
		_, offsetY := opts.GeoM.Apply(0, 0)
		drawScore(screen, g.score, cx-120, cy-10+offsetY, 0, "SCORE:")
		drawScore(screen, g.highScore, cx-120, cy+30+offsetY, 0, "HIGH SCORE:")
	}
}

func (g *Game) drawTransition(screen *ebiten.Image) {
	offsetUp := -float64(g.transitionY)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, offsetUp)
	g.drawEndScreen(screen, opts, true)

	offsetDown := float64(ScreenH) - float64(g.transitionY)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, offsetDown)
	g.drawStartScreen(screen, opts)
}

func drawBackground(screen *ebiten.Image) {
	w, h := BackgroundImage.Bounds().Dx(), BackgroundImage.Bounds().Dy()

	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(
		float64(sw)/float64(w),
		float64(sh)/float64(h),
	)

	screen.DrawImage(BackgroundImage, opts)
}

func drawScore(screen *ebiten.Image, score int, x, y, offsetY float64, prefix string) {
	s := fmt.Sprintf("%s %d", prefix, score)

	offsets := []struct{ x, y float64 }{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	for _, o := range offsets {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x+o.x, y+o.y+offsetY)
		text.DrawWithOptions(screen, s, ScoreFont, opts)
	}

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y+offsetY)
	text.DrawWithOptions(screen, s, ScoreFont, opts)
}
