package files

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

var Assets embed.FS
var PipeSprite *ebiten.Image
var BackgroundImage *ebiten.Image

func Init(fs embed.FS) {
	Assets = fs
}

func InitFeatures() {
	PipeSprite = loadAsset("assets/pipe-green.png")
	BackgroundImage = loadAsset("assets/background-night.png")
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
