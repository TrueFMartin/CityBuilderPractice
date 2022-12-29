package systems

import (
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
	highways     []*Highway
	vertTile     common.RenderComponent
	horizTile    common.RenderComponent
}

func (h *HighwaySystem) New(w *ecs.World) {
	h.world = w

	h.horizTile.Drawable = Spritesheet.Cell(716)
	h.horizTile.Scale = engo.Point{X: 64 / 16, Y: 64 / 16}
	h.vertTile.Drawable = Spritesheet.Cell(753)
	h.vertTile.Scale = engo.Point{X: 64 / 16, Y: 64 / 16}
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
		position := GetNearestPoint(engo.Point{
			X: h.mouseTracker.MouseX,
			Y: h.mouseTracker.MouseY,
		})

		possiblePositionIndex := h.isHighwayPresent(position)
		if possiblePositionIndex == -1 {
			highway := Highway{BasicEntity: ecs.NewBasic()}
			highway.SpaceComponent = common.SpaceComponent{
				Position: position,
				Width:    64,
				Height:   64,
			}
			isVert := engo.Input.Button("AddRoadVert").JustPressed()
			if isVert {
				highway.RenderComponent = h.vertTile
			} else {
				highway.RenderComponent = h.horizTile
			}
			h.highways = append(h.highways, &highway)
			for _, system := range h.world.Systems() {
				switch sys := system.(type) {
				case *common.RenderSystem:
					sys.Add(&highway.BasicEntity, &highway.RenderComponent, &highway.SpaceComponent)
				}
			}
		} else { //Position already filled by highway
			ent := h.highways[possiblePositionIndex].BasicEntity
			//remove entity from renderer
			for _, system := range h.world.Systems() {
				switch sys := system.(type) {
				case *common.RenderSystem:
					sys.Remove(ent)
				}
			}
			//remove highway from slice of highways
			h.highways = append(h.highways[:possiblePositionIndex],
				h.highways[possiblePositionIndex+1:]...)
		}
	}
}

func (h *HighwaySystem) Remove(ecs.BasicEntity) {}

// checks each highway for matching point, if found, returns index, else -1
func (h *HighwaySystem) isHighwayPresent(possibleP engo.Point) int {
	for i, currentHighways := range h.highways {
		if currentHighways.Position.X == possibleP.X &&
			currentHighways.Position.Y == possibleP.Y {
			return i
		}
	}
	return -1
}
