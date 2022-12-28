package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"math/rand"
	"time"
)

var Spritesheet *common.Spritesheet

var cities = [...][12]int{
	{99, 100, 101,
		454, 269, 455,
		415, 195, 416,
		452, 306, 453,
	},
	{99, 100, 101,
		268, 269, 270,
		268, 269, 270,
		305, 306, 307,
	},
	{75, 76, 77,
		446, 261, 447,
		446, 261, 447,
		444, 298, 445,
	},
	{75, 76, 77,
		407, 187, 408,
		407, 187, 408,
		444, 298, 445,
	},
	{75, 76, 77,
		186, 150, 188,
		186, 150, 188,
		297, 191, 299,
	},
	{83, 84, 85,
		413, 228, 414,
		411, 191, 412,
		448, 302, 449,
	},
	{83, 84, 85,
		227, 228, 229,
		190, 191, 192,
		301, 302, 303,
	},
	{91, 92, 93,
		241, 242, 243,
		278, 279, 280,
		945, 946, 947,
	},
	{91, 92, 93,
		241, 242, 243,
		278, 279, 280,
		945, 803, 947,
	},
	{91, 92, 93,
		238, 239, 240,
		238, 239, 240,
		312, 313, 314,
	},
}

type City struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	timeBuilt, timeAlive float32
}

type MouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}
type CityBuildingSystem struct {
	world              *ecs.World
	mouseTracker       MouseTracker
	usedTiles          []int
	elapsed, buildTime float32
	built              int
	worldTime          float32
}

func (cb *CityBuildingSystem) New(w *ecs.World) {
	cb.world = w
	fmt.Println("CityBuildingSystem was added to Scene")

	cb.mouseTracker.BasicEntity = ecs.NewBasic()
	cb.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&cb.mouseTracker.BasicEntity, &cb.mouseTracker.MouseComponent, nil, nil)
		}
	}
	Spritesheet = common.NewSpritesheetWithBorderFromFile(
		"textures/citySheet.png",
		16,
		16,
		1,
		1,
	)
	rand.Seed(time.Now().UnixNano())
}

func (cb *CityBuildingSystem) Update(dt float32) {
	//add city at randomized times, prog. faster
	cb.elapsed += dt
	cb.worldTime += dt
	if cb.elapsed >= cb.buildTime {
		cb.generateCity()
		cb.elapsed = 0
		cb.updateBuildTime()
		cb.built++
	}
	// This is for adding citys on mouse postion w/ F1 press
	//#FIXME Removed for now, may add back later
	//if engo.Input.Button("AddCity").JustPressed() {
	//	fmt.Println("F1 Pressed")
	//	fmt.Printf("\nWindow W: %f, Game W: %f, CanvasW: %f",
	//		engo.WindowWidth(), engo.GameWidth(), engo.CanvasWidth())
	//	city := City{BasicEntity: ecs.NewBasic()}
	//	city.SpaceComponent = common.SpaceComponent{
	//		Position: engo.Point{
	//			X: cb.mouseTracker.MouseX,
	//			Y: cb.mouseTracker.MouseY,
	//		},
	//		Width:  30,
	//		Height: 64,
	//	}
	//	texture, err := common.LoadedSprite("textures/citySheet.png")
	//	if err != nil {
	//		log.Println("Load Texture failed: " + err.Error())
	//	}
	//	city.RenderComponent = common.RenderComponent{
	//		Scale:    engo.Point{X: 0.1, Y: 0.1},
	//		Drawable: texture,
	//	}
	//	for _, system := range cb.world.Systems() {
	//		switch sys := system.(type) {
	//		case *common.RenderSystem:
	//			sys.Add(&city.BasicEntity, &city.RenderComponent, &city.SpaceComponent)
	//		}
	//	}
	//}
}

func (*CityBuildingSystem) Remove(ecs.BasicEntity) {}

func (cb *CityBuildingSystem) generateCity() {
	x, y := rand.Intn(18), rand.Intn(18)
	t := x + y*18

	for cb.isTileUsed(t) {
		if len(cb.usedTiles) > 300 {
			break
		}
		x, y = rand.Intn(18), rand.Intn(18)
		t = x + y*18
	}
	cb.usedTiles = append(cb.usedTiles, t)

	city := rand.Intn(len(cities))
	cityTiles := make([]*City, 0)
	for i := 0; i < 3; i++ { //i is x axis for city building
		for j := 0; j < 4; j++ { //j is y-axis for city building
			tile := &City{BasicEntity: ecs.NewBasic()}
			tile.SpaceComponent.Position = engo.Point{
				X: float32(((x+1)*64)+8) + float32(i*16),
				Y: float32((y+1)*64) + float32(j*16),
			}
			tile.RenderComponent.Drawable = Spritesheet.Cell(cities[city][i+(3*j)])
			tile.RenderComponent.SetZIndex(1)
			cityTiles = append(cityTiles, tile)
		}
	}

	for _, system := range cb.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range cityTiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
	buildTimeTxt := fmt.Sprintf("Built at Gamehour: %.0f", cb.worldTime)
	engo.Mailbox.Dispatch(HUDTextMessage{
		BasicEntity: ecs.NewBasic(),
		SpaceComponent: common.SpaceComponent{
			Position: engo.Point{X: float32((x + 1) * 64), Y: float32((y + 1) * 64)},
			Width:    64,
			Height:   64,
		},
		MouseComponent: common.MouseComponent{},
		Line1:          "Town",
		Line2:          buildTimeTxt,
		Line3:          "A town generates",
		Line4:          "$100 per day.",
	})

	engo.Mailbox.Dispatch(CityUpdateMessage{
		New: CityTypeNew,
	})
}

func (cb *CityBuildingSystem) isTileUsed(tile int) bool {
	for _, t := range cb.usedTiles {
		if tile == t {
			return true
		}
	}
	return false
}

func (cb *CityBuildingSystem) updateBuildTime() {
	switch {
	case cb.built < 2:
		// 10 to 15 seconds
		cb.buildTime = 5*rand.Float32() + 10
	case cb.built < 5:
		// 60 to 90 seconds
		cb.buildTime = 30*rand.Float32() + 60
	case cb.built < 10:
		// 30 to 90 seconds
		cb.buildTime = 60*rand.Float32() + 30
	case cb.built < 20:
		// 30 to 65 seconds
		cb.buildTime = 35*rand.Float32() + 30
	case cb.built < 25:
		// 30 to 60 seconds
		cb.buildTime = 30*rand.Float32() + 30
	default:
		// 20 to 40 seconds
		cb.buildTime = 20*rand.Float32() + 20
	}
}
