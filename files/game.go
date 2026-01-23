package files

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var ScreenW, ScreenH int

type Game struct {
	player         *Player
	pipes          []*PipePair
	coins          []*Coins
	score          int
	highScore      int
	pipeSpeed      float64
	pipeSpeedTimer *Timer
	endScreenTimer *Timer

	gameOver, startScreen bool
	inTransition          bool
	transitionY           int
}

func NewGame() *Game {
	return &Game{
		player:         NewPlayer(),
		pipes:          []*PipePair{},
		coins:          []*Coins{},
		score:          0,
		highScore:      0,
		pipeSpeed:      2.5,
		pipeSpeedTimer: NewTimer(10 * time.Second),
		endScreenTimer: NewTimer(3 * time.Second),
		gameOver:       false,
		startScreen:    true,

		inTransition: false,
		transitionY:  0,
	}
}

func (g *Game) updatePipes(active bool) {
	if active && (len(g.pipes) == 0 || g.pipes[len(g.pipes)-1].Top.X < 200) {
		pipes, shift := NewPipePair(ScreenH, ScreenW)
		g.pipes = append(g.pipes, pipes)
		g.coins = append(g.coins, NewCoins(shift))
	}
}

func (g *Game) Update() error {
	if g.startScreen && !g.inTransition {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.startScreen = false
		}
		return nil
	}

	if g.gameOver && !g.inTransition {
		if !g.endScreenTimer.IsActive() {
			g.endScreenTimer.Start()
		}

		if g.endScreenTimer.IsReady() {
			g.inTransition = true
			g.transitionY = 0
		}
		return nil
	}

	if g.inTransition {
		g.transitionY += 5

		if g.transitionY >= ScreenH {
			g.inTransition = false
			g.endScreenTimer.Stop()
			g.resetScreen()
		}
		return nil
	}

	if !g.pipeSpeedTimer.IsActive() {
		g.pipeSpeedTimer.Start()
	}

	if g.pipeSpeedTimer.IsReady() {
		g.pipeSpeed += 1.0
		g.pipeSpeedTimer.Reset()
	}

	g.player.Update(true)
	g.updatePipes(true)

	var activePipes []*PipePair
	for _, pair := range g.pipes {
		pair.Update(g.pipeSpeed)
		if pair.active {
			activePipes = append(activePipes, pair)
		}
	}
	g.pipes = activePipes

	var activeCoins []*Coins
	for _, c := range g.coins {
		c.Update(g.pipeSpeed)
		if c.coins[4].active {
			activeCoins = append(activeCoins, c)
		}
	}
	g.coins = activeCoins

	for _, pair := range g.pipes {
		if pair.active &&
			(pair.Top.GetRect().Intersects(g.player.GetRect()) ||
				pair.Bottom.GetRect().Intersects(g.player.GetRect())) {
			g.highScore = max(g.score, g.highScore)
			g.gameOver = true
			return nil
		}
	}

	for _, cc := range g.coins {
		for _, c := range cc.coins {
			if c.active && c.GetRect().Intersects(g.player.GetRect()) {
				g.score++
				c.active = false
			}
		}
	}

	return nil
}

func (g *Game) resetScreen() {
	g.player = NewPlayer()
	g.pipes = []*PipePair{}
	g.coins = []*Coins{}
	g.score = 0
	g.pipeSpeed = 2.5
	g.gameOver = false
	g.startScreen = true
}

func (g *Game) Draw(screen *ebiten.Image) {
    drawBackground(screen)

    drawScore(screen, g.score, 20, 50, 0, "")

    if g.inTransition {
        g.drawTransition(screen)
        return
    }

    if g.startScreen {
        g.drawStartScreen(screen, &ebiten.DrawImageOptions{})
        return
    }

    g.player.Draw(screen)
    for _, pair := range g.pipes {
        pair.Draw(screen)
    }

    for _, c := range g.coins {
        c.Draw(screen)
    }

    if g.gameOver {
        g.drawEndScreen(screen, &ebiten.DrawImageOptions{}, true)  // Draw with scores
    }
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	ScreenW = outsideW
	ScreenH = outsideH
	return outsideW, outsideH
}

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

