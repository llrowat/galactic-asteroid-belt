package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"log"
	"math/rand"
	"time"
)

const (
	screenWidth   = 1028
	screenHeight  = 720
	fontSize      = 24
	titleFontSize = fontSize * 1.5
	smallFontSize = fontSize / 2
)

var (
	backgroundImage        *ebiten.Image
	shipImage              *ebiten.Image
	floorImage             *ebiten.Image
	topSpire               *ebiten.Image
	bottomSpire            *ebiten.Image
	asteroid1              *ebiten.Image
	asteroid2              *ebiten.Image
	asteroid3              *ebiten.Image
	asteroid4              *ebiten.Image
	asteroidExplosionImage *ebiten.Image
	starImage              *ebiten.Image
	shieldImage            *ebiten.Image
	titleFont              font.Face
	normalFont             font.Face
	smallFont              font.Face
)

// Seed random number generator
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Initialize all images
func init() {
	var err error
	backgroundImage, _, err = ebitenutil.NewImageFromFile("assets/background.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	shipImage, _, err = ebitenutil.NewImageFromFile("assets/spaceship.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	floorImage, _, err = ebitenutil.NewImageFromFile("assets/groundDirt.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	topSpire, _, err = ebitenutil.NewImageFromFile("assets/rock-top.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	bottomSpire, _, err = ebitenutil.NewImageFromFile("assets/rock-bottom.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	asteroid1, _, err = ebitenutil.NewImageFromFile("assets/meteorBrown_big1.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	asteroid2, _, err = ebitenutil.NewImageFromFile("assets/meteorBrown_big2.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	asteroid3, _, err = ebitenutil.NewImageFromFile("assets/meteorBrown_big3.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	asteroid4, _, err = ebitenutil.NewImageFromFile("assets/meteorBrown_big4.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	asteroidExplosionImage, _, err = ebitenutil.NewImageFromFile("assets/meteorExplosion.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	starImage, _, err = ebitenutil.NewImageFromFile("assets/starGold.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	shieldImage, _, err = ebitenutil.NewImageFromFile("assets/shield.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

// Initialize fonts
func init() {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	normalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Initialize game
func newGame() *Game {
	game := &Game{}
	game.init()
	return game
}

// Entry point
func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Galactic Asteroid Belt")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
