package files

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	player         *Player
	pipes          []*PipePair
	score          int
	pipeSpeed      float64
	pipeSpeedTimer *Timer
	endScreenTimer *Timer

	ScreenW, ScreenH int

	gameOver, startScreen bool
	inTransition bool 
	transitionY int 
}

func NewGame() *Game {
	return &Game{
		player:         NewPlayer(),
		pipes:          []*PipePair{},
		score:          0,
		pipeSpeed:      2.5,
		pipeSpeedTimer: NewTimer(10 * time.Second),
		endScreenTimer: NewTimer(3 * time.Second),
		gameOver:       false,
		startScreen:    true,

		inTransition: false,
		transitionY: 0,
	}
}

func (g *Game) updatePipes(active bool) {
	if active && (len(g.pipes) == 0 || g.pipes[len(g.pipes)-1].Top.X < 200) {
		g.pipes = append(g.pipes, NewPipePair(g.ScreenH, g.ScreenW))
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

		if g.transitionY >= g.ScreenH {
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

	for _, pair := range g.pipes {
		if pair.active &&
			(pair.Top.GetRect().Intersects(g.player.GetRect()) ||
				pair.Bottom.GetRect().Intersects(g.player.GetRect())) {
			g.gameOver = true
			return nil
		}
	}

	return nil
}

func (g *Game) resetScreen() {
	g.player = NewPlayer()
	g.pipes = []*PipePair{}
	g.score = 0 
	g.pipeSpeed = 2.5
	g.gameOver = false 
	g.startScreen = true 
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen)
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

	if g.gameOver {
		g.drawEndScreen(screen, &ebiten.DrawImageOptions{})
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	g.ScreenW = outsideW
	g.ScreenH = outsideH
	return outsideW, outsideH
}

func (g *Game) drawStartScreen(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	startScreenText := "assets/message.png"
	overlay := loadAsset(startScreenText)

	local := *opts 
	local.GeoM.Translate(250, 100)
	screen.DrawImage(overlay, &local)
}

func (g *Game) drawEndScreen(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	endScreenText := "assets/gameover.png"
	overlay := loadAsset(endScreenText)

	local := *opts
 	local.GeoM.Translate(float64(g.ScreenW)/2-80, float64(g.ScreenH)/2)
	screen.DrawImage(overlay, &local)
}

func (g *Game) drawTransition(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, -float64(g.transitionY))
	g.drawEndScreen(screen, opts)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, float64(g.ScreenH) - float64(g.transitionY)) 
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
