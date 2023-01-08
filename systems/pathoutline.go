package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"image/color"
)

// Outline A double outline(left and right or top and bottom)
// to paths based on being complete/incomplete
type Outline struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	common.Rectangle
	row, col int
}

// 20x20 matrix/grid/map of where each
var outlineMatrix [20][20]Outline

// OutlineSystem Groups all outlines together, decides which display
type OutlineSystem struct {
	world      *ecs.World
	outlineRec Outline //Just a base rectangle that individual outlines can copy
}

type Teal struct {
	r, g, b, a uint8
}

var (
	pathlessTownColor  = color.RGBA{0, 255, 255, 255} //teal
	roadPathColor      = color.RGBA{240, 255, 0, 255} //yellow
	divergentPathColor = color.RGBA{220, 0, 220, 255} // purple
	townColor          = color.RGBA{0, 255, 0, 255}   //green
	cityColor          = color.RGBA{0, 0, 255, 255}   //blue
	metroColor         = color.RGBA{255, 0, 0, 255}   //red
)

const (
	recBorderWidth float32 = 2
	recSize        float32 = 64
)

func (ol *OutlineSystem) New(w *ecs.World) {
	ol.world = w
	baseRect := common.Rectangle{
		BorderWidth: recBorderWidth,
		BorderColor: pathlessTownColor, //Create teal boundary on rec
	}
	ol.outlineRec = Outline{} //All grid rect's will be based off this,
	ol.outlineRec.SpaceComponent.Width = recSize
	ol.outlineRec.SpaceComponent.Height = recSize
	ol.outlineRec.Drawable = baseRect
	ol.outlineRec.Color = color.Transparent //Make inner rectangle transparent
	ol.outlineRec.Hidden = true             // start hidden
	ol.outlineRec.SetZIndex(11)             // draw on top of regular grid

	//Create and add each rectangle of matrix to rendersystem
	for _, system := range ol.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			//Fill matrix w/ hidden, transparent, white-bordered rectangles
			for r := 0; r < len(outlineMatrix); r++ {
				for c := 0; c < len(outlineMatrix); c++ {
					outlineMatrix[r][c] = ol.outlineRec
					outlineMatrix[r][c].BasicEntity = ecs.NewBasic()
					outlineMatrix[r][c].Position = engo.Point{
						X: float32(c * 64),
						Y: float32(r * 64),
					}
					outlineMatrix[r][c].row = r
					outlineMatrix[r][c].col = c

					//FIXME May be better to only add matrix locations that have been activated
					//Add that matrix point to render system(currently hidden)
					sys.Add(&outlineMatrix[r][c].BasicEntity,
						&outlineMatrix[r][c].RenderComponent,
						&outlineMatrix[r][c].SpaceComponent)
				}
			}
		}
	}
	//Add message system to change which outlines are hidden/displayed
	engo.Mailbox.Listen(UpdateOutlineMessageType, func(m engo.Message) {
		msg, ok := m.(UpdateOutlineMessage)
		if !ok {
			return
		}
		r, c := msg.row, msg.col
		if msg.color != nil {
			tempRect := common.Rectangle{
				BorderWidth: recBorderWidth,
				BorderColor: msg.color, //Change color of rect's border
			}
			outlineMatrix[r][c].Drawable = tempRect
		}

		if msg.isAdding {
			outlineMatrix[r][c].Hidden = false
		} else {
			outlineMatrix[r][c].Hidden = true
		}

	})
}

func (ol *OutlineSystem) Update(dt float32) {}

func (ol *OutlineSystem) Remove(ecs.BasicEntity) {}

// UpdateOutlineMessage Is dispatched whenever pathSystem adds/removes a path
type UpdateOutlineMessage struct {
	row, col int
	isAdding bool
	color    *color.RGBA
}

const UpdateOutlineMessageType string = "UpdateOutlineMessage"

func (UpdateOutlineMessage) Type() string {
	return UpdateOutlineMessageType
}
