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
	world         *ecs.World
	mouseTracker  MouseTracker
	entities      []HighwayEntity
	highways      []*Highway
	vertTile      common.RenderComponent
	horizTile     common.RenderComponent
	toRemoveIndex int
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
	//If user inputs F1 or F2, add/remove road
	if engo.Input.Button("AddRoadVert").JustPressed() ||
		engo.Input.Button("AddRoadHoriz").JustPressed() {
		//get nearest grid intersection point to mouse position on F1/F2 input
		position := GetNearestPoint(engo.Point{
			X: h.mouseTracker.MouseX,
			Y: h.mouseTracker.MouseY,
		})
		//check if already road at position, if not add road, else remove
		possiblePositionIndex := h.isHighwayPresent(position)
		if possiblePositionIndex == -1 { //-1 means no road present at position
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
			//update money Amount from cost of road
			engo.Mailbox.Dispatch(RoadCostMessage{Amount: -50})
			//let path manager know about new road
			engo.Mailbox.Dispatch(UpdatePointMessage{
				point:     highway.Position,
				pointType: PointTypeRoad,
				isAdding:  true,
			})
		} else { //Position already filled by highway
			ent := h.highways[possiblePositionIndex].BasicEntity
			h.toRemoveIndex = possiblePositionIndex
			h.Remove(ent)
		}
	}
}

func (h *HighwaySystem) Remove(ent ecs.BasicEntity) {
	//remove entity from renderer
	for _, system := range h.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(ent)
		}
	}
	//remove highway from slice of highways
	engo.Mailbox.Dispatch(UpdatePointMessage{
		point:     h.highways[h.toRemoveIndex].Position,
		pointType: PointTypeRoad,
		isAdding:  false,
	})
	//remove highway from h.highways slice
	h.highways = append(h.highways[:h.toRemoveIndex],
		h.highways[h.toRemoveIndex+1:]...)
	//update money Amount from selling road
	engo.Mailbox.Dispatch(RoadCostMessage{Amount: 50})

}

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

type RoadCostMessage struct {
	Amount int
}

const RoadCostMessageType string = "RoadCostMessage"

func (RoadCostMessage) Type() string {
	return RoadCostMessageType
}
