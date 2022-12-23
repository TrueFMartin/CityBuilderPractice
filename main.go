package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/TrueFMartin/engotut/ents"
	"github.com/TrueFMartin/engotut/systems"
	"image/color"
	"log"
)

type myScene struct{}

func (*myScene) Type() string {
	return "myGame"
}

func (*myScene) Preload() {
	err := engo.Files.Load("textures/citySheet.png", "tilemap/TrafficMap.tmx")
	if err != nil {
		log.Fatal(err)
	}
}

func (*myScene) Setup(updater engo.Updater) {
	engo.Input.RegisterButton("AddCity", engo.KeyF1)
	common.SetBackground(color.White)
	world, _ := updater.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})
	kbs := common.NewKeyboardScroller(
		400, engo.DefaultHorizontalAxis, engo.DefaultVerticalAxis)
	world.AddSystem(kbs)
	//self created method, loads tile .tmx, add to render, and sets camera bounds
	ents.AddTile(world)
	world.AddSystem(&systems.CityBuildingSystem{})
	//self created method
	ents.AddHud(world)
}

func main() {
	opts := engo.RunOptions{
		Title:          "Hello World",
		Width:          800,
		Height:         800,
		StandardInputs: true,
		NotResizable:   false,
	}
	engo.Run(opts, &myScene{})
}
