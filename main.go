package main

import (
	"log"
	"os"
	"time"

	"github.com/jroimartin/gocui"
)

func main() {
	file, err := os.OpenFile("exec.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	g, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	worldX := 120
	worldY := 100

	game := GenerateGame(worldY, worldX)
	d := NewGameWindowManager(game, g)

	done := make(chan struct{})
	go loopDisplay(d, &done)
	go loopUpdateState(game)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func loopDisplay(d *GameWindowManager, done *chan struct{}) {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()

	count := 0
	for {
		select {
		case <-*done:
			return
		case <-ticker.C:
			count++
			count %= 10

			if count == 0 {
				d.Update()
			}
		}
	}
}

func loopUpdateState(game *Game) {
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	count := 10
	for {
		<-ticker.C

		if count == 0 {
			game.Tick()
			count = 10
		} else {
			count--
		}

	}
}
