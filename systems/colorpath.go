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
var colorMap [20][20]Outline

// OutlineSystem Groups all outlines together, decides which display
type OutlineSystem struct {
	world      *ecs.World
	outlines   []Outline
	outlineRec Outline
}

const (
	recBorderWidth float32 = 2
	recSize        float32 = 64
)

func (ol *OutlineSystem) New(w *ecs.World) {
	ol.world = w
	baseRect := common.Rectangle{
		BorderWidth: recBorderWidth,
		BorderColor: color.RGBA{0, 255, 255, 255}, //Create teal boundary on rec
	}
	ol.outlineRec = Outline{} //All grid rect's will be based off this,
	ol.outlineRec.SpaceComponent.Width = recSize
	ol.outlineRec.SpaceComponent.Height = recSize
	ol.outlineRec.Drawable = baseRect
	ol.outlineRec.Color = color.Transparent //Make inner rectangle transparent
	ol.outlineRec.Hidden = true             // start hidden
	ol.outlineRec.SetZIndex(11)             // draw ontop of regular grid

	for _, system := range ol.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			//Fill matrix w/ hidden, transparent, white-bordered rectangles
			for r := 0; r < len(colorMap); r++ {
				for c := 0; c < len(colorMap); c++ {
					colorMap[r][c] = ol.outlineRec
					colorMap[r][c].BasicEntity = ecs.NewBasic()
					colorMap[r][c].Position = engo.Point{
						X: float32(c * 64),
						Y: float32(r * 64),
					}
					colorMap[r][c].row = r
					colorMap[r][c].col = c

					//FIXME May be better to only add matrix locations that have been activated
					//Add that matrix point to render system(currently hidden)
					sys.Add(&colorMap[r][c].BasicEntity,
						&colorMap[r][c].RenderComponent,
						&colorMap[r][c].SpaceComponent)
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
			baseRect := common.Rectangle{
				BorderWidth: recBorderWidth,
				BorderColor: msg.color, //Create teal boundary on rec
			}
			//ol.outlineRec.Color = color.Transparent //Make inner rectangle transparent
			colorMap[r][c].Drawable = baseRect
		}
		if msg.isAdding {
			colorMap[r][c].Hidden = false
		} else {
			colorMap[r][c].Hidden = true
		}

	})
}
func (ol *OutlineSystem) Update(dt float32) {
	if int(dt)%10 < 5 {
		colorMap[1][1].BorderColor = color.RGBA{0, 0, 0, 255}
	} else {
		colorMap[1][1].BorderColor = color.RGBA{255, 255, 255, 255}
	}
}
func (ol *OutlineSystem) Remove(ecs.BasicEntity) {}

type UpdateOutlineMessage struct {
	row, col int
	isAdding bool
	color    *color.RGBA
}

const UpdateOutlineMessageType string = "UpdateOutlineMessage"

func (UpdateOutlineMessage) Type() string {
	return UpdateOutlineMessageType
}
