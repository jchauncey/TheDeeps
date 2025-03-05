package mapgen

import (
	"math/rand"
)

// GenerateDungeon creates a new dungeon with the specified configuration
func GenerateDungeon(config *DungeonConfig) *Dungeon {
	if config == nil {
		config = NewDefaultConfig()
	}

	dungeon := &Dungeon{
		Floors: make([]Floor, config.TotalFloors),
	}

	for i := range make([]struct{}, config.TotalFloors) {
		// Keep generating floors until we get a valid one
		for {
			floor := createFloor(config, i)
			if validateFloor(floor) {
				dungeon.Floors[i] = *floor
				break
			}
		}
	}

	return dungeon
}

func createFloor(config *DungeonConfig, floorNum int) *Floor {
	level := &Level{
		Width:  config.Width,
		Height: config.Height,
		Tiles:  make([][]Tile, config.Height),
	}

	// Initialize all tiles as walls with varied symbols
	for y := range make([]struct{}, config.Height) {
		level.Tiles[y] = make([]Tile, config.Width)
		for x := range make([]struct{}, config.Width) {
			symbol := '|'
			if y == 0 || y == config.Height-1 {
				symbol = '-'
			}
			if (y == 0 || y == config.Height-1) && (x == 0 || x == config.Width-1) {
				symbol = '+'
			}
			level.Tiles[y][x] = Tile{Symbol: symbol, Walkable: false, Type: TileWall}
		}
	}

	// Generate rooms with different styles based on floor number
	var rooms []Room
	for len(rooms) < 3 {
		rooms = generateRooms(config, floorNum)
	}

	// Connect rooms with hallways using minimum spanning tree
	connectRoomsWithMST(level, rooms)

	// Add room features based on floor number
	addRoomFeatures(level, rooms, floorNum)

	// Find the two most distant rooms for entrance and exit
	entrance, exit := findEntranceExitRooms(rooms)

	// Place entrance and exit with special surrounding tiles
	placeEntranceExit(level, entrance, exit)

	return &Floor{
		Level:    *level,
		Entrance: Position{X: entrance.X + entrance.Width/2, Y: entrance.Y + entrance.Height/2},
		Exit:     Position{X: exit.X + exit.Width/2, Y: exit.Y + exit.Height/2},
	}
}

// validateFloor ensures the floor is properly connected
func validateFloor(floor *Floor) bool {
	// Check if we have both entrance and exit
	if floor.Entrance == (Position{}) || floor.Exit == (Position{}) {
		return false
	}

	// Use flood fill to verify connectivity
	visited := make(map[Position]bool)
	floodFill(&floor.Level, floor.Entrance, visited)

	// Check if exit is reachable
	return visited[floor.Exit]
}

// floodFill marks all reachable positions from start
func floodFill(level *Level, start Position, visited map[Position]bool) {
	if start.X < 0 || start.X >= level.Width ||
		start.Y < 0 || start.Y >= level.Height ||
		!level.Tiles[start.Y][start.X].Walkable ||
		visited[start] {
		return
	}

	visited[start] = true

	// Check all adjacent tiles
	directions := []Position{
		{X: 0, Y: -1}, // up
		{X: 0, Y: 1},  // down
		{X: -1, Y: 0}, // left
		{X: 1, Y: 0},  // right
	}

	for _, dir := range directions {
		next := Position{X: start.X + dir.X, Y: start.Y + dir.Y}
		floodFill(level, next, visited)
	}
}

func findEntranceExitRooms(rooms []Room) (entrance, exit Room) {
	// Find the two rooms with maximum Manhattan distance
	maxDist := 0
	for i, r1 := range rooms {
		for j, r2 := range rooms {
			if i == j {
				continue
			}
			dist := manhattanDist(
				Position{X: r1.X + r1.Width/2, Y: r1.Y + r1.Height/2},
				Position{X: r2.X + r2.Width/2, Y: r2.Y + r2.Height/2},
			)
			if dist > maxDist {
				maxDist = dist
				entrance = r1
				exit = r2
			}
		}
	}
	return entrance, exit
}

func manhattanDist(p1, p2 Position) int {
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func connectRoomsWithMST(level *Level, rooms []Room) {
	if len(rooms) == 0 {
		return
	}

	// Create a list of all possible connections between rooms
	type edge struct {
		room1, room2 int
		distance     int
	}

	edges := make([]edge, 0)
	for i := 0; i < len(rooms); i++ {
		for j := i + 1; j < len(rooms); j++ {
			dist := manhattanDist(
				Position{X: rooms[i].X + rooms[i].Width/2, Y: rooms[i].Y + rooms[i].Height/2},
				Position{X: rooms[j].X + rooms[j].Width/2, Y: rooms[j].Y + rooms[j].Height/2},
			)
			edges = append(edges, edge{i, j, dist})
		}
	}

	// Sort edges by distance
	for i := 0; i < len(edges)-1; i++ {
		for j := i + 1; j < len(edges); j++ {
			if edges[i].distance > edges[j].distance {
				edges[i], edges[j] = edges[j], edges[i]
			}
		}
	}

	// Create disjoint set for MST
	parent := make([]int, len(rooms))
	for i := range parent {
		parent[i] = i
	}

	// Find with path compression
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	// Union by rank
	union := func(x, y int) {
		parent[find(x)] = find(y)
	}

	// Build MST and create hallways
	for _, e := range edges {
		if find(e.room1) != find(e.room2) {
			union(e.room1, e.room2)
			start := Position{
				X: rooms[e.room1].X + rooms[e.room1].Width/2,
				Y: rooms[e.room1].Y + rooms[e.room1].Height/2,
			}
			end := Position{
				X: rooms[e.room2].X + rooms[e.room2].Width/2,
				Y: rooms[e.room2].Y + rooms[e.room2].Height/2,
			}
			createHallway(level, start, end)
		}
	}

	// Carve out the rooms after creating hallways
	for _, room := range rooms {
		for y := room.Y; y < room.Y+room.Height; y++ {
			for x := room.X; x < room.X+room.Width; x++ {
				level.Tiles[y][x] = Tile{Symbol: '.', Walkable: true, Type: TileFloor}
			}
		}
	}

	// Add some random additional connections for variety (about 15% more)
	extraConnections := len(rooms) / 7
	for i := 0; i < extraConnections; i++ {
		r1 := rand.Intn(len(rooms))
		r2 := rand.Intn(len(rooms))
		if r1 != r2 {
			start := Position{
				X: rooms[r1].X + rooms[r1].Width/2,
				Y: rooms[r1].Y + rooms[r1].Height/2,
			}
			end := Position{
				X: rooms[r2].X + rooms[r2].Width/2,
				Y: rooms[r2].Y + rooms[r2].Height/2,
			}
			createHallway(level, start, end)
		}
	}
}

func generateRooms(config *DungeonConfig, floorNum int) []Room {
	rooms := make([]Room, 0)
	maxAttempts := 200 // Increased attempts to try harder to place rooms

	// Calculate target number of rooms based on floor size
	// Assuming minimum room size with 1 tile padding between rooms
	effectiveWidth := config.Width - 2 // Account for walls
	effectiveHeight := config.Height - 2
	minRoomWithPadding := config.MinRoomSize + 2 // Account for padding
	targetRooms := (effectiveWidth / minRoomWithPadding) * (effectiveHeight / minRoomWithPadding)

	// Cap target rooms to a reasonable number
	if targetRooms > 20 {
		targetRooms = 20
	} else if targetRooms < 5 {
		targetRooms = 5
	}

	// Adjust room sizes based on floor number, but keep them smaller to fit more
	minRoomSize := config.MinRoomSize
	maxRoomSize := config.MinRoomSize + 3 + (floorNum % 3) // Keep max size relatively small
	padding := 1

	// Try to place rooms until we hit target or run out of attempts
	attempts := 0
	for len(rooms) < targetRooms && attempts < maxAttempts {
		// Randomly choose room size, biased towards smaller rooms
		roomWidth := minRoomSize + rand.Intn(maxRoomSize-minRoomSize+1)
		roomHeight := minRoomSize + rand.Intn(maxRoomSize-minRoomSize+1)

		// Try to place room in random location
		x := 1 + rand.Intn(config.Width-roomWidth-2)
		y := 1 + rand.Intn(config.Height-roomHeight-2)

		newRoom := Room{X: x, Y: y, Width: roomWidth, Height: roomHeight}

		// Check if room overlaps with existing rooms
		if !roomOverlapsAny(newRoom, rooms, padding) {
			rooms = append(rooms, newRoom)

			// After successfully placing a room, try to place a connected room nearby
			if rand.Float32() < 0.5 { // 50% chance to try adjacent room
				adjacentAttempts := 0
				for adjacentAttempts < 10 {
					// Try to place room adjacent to the one we just placed
					side := rand.Intn(4) // 0: top, 1: right, 2: bottom, 3: left
					adjWidth := minRoomSize + rand.Intn(maxRoomSize-minRoomSize+1)
					adjHeight := minRoomSize + rand.Intn(maxRoomSize-minRoomSize+1)

					var adjX, adjY int
					switch side {
					case 0: // top
						adjX = x + rand.Intn(roomWidth-2) - (adjWidth - 2)
						adjY = y - adjHeight - 1
					case 1: // right
						adjX = x + roomWidth + 1
						adjY = y + rand.Intn(roomHeight-2) - (adjHeight - 2)
					case 2: // bottom
						adjX = x + rand.Intn(roomWidth-2) - (adjWidth - 2)
						adjY = y + roomHeight + 1
					case 3: // left
						adjX = x - adjWidth - 1
						adjY = y + rand.Intn(roomHeight-2) - (adjHeight - 2)
					}

					// Verify room is within bounds
					if adjX > 0 && adjY > 0 &&
						adjX+adjWidth < config.Width-1 &&
						adjY+adjHeight < config.Height-1 {

						adjRoom := Room{X: adjX, Y: adjY, Width: adjWidth, Height: adjHeight}
						if !roomOverlapsAny(adjRoom, rooms, padding) {
							rooms = append(rooms, adjRoom)
							break
						}
					}
					adjacentAttempts++
				}
			}
		}
		attempts++
	}

	// If we don't have minimum required rooms, try one last time with smaller sizes
	if len(rooms) < 3 {
		minRoomSize = config.MinRoomSize - 1
		maxRoomSize = config.MinRoomSize + 2
		for attempts := 0; attempts < 50 && len(rooms) < 3; attempts++ {
			roomWidth := minRoomSize + rand.Intn(maxRoomSize-minRoomSize+1)
			roomHeight := minRoomSize + rand.Intn(maxRoomSize-minRoomSize+1)
			x := 1 + rand.Intn(config.Width-roomWidth-2)
			y := 1 + rand.Intn(config.Height-roomHeight-2)

			newRoom := Room{X: x, Y: y, Width: roomWidth, Height: roomHeight}
			if !roomOverlapsAny(newRoom, rooms, padding) {
				rooms = append(rooms, newRoom)
			}
		}
	}

	return rooms
}

func roomsOverlapWithPadding(r1, r2 Room, padding int) bool {
	return !((r1.X+r1.Width+padding) < r2.X ||
		(r2.X+r2.Width+padding) < r1.X ||
		(r1.Y+r1.Height+padding) < r2.Y ||
		(r2.Y+r2.Height+padding) < r1.Y)
}

func createHallway(level *Level, start, end Position) {
	// Randomly choose between L-shaped and S-shaped hallways
	if rand.Float32() < 0.3 { // 30% chance for S-shaped
		// Create S-shaped hallway
		midX := (start.X + end.X) / 2

		// First horizontal section to midpoint
		x := start.X
		for x != midX {
			if x < midX {
				x++
			} else {
				x--
			}
			level.Tiles[start.Y][x] = Tile{Symbol: '.', Walkable: true, Type: TileHallway}
		}

		// Vertical section at midpoint
		y := start.Y
		for y != end.Y {
			if y < end.Y {
				y++
			} else {
				y--
			}
			level.Tiles[y][midX] = Tile{Symbol: '.', Walkable: true, Type: TileHallway}
		}

		// Final horizontal section from midpoint to end
		x = midX
		for x != end.X {
			if x < end.X {
				x++
			} else {
				x--
			}
			level.Tiles[end.Y][x] = Tile{Symbol: '.', Walkable: true, Type: TileHallway}
		}
	} else {
		// L-shaped hallway (existing code)
		x := start.X
		for x != end.X {
			if x < end.X {
				x++
			} else {
				x--
			}
			level.Tiles[start.Y][x] = Tile{Symbol: '.', Walkable: true, Type: TileHallway}
		}

		y := start.Y
		for y != end.Y {
			if y < end.Y {
				y++
			} else {
				y--
			}
			level.Tiles[y][x] = Tile{Symbol: '.', Walkable: true, Type: TileHallway}
		}
	}
}

func addRoomFeatures(level *Level, rooms []Room, floorNum int) {
	for _, room := range rooms {
		// Basic room carving
		for y := room.Y; y < room.Y+room.Height; y++ {
			for x := room.X; x < room.X+room.Width; x++ {
				level.Tiles[y][x] = Tile{Symbol: '.', Walkable: true, Type: TileFloor}
			}
		}

		// Add features based on room size and floor number
		if room.Width >= 8 && room.Height >= 8 {
			// Add pillars in larger rooms
			if floorNum%4 == 0 {
				addPillars(level, room)
			} else if floorNum%4 == 1 {
				addWaterFeature(level, room)
			} else if floorNum%4 == 2 {
				addRubble(level, room)
			}
		}
	}
}

func addPillars(level *Level, room Room) {
	// Add pillars in corners of large rooms
	pillarPositions := []Position{
		{X: room.X + 1, Y: room.Y + 1},
		{X: room.X + room.Width - 2, Y: room.Y + 1},
		{X: room.X + 1, Y: room.Y + room.Height - 2},
		{X: room.X + room.Width - 2, Y: room.Y + room.Height - 2},
	}

	for _, pos := range pillarPositions {
		level.Tiles[pos.Y][pos.X] = Tile{Symbol: 'O', Walkable: false, Type: TilePillar}
	}
}

func addWaterFeature(level *Level, room Room) {
	// Add a small pond or stream
	centerX := room.X + room.Width/2
	centerY := room.Y + room.Height/2

	for y := centerY - 1; y <= centerY+1; y++ {
		for x := centerX - 2; x <= centerX+2; x++ {
			if y >= room.Y && y < room.Y+room.Height &&
				x >= room.X && x < room.X+room.Width {
				level.Tiles[y][x] = Tile{Symbol: '~', Walkable: false, Type: TileWater}
			}
		}
	}
}

func addRubble(level *Level, room Room) {
	// Add random rubble piles
	numRubble := rand.Intn(4) + 2
	for i := 0; i < numRubble; i++ {
		x := room.X + rand.Intn(room.Width-2) + 1
		y := room.Y + rand.Intn(room.Height-2) + 1
		level.Tiles[y][x] = Tile{Symbol: '%', Walkable: false, Type: TileRubble}
	}
}

func roomOverlapsAny(newRoom Room, rooms []Room, padding int) bool {
	for _, room := range rooms {
		if roomsOverlapWithPadding(newRoom, room, padding) {
			return true
		}
	}
	return false
}

func placeEntranceExit(level *Level, entrance, exit Room) {
	// Place entrance with surrounding floor tiles
	entrancePos := Position{X: entrance.X + entrance.Width/2, Y: entrance.Y + entrance.Height/2}
	exitPos := Position{X: exit.X + exit.Width/2, Y: exit.Y + exit.Height/2}

	level.Tiles[entrancePos.Y][entrancePos.X] = Tile{Symbol: '<', Walkable: true, Type: TileEntrance}
	level.Tiles[exitPos.Y][exitPos.X] = Tile{Symbol: '>', Walkable: true, Type: TileExit}

	// Add floor tiles around entrance/exit
	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			if x == 0 && y == 0 {
				continue
			}
			makeFloorTile(level, entrancePos.Y+y, entrancePos.X+x)
			makeFloorTile(level, exitPos.Y+y, exitPos.X+x)
		}
	}
}

func makeFloorTile(level *Level, y, x int) {
	if y >= 0 && y < level.Height && x >= 0 && x < level.Width {
		level.Tiles[y][x] = Tile{Symbol: '.', Walkable: true, Type: TileFloor}
	}
}
