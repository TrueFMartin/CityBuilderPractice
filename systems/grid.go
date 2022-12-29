package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Grid Vertical or Horizontal grid bar
type Grid struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

// GridSystem collection of Grids (X and Y)
type GridSystem struct {
	world *ecs.World
	gridX []*Grid
	gridY []*Grid
	image common.RenderComponent
}

func (g *GridSystem) New(w *ecs.World) {
	g.world = w
	//Set image to be used by each grid, solid black tile
	g.image.Drawable = Spritesheet.Cell(898)
	//Create 20x20 Grid pattern
	for i := 1; i < 20; i++ {
		barX := Grid{BasicEntity: ecs.NewBasic()}
		barY := Grid{BasicEntity: ecs.NewBasic()}

		//At intervals of 64 pixels
		pos := float32(i) * 64
		//Set position of each bar, X then Y
		barX.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{X: 0, Y: pos},
			Width:    1260,
			Height:   1,
		}
		barY.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{X: pos, Y: 0},
			Width:    1,
			Height:   1260,
		}
		//set image of bar to black tile
		barX.RenderComponent.Drawable = g.image.Drawable
		barY.RenderComponent.Drawable = g.image.Drawable

		//Stretch black tile Horizontal for barX and vertical for barY
		barX.RenderComponent.Scale = engo.Point{X: 80, Y: .1}
		barY.RenderComponent.Scale = engo.Point{X: .1, Y: 80}

		//Set grid pattern to be above other entities
		barX.RenderComponent.SetZIndex(10)
		barY.RenderComponent.SetZIndex(10)

		//Update system w/ reference to each bar
		g.gridX = append(g.gridX, &barX)
		g.gridY = append(g.gridY, &barY)
	}

	for _, system := range g.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			//Add each barX of horizontal grid to renderer
			for _, grid := range g.gridX {
				sys.Add(&grid.BasicEntity, &grid.RenderComponent, &grid.SpaceComponent)
			}
			//Add each barY of vertical grid to renderer
			for _, grid := range g.gridY {
				sys.Add(&grid.BasicEntity, &grid.RenderComponent, &grid.SpaceComponent)
			}
		}
	}
}

func (g *GridSystem) Update(dt float32) {}

func (g *GridSystem) Remove(entity ecs.BasicEntity) {}

func GetNearestPoint(p engo.Point) engo.Point {
	x := int(p.X/64) * 64
	y := int(p.X/64) * 64
	return engo.Point{X: float32(x), Y: float32(y)}
}
