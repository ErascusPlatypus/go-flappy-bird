package main

import (
	"embed"
	"pro16_flappybird/files"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/**
var assets embed.FS

func main() {
	files.Init(assets)
	files.InitFeatures()

	g := files.NewGame()
	// ebiten.SetWindowSize()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
