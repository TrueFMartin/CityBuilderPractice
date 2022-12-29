package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

//var Spritesheet *common.Spritesheet

//type MouseTracker struct {
//	ecs.BasicEntity
//	common.MouseComponent
//}

type Highway struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	isVert bool
}
type HighwayEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.RenderComponent
	*common.MouseComponent
}

type HighwaySystem struct {
	world        *ecs.World
	mouseTracker MouseTracker
	entities     []HighwayEntity
	highways     []Highway
	vertTile     common.RenderComponent
	horizTile    common.RenderComponent
}

func (h *HighwaySystem) New(w *ecs.World) {
	h.world = w

	h.horizTile.Drawable = Spritesheet.Cell(716)
	h.horizTile.Scale = engo.Point{X: 3, Y: 3}
	h.vertTile.Drawable = Spritesheet.Cell(753)
	h.vertTile.Scale = engo.Point{X: 3, Y: 3}
	h.mouseTracker.BasicEntity = ecs.NewBasic()
	h.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&h.mouseTracker.BasicEntity, &h.mouseTracker.MouseComponent, nil, nil)
		}
	}
}

func (h *HighwaySystem) Update(dt float32) {
	//fmt.Println(h.mouseTracker.MouseX, " ", h.mouseTracker.MouseY)
	if engo.Input.Button("AddRoadVert").JustPressed() ||
		engo.Input.Button("AddRoadHoriz").JustPressed() {
		highway := Highway{BasicEntity: ecs.NewBasic()}
		highway.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{
				X: h.mouseTracker.MouseX,
				Y: h.mouseTracker.MouseY,
			},
			Width:  32,
			Height: 32,
		}
		isVert := engo.Input.Button("AddRoadVert").JustPressed()
		if isVert {
			highway.RenderComponent = h.vertTile
		} else {
			highway.RenderComponent = h.horizTile
		}
		for _, system := range h.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				fmt.Println(highway.SpaceComponent)
				fmt.Println(highway.RenderComponent)
				fmt.Println(highway.BasicEntity)
				sys.Add(&highway.BasicEntity, &highway.RenderComponent, &highway.SpaceComponent)
			}
		}
	}
	//for _, system := range h.world.Systems() {
	//	switch sys := system.(type) {
	//	case *common.RenderSystem:
	//		for _, v := range h.entities {
	//			sys.Add(v.BasicEntity, v.RenderComponent, v.SpaceComponent)
	//		}
	//	}
	//}
}

func (h *HighwaySystem) Remove(ecs.BasicEntity) {}
