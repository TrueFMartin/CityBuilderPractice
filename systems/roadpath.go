package systems

import (
	"container/list"
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"image/color"
	"log"
)

type Path struct {
	ecs.BasicEntity
	path *list.List
}

// Used to ID locations on 2D array, empty location, city location, road location
const (
	// emptyLoc
	emptyLoc = iota
	townLoc
	roadLoc
	pathRoadLoc //Road that has been added to path
	pathTownLoc //town that has been added to path
	pathCityLoc
	pathMetroLoc
)

// 2d Array of Y cord, and X cord of 20x20 matrix grid
var pathMatrix [20][20]*pathMatrixLocation

type pathMatrixLocation struct {
	row, col         int
	locType          int
	pathsBelongingTo []int
}
type PathSystem struct {
	paths       []*list.List
	world       *ecs.World
	newLocation *pathMatrixLocation
}

func (ps *PathSystem) New(w *ecs.World) {
	ps.world = w
	for r := 0; r < len(pathMatrix); r++ {
		for c := 0; c < len(pathMatrix); c++ {
			pathMatrix[r][c] = &pathMatrixLocation{
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
		//Convert x,y point to 20x20 matrix positions and ints
		x, y := floatPointToPathInt(msg.point)
		var msgLocType int
		if msg.pointType == PointTypeTown {
			msgLocType = townLoc
		} else {
			msgLocType = roadLoc
		}

		if msg.isAdding {
			ps.newLocation = pathMatrix[y][x]
			ps.newLocation.locType = msgLocType
		} else if !msg.isAdding && msg.pointType == PointTypeRoad {
			ps.newLocation = pathMatrix[y][x]
			ps.newLocation.locType = emptyLoc
			dispatchPathOutline(false, nil, y, x)
			//FIXME add check to remove rest of path
		}
		//See if new city/road added new potential path
		ps.checkForNewPaths()
		//ps.upgradeConnectedLocations()
	})
}

func (ps *PathSystem) Update(dt float32)      {}
func (ps *PathSystem) Remove(ecs.BasicEntity) {}
func (ps *PathSystem) checkForNewPaths() {
	switch ps.newLocation.locType {
	case townLoc:
		adjacents := getNonEmptyAdjacents(ps.newLocation)
		for _, adjacent := range adjacents {
			switch adjacent.locType {
			case pathRoadLoc:
				//Check each path that the adjacent path'ed road belongs to
				//Then upgrade each attached city/town
				ps.upgradeConnectedLocations(adjacent)
			case roadLoc:
				//Connects a lone road with a lone city, creates a new path
				ps.connectTownAndRoad(adjacent, ps.newLocation)
			}
		}

	case roadLoc:
		//Get adjacent grid points with items, path'ed roads will be first
		adjacents := getNonEmptyAdjacents(ps.newLocation)
		for _, adjacent := range adjacents {
			switch adjacent.locType {
			case pathRoadLoc: // Add new road to front of adjacent path
				ps.pathRoadToPath(adjacent) //connect paths and color new location
				ps.newLocation.locType = pathRoadLoc
			case roadLoc: // Combine two lone roads
				newPath, index := ps.AddNewPath(adjacent) //front of path will be the adjacent road
				newPath.PushFront(ps.newLocation)
				//Add index(ps.path's path index) to newlocation's pathsBelongingTo
				ps.newLocation.pathsBelongingTo = append(ps.newLocation.pathsBelongingTo, index)
				adjacent.locType = pathRoadLoc //label both roads as a path'ed road
				ps.newLocation.locType = pathRoadLoc
				//color both
				dispatchPathOutline(true, &roadPathColor, adjacent.row, adjacent.col)
				dispatchPathOutline(true, &roadPathColor, ps.newLocation.row, ps.newLocation.col)
			case townLoc: //new road is adjacent to an unpath'ed town
				//If a previous loop of the adjacent's hasn't added
				//the new road into a previous path
				if ps.newLocation.locType == roadLoc {
					ps.connectTownAndRoad(ps.newLocation, adjacent)
				} else { //meaning that road has been made into a path'ed road

				}
			default:

			}
			//if adjacent.locType == pathRoadLoc {
			//	ps.pathRoadToPath(adjacent)
			//	ps.newLocation.locType = pathRoadLoc
			//} else {
			//	if ps.newLocation.locType != pathRoadLoc {
			//		path := ps.AddNewPath(adjacent)
			//		path.PushFront(ps.newLocation)
			//	}
			//	dispatchPathOutline(true, nil, ps.newLocation.row, ps.newLocation.col)
			//	dispatchPathOutline(true, &townColor, adjacent.row, adjacent.col)
			//	ps.newLocation.locType = pathRoadLoc
			//	adjacent.locType = pathTownLoc
			//	ps.upgradeConnectedLocations()
			//}
		}
	}
}

// AddNewPath Creates a path starting at 'start',
// returns the path's list and index of list in ps.paths
func (ps *PathSystem) AddNewPath(start *pathMatrixLocation) (*list.List, int) {
	path := list.New()
	path.PushFront(start)
	pathIndex := len(ps.paths)
	//Update locations belongingTo with new path index
	start.pathsBelongingTo = append(start.pathsBelongingTo, pathIndex)
	ps.paths = append(ps.paths, path)
	//dispatchPathOutline(true, nil, start.row, start.col)
	return path, pathIndex
}

// connectTownAndRoad is a helper method to connect a single lone road to a new city
func (ps *PathSystem) connectTownAndRoad(road, city *pathMatrixLocation) {
	path, index := ps.AddNewPath(city)
	path.PushFront(road)
	//update  road's pathsBelongingTo to include new path's index
	road.pathsBelongingTo = append(road.pathsBelongingTo, index)
	dispatchPathOutline(true, &townColor, city.row, city.col)
	dispatchPathOutline(true, &roadPathColor, road.row, road.col)
	city.locType = pathTownLoc
	road.locType = pathRoadLoc
}

func (ps *PathSystem) pathRoadToPath(adjacentRoad *pathMatrixLocation) {
	for _, pathIndex := range adjacentRoad.pathsBelongingTo {
		front := validateListType(ps.paths[pathIndex].Front())
		if front == adjacentRoad {
			ps.paths[pathIndex].PushFront(ps.newLocation)
			ps.newLocation.pathsBelongingTo = append(ps.newLocation.pathsBelongingTo, pathIndex)
			dispatchPathOutline(true, &roadPathColor, ps.newLocation.row, ps.newLocation.col)
			return
		} //FIXME Currenlty I need to create a new path for every front and back of all related paths
	} //FIXME This seems really complicated, maybe make it so roads can only be started at cities
	//FIXME If I do this, never-mined defintely do this
	//First check front of each existing path
	for i, oldPath := range ps.paths {
		//Adding new road to front of paths
		possibleLoc := validateListType(oldPath.Front())
		if possibleLoc == adjacentRoad {
			oldPath.PushFront(ps.newLocation)
			ps.newLocation.pathsBelongingTo = append(ps.newLocation.pathsBelongingTo, i)
			dispatchPathOutline(true, &roadPathColor, ps.newLocation.row, ps.newLocation.col)
			return //FIXME if two path converge with new location, wont register
		}
	}
	//Since adjacent location isn't at front of all current paths, check inside of paths
	for _, oldPath := range ps.paths {
		// if adjacent location is an inner part of old path
		if ps.checkInnerPathParts(oldPath, adjacentRoad) {
			//return the index of that path's
			return
		}

	}
	return
}

func (ps *PathSystem) checkInnerPathParts(oldPath *list.List, adjacentLoc *pathMatrixLocation) bool {
	//Starting at front of path(Most recently added), step to next until match found
	for outerElm := oldPath.Front(); outerElm != nil; outerElm = outerElm.Next() {
		possibleLoc := validateListType(outerElm)
		if possibleLoc == adjacentLoc {
			length := len(ps.paths) //Used to add ".pathsBelongingTo" index #
			//Start empty path list
			newPath := CopyListUntil(oldPath, outerElm)
			newPath.PushFront(ps.newLocation)
			//For each element in this divergent path,
			//update its "pathsBelongingTo" to include new path index
			for e := oldPath.Back(); e != nil; e = e.Prev() {
				location := validateListType(e)
				location.pathsBelongingTo = append(location.pathsBelongingTo, length)
			}
			//add this new divergent path to list of paths
			ps.paths = append(ps.paths, newPath)
			//color the new element to show seperate from rest of path
			dispatchPathOutline(true, &divergentPathColor,
				ps.newLocation.row, ps.newLocation.col)
			return true
		}
	}
	return false
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

// upgradeConnectedLocations checks all paths that a point belongs to,
// then checks all cities at the front/back of the paths.
// This is only called when a new city is connected. Towns will be upgraded to cities, cities to metros
func (ps *PathSystem) upgradeConnectedLocations(connectLocation *pathMatrixLocation) bool {
	locationsToUpdate := make([]*pathMatrixLocation, 0)
	//Get all locations at end of paths that are a town or greater
	for _, pathsIndex := range connectLocation.pathsBelongingTo {
		path := ps.paths[pathsIndex]
		pathFront := validateListType(path.Front())
		pathBack := validateListType(path.Back())
		if pathFront.locType >= pathTownLoc {
			locationsToUpdate = append(locationsToUpdate, pathFront)
		}
		if pathBack.locType >= pathTownLoc {
			locationsToUpdate = append(locationsToUpdate, pathBack)
		}
	}
	for _, l := range locationsToUpdate {
		//If city is metro, can't upgrade further
		upgradeALocation(l)
	}
	return false //FIXME remove return value
}

// For a locations that need to be updated, send out a hudText message and a color message
func upgradeALocation(l *pathMatrixLocation) {
	if l.locType != pathMetroLoc { //FIXME may be cleaner to just do "if city, do this","if town..."
		startLocType := CityType(l.locType)
		endLocType := CityType(l.locType + 1)
		if !(endLocType == pathMetroLoc || endLocType == pathCityLoc) {
			log.Fatal("upgradeConnectedLocations, upgrade is not a city or metro"+
				"instead is trying to upgrade to ", endLocType)
		}
		var upgradeColor *color.RGBA
		switch endLocType {
		case pathCityLoc:
			upgradeColor = &cityColor
		case pathMetroLoc:
			upgradeColor = &metroColor
		default:
			log.Fatal("upgradeConnectedLocations, invalid type for upgradeColor assignment")
		}
		engo.Mailbox.Dispatch(CityUpdateMessage{
			Old: startLocType,
			New: endLocType,
		})
		engo.Mailbox.Dispatch(UpdateOutlineMessage{
			row:      l.row,
			col:      l.col,
			isAdding: true,
			color:    upgradeColor,
		})
	}
}

// FIXME probably remove
func checkForAdjacentCity(location pathMatrixLocation) pathMatrixLocation {
	return pathMatrixLocation{}
}

// Returns a slice of pathMatrixLocations, each one being a filled location on 20x20 grid
// Only returns directly adjacent, non-empty, locations
// Already Path'ed ROADS are added to front of return slice
func getNonEmptyAdjacents(location *pathMatrixLocation) (adjacents []*pathMatrixLocation) {
	r, c := location.row, location.col

	var (
		above *pathMatrixLocation
		below *pathMatrixLocation
		left  *pathMatrixLocation
		right *pathMatrixLocation
	)
	tempAdjacents := make([]*pathMatrixLocation, 0)
	if r > 0 {
		above = pathMatrix[r-1][c]
		if above.locType != emptyLoc {
			tempAdjacents = append(tempAdjacents, above)
		}
	}

	if c > 0 {
		left = pathMatrix[r][c-1]
		if left.locType != emptyLoc {
			tempAdjacents = append(tempAdjacents, left)
		}
	}

	below = pathMatrix[r+1][c]
	if below.locType != emptyLoc {
		tempAdjacents = append(tempAdjacents, below)
	}

	right = pathMatrix[r][c+1]
	if right.locType != emptyLoc {
		tempAdjacents = append(tempAdjacents, right)
	}
	//place roads first
	for _, v := range tempAdjacents {
		if v.locType == pathRoadLoc {
			adjacents = append(adjacents, v)
		}
	}
	//add the rest of adjacents to
	for _, v := range tempAdjacents {
		if v.locType != pathRoadLoc {
			adjacents = append(adjacents, v)
		}
	}
	fmt.Println(adjacents)
	return
}

// Confirm type assertion of list element's value as pathMatrixLocation
func validateListType(element *list.Element) *pathMatrixLocation {
	pathLocation, ok := element.Value.(*pathMatrixLocation)
	if !ok {
		log.Fatal("Invalid type assertion on Validating list element type'")
	}
	return pathLocation
}

// Returns true if two pathMatrixLocations share same point in 20x20 grid
func areLocationsEqual(loc1, loc2 pathMatrixLocation) bool {
	return loc1.row == loc2.row && loc1.col == loc2.col
}

// Print each item of list
func printPathList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		pathLocation := validateListType(e)
		fmt.Println("Col: ", pathLocation.col, " Row: ", pathLocation.row, "Type: ", pathLocation.locType)
	}
}

func printMatrix() {
	for r := 0; r < len(pathMatrix); r++ {
		for c := 0; c < len(pathMatrix); c++ {
			fmt.Printf("%v ", pathMatrix[r][c].locType)
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

// Dispatch a message creating outline around 20x20 matrix point, nil color* defaults teal border
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

type UpdatePointMessage struct {
	point     engo.Point
	pointType string
	isAdding  bool
}

const UpdatePointMessageType string = "UpdatePointMessage"

// PointTypeTown Used for 'pointType' in UpdatePointMessage
const PointTypeTown string = "town"

// PointTypeRoad Used for 'pointType' in UpdatePointMessage
const PointTypeRoad string = "road"

func (UpdatePointMessage) Type() string {
	return UpdatePointMessageType
}
