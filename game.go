package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/llrowat/spriteutils"
	"image/color"
	"math"
	"math/rand"
	"time"
)

const (
	// outofBoundsX represents the location when sprites are considered out of bounds and will be destroyed
	outOfBoundsX = -200
)

// Game represents the game state
type Game struct {
	// mode is the current game mode
	mode Mode

	// ship is the main character ship sprite
	ship   *spriteutils.Sprite
	// shield is the main character's ship shield sprite
	shield *spriteutils.Sprite

	// topGroundTiles are the floor tile sprites at the top of the screen
	topGroundTiles    []*spriteutils.Sprite
	// bottomGroundTiles are the floor tile sprites at the bottom of the screen
	bottomGroundTiles []*spriteutils.Sprite

	// topSpireFactory is a factory for generating spires at the top of the screen
	topSpireFactory    *spriteutils.SpriteFactory
	// bottomSpireFactory is the factory for generating spires at the bottom of the screen
	bottomSpireFactory *spriteutils.SpriteFactory
	// asteroidFactory is the factory for generating asteroids
	asteroidFactory    *spriteutils.SpriteFactory
	// starFactory is the factory for generating stars
	starFactory        *spriteutils.SpriteFactory
	// spires are all the spire sprites currently in the game
	spires             []*spriteutils.Sprite
	// asteroids are all the asteroid sprites currently in the game
	asteroids          []*spriteutils.Sprite
	// asteroidExplosions are transient sprites that exist temporarily when asteroids collide with other objects
	asteroidExplosions []*spriteutils.TransientSprite
	// stars are all the star sprites currently in the game
	stars              []*spriteutils.Sprite

	// distanceTravelled represents the current distance travelled in game (basically the score)
	distanceTravelled      int
	// speed represents how fast the player moves through the world (or actually how fast the world moves around the player)
	speed                  float64
	// speedIncreaseThreshold represents the distance that the next speed increase will occur
	speedIncreaseThreshold int
	// boostFactor is the amount speed will increase when the player hits a star
	boostFactor            float64
	// boostSeconds is how long the boost will last
	boostSeconds           int64
	// isBoosting represents whether the player is currently undergoing a boost
	isBoosting             bool
	// lastBoostTime represents the start duration (time that has elapsed since game started) of the last boost
	lastBoostTime          time.Duration
	// spireSpawnThreshold represents the distance that the next spire will spawn
	spireSpawnThreshold    int
	// asteroidSpawnThreshold represents the distance that the next asteroid will spawn
	asteroidSpawnThreshold int
	// starSpawnThreshold represents the distance that the next star will spawn
	starSpawnThreshold     int

	// frameCount is the current frame the game is on since it has started
	frameCount int64
}

// Initialize by resetting game state to initial
func (g *Game) init() {
	g.resetGame()
}

// resetGame Resets game start to initial state
func (g *Game) resetGame() {
	g.ship = &spriteutils.Sprite{
		Image:     shipImage,
		X:         screenWidth / 4,
		Y:         screenHeight / 2,
		XVelocity: 0,
		YVelocity: 0,
		Rotation:  0,
	}
	g.shield = nil

	g.distanceTravelled = 0
	g.frameCount = 0
	g.isBoosting = false
	g.boostFactor = 2
	g.lastBoostTime = 0
	g.speed = 1
	g.speedIncreaseThreshold = 500
	g.spireSpawnThreshold = 600
	g.asteroidSpawnThreshold = 200
	g.starSpawnThreshold = 50

	g.asteroidExplosions = nil

	g.initializeGround()
	g.initializeSpireFactories()
	g.initializeAsteroidFactories()
	g.initializeStarFactory()
}

// Update runs the game loop logic
func (g *Game) Update(screen *ebiten.Image) error {

	switch g.mode {
	case ModeTitle:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.mode = ModeGame
		}
	case ModeGame:
		{
			// Increase speed periodically
			g.distanceTravelled += int(g.speed)
			if g.distanceTravelled > g.speedIncreaseThreshold {
				g.speedIncreaseThreshold += g.speedIncreaseThreshold
				g.speed++
			}

			// Check whether boost duration has elapsed
			if g.isBoosting && time.Duration(g.frameCount)*time.Second/60-g.lastBoostTime > time.Second*5 {
				g.speed -= g.boostFactor
				g.isBoosting = false
			}

			g.shipMovement()
			g.checkShieldOn()

			g.updateGround()
			g.updateSpires()
			g.updateAsteroids()
			g.updateStars()

			g.checkCollisions()

			// Generate Spires
			if g.distanceTravelled > g.spireSpawnThreshold {
				if rand.Intn(2)%2 == 0 {
					g.spires = append(g.spires, g.topSpireFactory.GenerateSprite())
				} else {
					g.spires = append(g.spires, g.bottomSpireFactory.GenerateSprite())
				}
				g.spireSpawnThreshold += 600
			}

			// Generate asteroids and apply random impulse
			if g.distanceTravelled > g.asteroidSpawnThreshold {
				g.asteroids = append(g.asteroids, g.asteroidFactory.GenerateSprite())
				g.asteroids[len(g.asteroids)-1].ApplyImpulse(float64(rand.Intn(10))-15, float64(rand.Intn(6))-3)
				g.asteroidSpawnThreshold += 200
			}

			// Generate Stars
			if g.distanceTravelled > g.starSpawnThreshold {
				g.stars = append(g.stars, g.starFactory.GenerateSprite())
				g.starSpawnThreshold += 2000
			}

			// Handle explosions
			temp := g.asteroidExplosions[:0]
			for _, asteroidExplosion := range g.asteroidExplosions {
				if !asteroidExplosion.IsExpired {
					temp = append(temp, asteroidExplosion)
				}
				asteroidExplosion.Update(time.Duration(g.frameCount) * time.Second / 60)
			}
			g.asteroidExplosions = temp

			g.frameCount++
		}
	case ModeGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.resetGame()
			g.mode = ModeTitle
		}
	}

	return nil
}

// Draw draws all the game assets to screen
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)

	// Draw all stars
	for _, star := range g.stars {
		star.Draw(screen)
	}

	// Draw all spires
	for _, spire := range g.spires {
		spire.Draw(screen)
	}

	// Draw  floor tiles
	for _, tile := range g.topGroundTiles {
		tile.Draw(screen)
	}
	for _, tile := range g.bottomGroundTiles {
		tile.Draw(screen)
	}

	// Draw all asteroids
	for _, asteroid := range g.asteroids {
		asteroid.Draw(screen)
	}

	// Draw all asteroid explosions
	for _, asteroidExplosion := range g.asteroidExplosions {
		asteroidExplosion.Draw(screen)
	}

	// Draw ship and shield is it is enabled
	g.ship.Draw(screen)
	if g.shield != nil {
		g.shield.Draw(screen)
	}

	var titleTexts []string
	var texts []string

	// Draw game text
	switch g.mode {
	case ModeTitle:
		titleTexts = []string{"GALACTIC ASTEROID BELT"}
		texts = []string{"", "", "", "", "", "", "", "PRESS SPACE KEY"}
	case ModeGame:
		g.drawScore(screen)
	case ModeGameOver:
		titleTexts = []string{"GAME OVER!"}
		texts = []string{"", "", "", "", "", "", fmt.Sprintf("DISTANCE TRAVELLED: %d M", g.distanceTravelled), "", "", "", "PRESS 'R' KEY TO RESTART"}
	}
	for i, l := range titleTexts {
		x := (screenWidth - len(l)/2*titleFontSize) / 2
		text.Draw(screen, l, titleFont, x, screenHeight/4+(i+4)*titleFontSize, color.White)
	}
	for i, l := range texts {
		x := (screenWidth - len(l)/2*fontSize) / 2
		text.Draw(screen, l, normalFont, x, screenHeight/4+(i+4)*fontSize, color.White)
	}


	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}

// Layout scales the logical game size with the window size.  We don't do anything here, just return the fixed window size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// shipMovement handles all the logic for moving the player character ship
func (g *Game) shipMovement() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.ship.YVelocity -= 0.5
	}

	// Gravity
	g.ship.YVelocity += 0.25

	g.ship.Update()

	// The ship rotates a little bit when moving up/down to give it some "floatiness"
	g.ship.Rotation = float64(g.ship.YVelocity) / 96.0 * math.Pi / 2
}

// checkShieldOn checks whether the ship shield should be enabled
func (g *Game) checkShieldOn() {
	if g.isBoosting {
		g.shield = &spriteutils.Sprite{
			Image:     shieldImage,
			X:         g.ship.X - 17,
			Y:         g.ship.Y - 15,
		}
	} else {
		g.shield = nil
	}
}

// drawBackground draws the background image
func (g *Game) drawBackground(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	imageWidth, imageHeight := backgroundImage.Size()
	maxScale := math.Max(float64(screenWidth)/float64(imageWidth), float64(screenHeight)/float64(imageHeight))
	op.GeoM.Scale(maxScale, maxScale)
	screen.DrawImage(backgroundImage, op)
}

// initializeGround sets the initial state of the floor tiles
func (g *Game) initializeGround() {
	imageWidth, imageHeight := floorImage.Size()

	g.topGroundTiles = nil
	g.bottomGroundTiles = nil

	for i := 0; (i * imageWidth) < (screenWidth + imageWidth*2); i++ {
		topTile := &spriteutils.Sprite{
			Image:     floorImage,
			X:         imageWidth*i - g.distanceTravelled,
			XVelocity: -g.speed,
			Rotation:  math.Pi,
		}

		bottomTile := &spriteutils.Sprite{
			Image:     floorImage,
			X:         imageWidth*i - g.distanceTravelled,
			Y:         screenHeight - imageHeight,
			XVelocity: -g.speed,
		}

		g.topGroundTiles = append(g.topGroundTiles, topTile)
		g.bottomGroundTiles = append(g.bottomGroundTiles, bottomTile)
	}
}

// initializeSpireFactories sets the options of the spire sprite factory
func (g *Game) initializeSpireFactories() {
	_, spireHeight := topSpire.Size()

	g.spires = nil
	g.topSpireFactory = &spriteutils.SpriteFactory{
		Images: []*ebiten.Image{topSpire},
		MaxX:   screenWidth + 150,
		MinX:   screenWidth + 150,
		MaxY:   0,
		MinY:   -200,
	}

	g.bottomSpireFactory = &spriteutils.SpriteFactory{
		Images: []*ebiten.Image{bottomSpire},
		MaxX:   screenWidth + 150,
		MinX:   screenWidth + 150,
		MaxY:   screenHeight - spireHeight + 200,
		MinY:   screenHeight - spireHeight,
	}
}

// initializeAsteroidFactories sets the options of the asteroid sprite factory
func (g *Game) initializeAsteroidFactories() {
	g.asteroids = nil
	g.asteroidFactory = &spriteutils.SpriteFactory{
		Images: []*ebiten.Image{asteroid1, asteroid2, asteroid3, asteroid4},
		MaxX:   screenWidth + 100,
		MinX:   screenWidth + 100,
		MaxY:   screenHeight - 100,
		MinY:   100,
	}
}

// initializeStarFactory sets the options of the star sprite factory
func (g *Game) initializeStarFactory() {
	g.stars = nil
	g.starFactory = &spriteutils.SpriteFactory{
		Images: []*ebiten.Image{starImage},
		MaxX:   screenWidth + 100,
		MinX:   screenWidth + 100,
		MaxY:   screenHeight - 100,
		MinY:   100,
	}
}

// drawScore draws the score (distance travelled)
func (g *Game) drawScore(screen *ebiten.Image) {
	scoreStr := fmt.Sprintf("Distance: %8d m", g.distanceTravelled)
	text.Draw(screen, scoreStr, normalFont, screenWidth-(len(scoreStr)*fontSize/2), fontSize, color.White)
}

// createAsteroidExplosion creates the sprites for asteroid explosion, given an asteroid
func (g *Game) createAsteroidExplosion(asteroid *spriteutils.Sprite) *spriteutils.TransientSprite {
	return &spriteutils.TransientSprite{
		CreatedAtGameTime: time.Duration(g.frameCount) * time.Second / 60,
		LifetimeDuration:  time.Millisecond * 100,
		Sprite: &spriteutils.Sprite{
			Image:     asteroidExplosionImage,
			X:         asteroid.X,
			Y:         asteroid.Y,
			XVelocity: -g.speed,
			Rotation:  rand.Float64() * math.Pi,
		},
	}
}

// checkCollisions does all the collision handling logic
func (g *Game) checkCollisions() {
	// Ground collisions
	for _, tile := range g.topGroundTiles {
		if g.ship.IsColliding(tile) {
			g.mode = ModeGameOver
		}

		for i, asteroid := range g.asteroids {
			if asteroid.IsColliding(tile) {
				g.asteroidExplosions = append(g.asteroidExplosions, g.createAsteroidExplosion(asteroid))
				g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
			}
		}
	}

	for _, tile := range g.bottomGroundTiles {
		if g.ship.IsColliding(tile) {
			g.mode = ModeGameOver
		}

		for i, asteroid := range g.asteroids {
			if asteroid.IsColliding(tile) {
				g.asteroidExplosions = append(g.asteroidExplosions, g.createAsteroidExplosion(asteroid))
				g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
			}
		}
	}

	// spire collisions
	for _, spire := range g.spires {
		if g.ship.IsColliding(spire) {
			g.mode = ModeGameOver
		}

		for i, asteroid := range g.asteroids {
			if asteroid.IsColliding(spire) {
				g.asteroidExplosions = append(g.asteroidExplosions, g.createAsteroidExplosion(asteroid))
				g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
			}
		}
	}

	// asteroid collisions
	for i, asteroid := range g.asteroids {
		if g.shield != nil && g.shield.IsColliding(asteroid) {
			g.asteroidExplosions = append(g.asteroidExplosions, g.createAsteroidExplosion(asteroid))
			g.asteroids = append(g.asteroids[:i], g.asteroids[i+1:]...)
		}

		if g.ship.IsColliding(asteroid) {
			g.mode = ModeGameOver
		}
	}

	// star collisions
	for i, star := range g.stars {
		if g.ship.IsColliding(star) {
			g.stars = append(g.stars[:i], g.stars[i+1:]...)
			g.isBoosting = true
			g.lastBoostTime = time.Duration(g.frameCount) * time.Second / 60
			g.speed += g.boostFactor
		}
	}
}

//updateGround updates the ground positions and ensures that the ground loops properly
func (g *Game) updateGround() {
	imageWidth, imageHeight := floorImage.Size()

	for _, tile := range g.topGroundTiles {
		tile.XVelocity = -g.speed
		tile.Update()
	}

	for _, tile := range g.bottomGroundTiles {
		tile.XVelocity = -g.speed
		tile.Update()
	}

	if g.topGroundTiles[0].X <= -imageWidth {
		g.topGroundTiles = append(g.topGroundTiles[:0], g.topGroundTiles[1:]...)
		g.topGroundTiles = append(g.topGroundTiles, &spriteutils.Sprite{
			Image:     floorImage,
			X:         g.topGroundTiles[len(g.topGroundTiles)-1].X + imageWidth,
			XVelocity: -g.speed,
			Rotation:  math.Pi,
		})
	}

	if g.bottomGroundTiles[0].X <= -imageWidth {
		g.bottomGroundTiles = append(g.bottomGroundTiles[:0], g.bottomGroundTiles[1:]...)
		g.bottomGroundTiles = append(g.bottomGroundTiles, &spriteutils.Sprite{
			Image:     floorImage,
			X:         g.bottomGroundTiles[len(g.bottomGroundTiles)-1].X + imageWidth,
			Y:         screenHeight - imageHeight,
			XVelocity: -g.speed,
		})
	}
}

// updateSpires updates the spire positions and destroys out of bounds spires
func (g *Game) updateSpires() {
	temp := g.spires[:0]
	for _, spire := range g.spires {
		spire.XVelocity = -g.speed
		spire.Update()

		if spire.X > outOfBoundsX {
			temp = append(temp, spire)
		}
	}
	g.spires = temp
}

// updateAsteroids updates the asteroid positions and destroys out of bounds asteroids
func (g *Game) updateAsteroids() {
	temp := g.asteroids[:0]
	for _, asteroid := range g.asteroids {
		asteroid.Update()

		if asteroid.X > outOfBoundsX {
			temp = append(temp, asteroid)
		}
	}
	g.asteroids = temp
}

// updateStars updates the star positions and destroys out of bounds stars
func (g *Game) updateStars() {
	temp := g.stars[:0]
	for _, star := range g.stars {
		star.XVelocity = -g.speed
		star.Update()

		if star.X > outOfBoundsX {
			temp = append(temp, star)
		}
	}
	g.stars = temp
}
