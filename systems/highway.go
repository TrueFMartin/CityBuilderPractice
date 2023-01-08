package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Highway struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	row, col int
}
type matrixLocation struct {
	row, col int
}
type HighwaySystem struct {
	world                 *ecs.World
	mouseTracker          MouseTracker
	highwayMatrix         [20][20]*Highway
	highwayImage          common.RenderComponent
	locationToRemove      matrixLocation
	common.SpaceComponent //for mouse space
	availableFunds        int
}

func (h *HighwaySystem) New(w *ecs.World) {
	h.world = w

	h.highwayImage.Drawable = Spritesheet.Cell(898) //753
	h.highwayImage.Scale = engo.Point{X: 64 / 16, Y: 64 / 16}
	h.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    1216,
		Height:   1216,
	}
	h.mouseTracker.BasicEntity = ecs.NewBasic()
	h.mouseTracker.MouseComponent = common.MouseComponent{}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&h.mouseTracker.BasicEntity, &h.mouseTracker.MouseComponent, &h.SpaceComponent, nil)
		}
	}
	engo.Mailbox.Listen(HUDMoneyMessageType, func(m engo.Message) {
		msg, ok := m.(HUDMoneyMessage)
		if !ok {
			return
		}
		h.availableFunds = msg.CurrentAmount
	})
}

func (h *HighwaySystem) Update(dt float32) {
	//If user inputs F1 or F2, add/remove road
	//h.mouseTracker.
	if h.mouseTracker.Clicked {

		//get the nearest grid intersection point from mouse position
		position := GetNearestPoint(engo.Point{
			X: h.mouseTracker.MouseX,
			Y: h.mouseTracker.MouseY,
		})
		c, r := floatPointToPathInt(position)

		//check if already road at position, if not: add road, else remove
		if h.highwayMatrix[r][c] == nil { //means there is no pointer to a highway at matrix point
			//If user has enough money to build road
			if h.availableFunds >= 50 { //FIXME add "Insuff. funds available" pop-up
				highway := &Highway{BasicEntity: ecs.NewBasic()}
				highway.SpaceComponent = common.SpaceComponent{
					Position: position,
					Width:    64,
					Height:   64,
				}
				highway.RenderComponent = h.highwayImage
				highway.row = r
				highway.col = c
				tempPathLocation := pathMatrix[r][c]
				adjacentLocations := getNonEmptyAdjacents(tempPathLocation)
				//FIXME if road is later removed breaking path, this still returns true
				//If there is a road/city adjacent to this location, add road
				if len(adjacentLocations) > 0 {
					h.highwayMatrix[r][c] = highway

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
				} else { //Meaning there were no adjacent locations allowing a road to be built
					highway = nil
					fmt.Println("No adjacent road/city, unable to build here")
				}
			}

		} else { //Position already filled by highway
			ent := h.highwayMatrix[r][c].BasicEntity
			h.locationToRemove = matrixLocation{r, c}
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
	r, c := h.locationToRemove.row, h.locationToRemove.col
	engo.Mailbox.Dispatch(UpdatePointMessage{
		point:     h.highwayMatrix[r][c].Position,
		pointType: PointTypeRoad,
		isAdding:  false,
	})

	//remove highway from matrix of slice
	h.highwayMatrix[r][c] = nil

	//update money Amount from selling road
	engo.Mailbox.Dispatch(RoadCostMessage{Amount: 50})
}

type RoadCostMessage struct {
	Amount int
}

const RoadCostMessageType string = "RoadCostMessage"

func (RoadCostMessage) Type() string {
	return RoadCostMessageType
}
