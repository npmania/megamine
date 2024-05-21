package main

import (
	"log"
	"megamine/game"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	err := game.InitGame()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if err := ebiten.RunGame(&game.Game); err != nil {
		log.Fatal(err)
	}
}