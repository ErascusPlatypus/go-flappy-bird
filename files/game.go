package files

import (
	"math/rand"
	"time"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var ScreenW, ScreenH int

type Game struct {
	player             *Player
	pipes              []*PipePair
	coins              []*Coins
	magnets            []*Magnet
	magnetActive       bool
	score              int
	highScore          int
	pipeSpeed          float64
	pipeSpeedTimer     *Timer
	endScreenTimer     *Timer
	magnetSpawnTimer   *Timer
	magnetActiveTimer  *Timer
	abilityActiveTimer *Timer
	abilityTimer       *Timer
	abilityActive      bool
	ability            *Ability
	saviours           []*Saviour
	saviourTimer       *Timer
	enemies            []*Enemy
	enemySpawnTimer    *Timer

	gameOver, startScreen bool
	inTransition          bool
	transitionY           int
}

func NewGame() *Game {
	return &Game{
		player:    NewPlayer(),
		ability:   NewAbility(),
		pipes:     []*PipePair{},
		coins:     []*Coins{},
		score:     0,
		highScore: 0,

		pipeSpeed:      2.5,
		pipeSpeedTimer: NewTimer(10 * time.Second),

		magnetSpawnTimer:  NewTimer(8 * time.Second),
		magnetActiveTimer: NewTimer(3 * time.Second),
		magnetActive:      false,

		abilityActiveTimer: NewTimer(5 * time.Second),
		abilityTimer:       NewTimer(10 * time.Second),
		abilityActive:      false,

		saviours:     []*Saviour{},
		saviourTimer: NewTimer(5 * time.Second),

		enemies: []*Enemy{},
		enemySpawnTimer: NewTimer(6 * time.Second),

		gameOver:    false,
		startScreen: true,

		inTransition:   false,
		transitionY:    0,
		endScreenTimer: NewTimer(3 * time.Second),
	}
}

func (g *Game) spawnMagnet() {
	if !g.magnetSpawnTimer.IsActive() {
		g.magnetSpawnTimer.Start()
		return
	}

	if !g.magnetSpawnTimer.IsReady() {
		return
	}

	if rand.Float64() < 0.5 {
		g.magnets = append(g.magnets, NewMagnet())
	}

	g.magnetSpawnTimer.Reset()
}

func (g *Game) updatePipes(active bool) {
	if active && (len(g.pipes) == 0 || g.pipes[len(g.pipes)-1].Top.X < 200) {
		pipes, shift := NewPipePair(ScreenH, ScreenW)
		g.pipes = append(g.pipes, pipes)
		g.coins = append(g.coins, NewCoins(shift))
	}
}

func (g *Game) canSpawnEnemyAt(y float64) bool {
	for _, e := range g.enemies {
		if e.active && math.Abs(e.Y-y) < 100 {
			return false
		}
	}
	return true
}


func (g *Game) spawnEnemy() {
	if rand.Float64() < 0.3 {
		return
	}

	const maxAttempts = 5

	for i := 0; i < maxAttempts; i++ {
		posY := rand.Float64() * (float64(ScreenH) - 50)

		if g.canSpawnEnemyAt(posY) {
			g.enemies = append(
				g.enemies,
				NewEnemy(float64(ScreenW), posY),
			)
			return
		}
	}
}


func (g *Game) handleEnemy() {
	if !g.enemySpawnTimer.IsActive() {
		g.enemySpawnTimer.Start()
	}

	if g.enemySpawnTimer.IsReady() {
		chance := rand.Float64() 
		if chance <= 0.8 {
			for i := 0 ; i < 6 ; i++ {
				g.spawnEnemy() 
			}
		}
		g.enemySpawnTimer.Reset()
	}

	var active []*Enemy
    for _, e := range g.enemies {
        e.Update()
        if e.active {
            active = append(active, e)
        }
    }
    g.enemies = active

	for _, e := range g.enemies {
		e.Update()
	}
}

func (g *Game) activateAbility() {
	g.abilityActive = true
	if !g.abilityActiveTimer.IsActive() {
		g.abilityActiveTimer.Start()
	}
}

func (g *Game) handleAbilityState() {
	if g.abilityActiveTimer.IsReady() {
		g.abilityActiveTimer.Stop()
		g.abilityActive = false
		g.abilityTimer.Reset()
	}
}

func (g *Game) handleAbility() {
	if !g.abilityTimer.IsActive() {
		g.abilityTimer.Start()
	}

	if !g.startScreen && g.abilityTimer.IsReady() && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.activateAbility()
		g.abilityTimer.Stop()
	}
}

func (g *Game) spawnSaviours() {
	x, y := 0.0, 50.0
	for i := 0; i < 5; i++ {
		g.saviours = append(g.saviours, NewSaviour(x, y))
		y += 100.0
	}
}

func (g *Game) handleSaviours() {
	if !g.saviourTimer.IsActive() {
		g.saviourTimer.Start()
	}

	if !g.startScreen && g.saviourTimer.IsReady() && ebiten.IsKeyPressed(ebiten.KeyR) {
		g.spawnSaviours()
		g.saviourTimer.Reset()
	}

	for _, s := range g.saviours {
		s.Update()
	}
}

func (g *Game) handleScreens() bool {
	if g.startScreen && !g.inTransition {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.startScreen = false
		}
		return true
	}

	if g.gameOver && !g.inTransition {
		if !g.endScreenTimer.IsActive() {
			g.endScreenTimer.Start()
		}
		if g.endScreenTimer.IsReady() {
			g.inTransition = true
			g.transitionY = 0
		}
		return true
	}

	if g.inTransition {
		g.transitionY += 5
		if g.transitionY >= ScreenH {
			g.inTransition = false
			g.endScreenTimer.Stop()
			g.resetScreen()
		}
		return true
	}

	return false
}

func (g *Game) updateDifficulty() {
	if !g.pipeSpeedTimer.IsActive() {
		g.pipeSpeedTimer.Start()
	}

	if g.pipeSpeedTimer.IsReady() {
		g.pipeSpeed += 1.0
		g.pipeSpeedTimer.Reset()
	}
}

func (g *Game) updatePipesList() {
	var active []*PipePair
	for _, pp := range g.pipes {
		pp.Update(g.pipeSpeed, g.abilityActive)
		if pp.active {
			active = append(active, pp)
		}
	}
	g.pipes = active
}

func (g *Game) updateCoinsList() {
	var active []*Coins
	for _, c := range g.coins {
		c.Update(g.pipeSpeed, g.magnetActive, g.player.X, g.player.Y)
		if c.coins[4].active {
			active = append(active, c)
		}
	}
	g.coins = active
}

func (g *Game) updateMagnetsList() {
	var active []*Magnet
	for _, m := range g.magnets {
		m.Update(g.pipeSpeed, g.magnetActiveTimer.IsActive())
		if m.active {
			active = append(active, m)
		}
	}
	g.magnets = active
}

func (g *Game) updateEntities() {
	g.player.Update(true)

	g.updatePipes(true)
	g.handleEnemy()
	g.handleAbility()
	g.handleSaviours()
	g.spawnMagnet()

	g.updatePipesList()
	g.updateCoinsList()
	g.updateMagnetsList()
}

func (g *Game) handlePipeCollision() {
	for _, p := range g.pipes {
		if p.active &&
			(p.Top.GetRect().Intersects(g.player.GetRect()) ||
				p.Bottom.GetRect().Intersects(g.player.GetRect())) {

			g.highScore = max(g.score, g.highScore)
			g.gameOver = true
			return
		}
	}
}

func (g *Game) handleCoinCollision() {
	for _, cc := range g.coins {
		for _, c := range cc.coins {
			if c.active && c.GetRect().Intersects(g.player.GetRect()) {
				g.score++
				c.active = false
			}
		}
	}
}

func (g *Game) activateMagnet() {
	g.magnetActive = true
	if !g.magnetActiveTimer.IsActive() {
		g.magnetActiveTimer.Start()
	}
}

func (g *Game) handleMagnetState() {
	if g.magnetActiveTimer.IsReady() {
		g.magnetActiveTimer.Stop()
		g.magnetActive = false
	}
}

func (g *Game) handleMagnetPickup() {
	for _, m := range g.magnets {
		if m.active && m.GetRect().Intersects(g.player.GetRect()) {
			g.activateMagnet()
			if g.magnetActiveTimer.IsReady() {
				m.active = false
			}
		}
	}
}

func (g *Game) handleEnemyCollision() {
	for _, e := range g.enemies {
		for _, s := range g.saviours {
			if s.active && e.active && s.GetRect().Intersects(e.GetRect()) {
				e.active = false
			}
		}

		if e.GetRect().Intersects(g.player.GetRect()) {
			g.highScore = max(g.score, g.highScore)
			g.gameOver = true
			return
		}
	}
}

func (g *Game) handleCollisions() {
	g.handlePipeCollision()
	g.handleCoinCollision()
	g.handleMagnetPickup()
	g.handleEnemyCollision() 
}

func (g *Game) Update() error {
	if g.handleScreens() {
		return nil
	}

	g.updateDifficulty()
	g.updateEntities()
	g.handleCollisions()
	g.handleMagnetState()
	g.handleAbilityState()

	return nil
}

func (g *Game) resetScreen() {
	g.player = NewPlayer()
	g.pipes = []*PipePair{}
	g.coins = []*Coins{}
	g.magnets = []*Magnet{}
	g.score = 0
	g.pipeSpeed = 2.5
	g.gameOver = false
	g.startScreen = true

	g.pipeSpeedTimer.Stop()
	g.magnetSpawnTimer.Stop()
	g.magnetActiveTimer.Stop()
	g.magnetActive = false

	g.abilityTimer.Stop()
	g.abilityActiveTimer.Stop()
	g.abilityActive = false

	g.saviours = []*Saviour{} 
	g.saviourTimer.Stop()

	g.enemies = []*Enemy{} 
	g.enemySpawnTimer.Stop()
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

	for _, m := range g.magnets {
		if m.active {
			m.Draw(screen)
		}
	}

	for _, s := range g.saviours {
		if s.active {
			s.Draw(screen)
		}
	}

	for _, e := range g.enemies {
		if e.active {
			e.Draw(screen)
		}
	}

	if g.abilityTimer.IsReady() {
		g.ability.Draw(screen)
	}

	if g.gameOver {
		g.drawEndScreen(screen, &ebiten.DrawImageOptions{}, true)
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	ScreenW = outsideW
	ScreenH = outsideH
	return outsideW, outsideH
}
