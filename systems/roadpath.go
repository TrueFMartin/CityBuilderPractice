package systems

import (
	"container/list"
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"image/color"
	"log"
)

type Node struct {
	next  *Node
	prev  *Node
	point engo.Point
}
type LinkedList struct {
	head *Node
	tail *Node
}
type Path struct {
	ecs.BasicEntity
	path *list.List
}

// Used to ID locations on 2D array, empty location, city location, road location
const (
	emptyLoc = iota
	cityLoc
	roadLoc
	pathCityLoc //City that has been added to path
	pathRoadLoc //Road that has been added to path
)

// 2d Array of Y cord, and X cord of 20x20 map grid
var pathMap [20][20]*pathMapLocation

type pathMapLocation struct {
	row, col int
	locType  int
}
type PathSystem struct {
	paths       []*list.List
	world       *ecs.World
	newLocation *pathMapLocation
}

func (ps *PathSystem) New(w *ecs.World) {
	ps.world = w
	for r := 0; r < len(pathMap); r++ {
		for c := 0; c < len(pathMap); c++ {
			pathMap[r][c] = &pathMapLocation{
				row:     r,
				col:     c,
				locType: emptyLoc,
			}
		}
	}
	engo.Mailbox.Listen(UpdatePointMessageType, func(m engo.Message) {
		msg, ok := m.(UpdatePointMessage)
		if !ok {
			return
		}
		//Convert x,y point to 20x20 map positions and ints
		x, y := floatPointToPathInt(msg.point)
		var msgLocType int
		if msg.pointType == PointTypeCity {
			msgLocType = cityLoc
		} else {
			msgLocType = roadLoc
		}

		if msg.isAdding {
			ps.newLocation = pathMap[y][x]
			ps.newLocation.locType = msgLocType
		} else if !msg.isAdding && msg.pointType == PointTypeRoad {
			ps.newLocation = pathMap[y][x]
			ps.newLocation.locType = emptyLoc
			dispatchPathOutline(false, nil, y, x)
			//FIXME add check to remove rest of path
		}
		//See if new city/road added new potential path
		ps.checkForNewPaths()
		//ps.upgradeCompletedPaths()
	})
}

func (ps *PathSystem) Update(dt float32)      {}
func (ps *PathSystem) Remove(ecs.BasicEntity) {}
func (ps *PathSystem) checkForNewPaths() {
	switch ps.newLocation.locType {
	case cityLoc:
		adjacents := getNonEmptyAdjacents(ps.newLocation)
		for _, adjacent := range adjacents {
			if adjacent.locType == roadLoc {
				path := ps.AddNewPath(ps.newLocation)
				path.PushFront(adjacent)
				dispatchPathOutline(true, nil, adjacent.row, adjacent.col)
				ps.newLocation.locType = pathCityLoc
				adjacent.locType = pathRoadLoc
			}
			//FIXME if adjacent.loctype == pathedRoadLoc
		}
		//for _, v := range ps.paths {
		//	printPathList(v)
		//}
	case roadLoc:
		adjacents := getNonEmptyAdjacents(ps.newLocation)
		for _, adjacent := range adjacents {
			if adjacent.locType == cityLoc {
				path := ps.AddNewPath(adjacent)
				path.PushFront(ps.newLocation)
				dispatchPathOutline(true, nil, ps.newLocation.row, ps.newLocation.col)
				ps.newLocation.locType = pathRoadLoc
				adjacent.locType = pathCityLoc
				fmt.Println("Just hit city at ", adjacent.row, adjacent.col)
			}
			if adjacent.locType == pathRoadLoc {
				ps.pathRoadToRoad(adjacent)
				ps.newLocation.locType = pathRoadLoc
			} //FIXME ADD some visual change when path complete, helpful for debugging
		}
		//for _, v := range ps.paths {
		//	printPathList(v)
		//}

	}
	//compare each city and index
	//for cityI, city := range ps.cities {
	//	for roadI, road := range ps.roadPieces {
	//for xIndex := 0; xIndex < 20; xIndex++ {
	//	for yIndex := 0; yIndex < 20; yIndex++ {
	//		//if ps.isPointIndexed(cityI, roadI) {
	//		mapPoint := pathMap[yIndex][xIndex]
	//		if mapPoint == emptyLoc {
	//			continue
	//		}
	//		if mapPoint == pathCityLoc || mapPoint == pathRoadLoc {
	//
	//		}
	//		//if any corner of city and road are shared, create new path
	//		if engomath.Abs(city.X-road.X) <= 64 &&
	//			engomath.Abs(city.Y-road.Y) <= 64 {
	//			//create a path with the city as start, then push road to start
	//			newPath := ps.AddNewPath(city)
	//			newPath.path.PushFront(road)
	//			for i, v := range ps.paths {
	//				for e := v.path.Front(); e != nil; e = e.Next() {
	//					fmt.Println("at path ", i, " ", e.Value) // do something with e.Value
	//				}
	//			}
	//			ps.pathedCityIndexes[cityI] = true
	//			ps.pathedRoadIndexes[roadI] = true
	//		}
	//	}
	//}
}

func (ps *PathSystem) AddNewPath(start *pathMapLocation) *list.List {
	p := list.New()
	p.PushFront(start)
	ps.paths = append(ps.paths, p)
	dispatchPathOutline(true, nil, start.row, start.col)
	return p
}

func (ps *PathSystem) pathRoadToRoad(adjacentLoc *pathMapLocation) {
	for _, oldPath := range ps.paths {
		//Adding new road to front of path
		possibleLoc := validateListType(oldPath.Front())
		if possibleLoc == adjacentLoc {
			oldPath.PushFront(ps.newLocation)
			dispatchPathOutline(true, nil, ps.newLocation.row, ps.newLocation.col)
			printPathList(oldPath)
			return //FIXME if two path converge with new location, wont register
		}
		ps.checkInnerPathParts(oldPath, adjacentLoc)
		//creating new path with new road as front of new path
		//Means new road is added as offshoot(branch) of previous path
	}
	return
}

func (ps *PathSystem) checkInnerPathParts(oldPath *list.List, adjacentLoc *pathMapLocation) {
	for outerElm := oldPath.Front(); outerElm != nil; outerElm = outerElm.Next() {
		possibleLoc := validateListType(outerElm)
		if possibleLoc == adjacentLoc {
			//Start empty path list
			newPath := CopyListUntil(oldPath, outerElm)
			newPath.PushFront(ps.newLocation)
			ps.paths = append(ps.paths, newPath)
			rgba := &color.RGBA{200, 100, 0, 255}
			dispatchPathOutline(true, rgba,
				ps.newLocation.row, ps.newLocation.col)
		}
	}

}

// CopyListUntil Copies a list from list-back until element 'e' is reached
func CopyListUntil(l *list.List, elem *list.Element) *list.List {
	//Init a new empty list
	returnList := list.New()
	//listPos starts at back of list to be copied
	listPos := l.Back()
	//While we haven't reached the 'until' value,
	for listPos.Value != elem.Value {
		//Add listPos's value to the front of the returning list,
		returnList.PushFront(listPos.Value)
		//move listPos step closer to front of list
		listPos = listPos.Prev()
	}
	returnList.PushFront(listPos.Value)
	return returnList
}

func (ps *PathSystem) upgradeCompletedPaths() bool {
	for _, v := range ps.paths {
		pathFront := validateListType(v.Front())
		pathBack := validateListType(v.Back())
		if pathFront.locType == pathCityLoc && pathBack.locType == pathCityLoc {

		}
	}

	return false
}

// FIXME probably remove
func checkForAdjacentCity(location pathMapLocation) pathMapLocation {
	return pathMapLocation{}
}

// Returns a slice of pathMapLocations, each one being a filled location on 20x20 grid
// Only returns directly adjacent, non-empty, locations
func getNonEmptyAdjacents(location *pathMapLocation) (adjacents []*pathMapLocation) {
	r, c := location.row, location.col

	var (
		above *pathMapLocation
		below *pathMapLocation
		left  *pathMapLocation
		right *pathMapLocation
	)

	if r > 0 {
		above = pathMap[r-1][c]
		if above.locType != emptyLoc {
			adjacents = append(adjacents, above)
		}
	}

	if c > 0 {
		left = pathMap[r][c-1]
		if left.locType != emptyLoc {
			adjacents = append(adjacents, left)
		}
	}

	below = pathMap[r+1][c]
	if below.locType != emptyLoc {
		adjacents = append(adjacents, below)
	}

	right = pathMap[r][c+1]
	if right.locType != emptyLoc {
		adjacents = append(adjacents, right)
	}

	return
}

// Confirm type assertion of list element's value as pathMapLocation
func validateListType(element *list.Element) *pathMapLocation {
	pathLocation, ok := element.Value.(*pathMapLocation)
	if !ok {
		log.Fatal("Invalid type assertion on Validating list element type'")
	}
	return pathLocation
}

// Returns true if two pathMapLocations share same point in 20x20 grid
func areLocationsEqual(loc1, loc2 pathMapLocation) bool {
	return loc1.row == loc2.row && loc1.col == loc2.col
}

// Print each item of list
func printPathList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		pathLocation := validateListType(e)
		fmt.Println("Col: ", pathLocation.col, " Row: ", pathLocation.row, "Type: ", pathLocation.locType)
	}
}

func printMap() {
	for r := 0; r < len(pathMap); r++ {
		for c := 0; c < len(pathMap); c++ {
			fmt.Printf("%v ", pathMap[r][c].locType)
		}
		fmt.Printf("\n")
	}
}

// Takes a point, returns a position on 20x20 grid
func floatPointToPathInt(p engo.Point) (x, y int) {
	xFlo, yFlo := p.X, p.Y
	x, y = int(xFlo/64), int(yFlo/64)
	return
}

// Dispatch a message creating outline around 20x20 matrix point, nil color* defaults white border
func dispatchPathOutline(isAdding bool, color *color.RGBA, r, c int) {
	if color == nil {
		engo.Mailbox.Dispatch(UpdateOutlineMessage{
			row:      r,
			col:      c,
			isAdding: isAdding,
		})
	} else {
		engo.Mailbox.Dispatch(UpdateOutlineMessage{
			row:      r,
			col:      c,
			isAdding: isAdding,
			color:    color,
		})
	}
}
func (l *LinkedList) PushBack(n *Node) {

	if l.head == nil {
		l.head = n
		l.tail = n
	} else {
		l.tail.next = n
		l.tail = n
	}

}

type UpdatePointMessage struct {
	point     engo.Point
	pointType string
	isAdding  bool
}

const UpdatePointMessageType string = "UpdatePointMessage"

// Used for 'pointType' in UpdatePointMessage
const PointTypeCity string = "city"
const PointTypeRoad string = "road"

func (UpdatePointMessage) Type() string {
	return UpdatePointMessageType
}
