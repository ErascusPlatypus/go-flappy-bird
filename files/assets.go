package files

import (
	"bytes"
	"embed"
	"image"
	"log"
	_ "image/png"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

var Assets embed.FS
var PipeSprite *ebiten.Image
var BackgroundImage *ebiten.Image
var CoinSprites [] *ebiten.Image

func Init(fs embed.FS) {
	Assets = fs
}

func InitFeatures() {
	PipeSprite = loadAsset("assets/pipe-green.png")
	BackgroundImage = loadAsset("assets/background-night.png")
	CoinSprites = loadAssets("assets/coins_*.png")
}

func loadAsset(path string) *ebiten.Image {
	data, err := Assets.ReadFile(path)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func loadAssets(path string) [] *ebiten.Image {
	files, err := filepath.Glob(path)

	if err != nil {
		log.Fatal(err)
	}

	var res [] *ebiten.Image 

	for _, f := range files {
		res = append(res, loadAsset(f))
	}

	return res 
}