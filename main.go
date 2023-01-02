package main

import (
	"bytes"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/TrueFMartin/engotut/ents"
	"github.com/TrueFMartin/engotut/systems"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"image/color"
	"log"
)

type myScene struct{}

func (*myScene) Type() string {
	return "myGame"
}

func (*myScene) Preload() {
	err := engo.Files.Load("textures/citySheet.png", "tilemap/TrafficMap.tmx")
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	if err != nil {
		log.Fatal(err)
	}
}

func (*myScene) Setup(updater engo.Updater) {
	engo.Input.RegisterButton("AddRoadVert", engo.KeyF1)
	engo.Input.RegisterButton("AddRoadHoriz", engo.KeyF2)

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
	world.AddSystem(&systems.GridSystem{})
	world.AddSystem(&systems.HUDTextSystem{})
	world.AddSystem(&systems.MoneySystem{})
	world.AddSystem(&systems.HighwaySystem{})
	world.AddSystem(&systems.PathSystem{})
	world.AddSystem(&systems.OutlineSystem{})

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
