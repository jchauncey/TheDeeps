package game

import (
	"log"
	"math/rand"

	"github.com/jchauncey/TheDeeps/server/models"
)

// MobSpawner handles the spawning of mobs on the map
type MobSpawner struct {
	game *Game
}

// NewMobSpawner creates a new mob spawner
func NewMobSpawner(game *Game) *MobSpawner {
	return &MobSpawner{
		game: game,
	}
}

// SpawnMobsOnFloor spawns mobs on the specified floor
func (ms *MobSpawner) SpawnMobsOnFloor(floorLevel int) {
	if floorLevel < 0 || floorLevel >= len(ms.game.Dungeon.Floors) {
		log.Printf("Invalid floor level for spawning mobs: %d", floorLevel)
		return
	}

	floor := ms.game.Dungeon.Floors[floorLevel]

	// Clear existing entities
	floor.Entities = []models.Entity{}

	// Generate new entities
	floor.Entities = ms.generateEntities(floor, floorLevel+1)

	log.Printf("Spawned %d mobs on floor %d", len(floor.Entities), floorLevel+1)
}

// generateEntities generates random entities for the floor
func (ms *MobSpawner) generateEntities(floor *models.Floor, level int) []models.Entity {
	// Scale number of entities with floor level and size
	numEntities := 5 + level*2 + len(floor.Rooms)/2
	entities := make([]models.Entity, 0, numEntities)

	// Get available mob types for this floor level
	availableMobTypes := models.GetMobsForFloorLevel(level)
	if len(availableMobTypes) == 0 {
		// Fallback to basic mobs if none are available
		availableMobTypes = []models.MobType{models.MobRatman, models.MobGoblin, models.MobSkeleton}
	}

	// Generate entities
	for i := 0; i < numEntities; i++ {
		// Pick a random room (excluding the first room which is the player's starting point)
		roomIndex := 0
		if len(floor.Rooms) > 1 {
			roomIndex = 1 + rand.Intn(len(floor.Rooms)-1)
		}
		room := floor.Rooms[roomIndex]

		// Pick a random position within the room
		x := room.X + rand.Intn(room.Width)
		y := room.Y + rand.Intn(room.Height)

		// Pick a random mob type from available types
		mobTypeIndex := rand.Intn(len(availableMobTypes))
		mobType := availableMobTypes[mobTypeIndex]

		// Determine difficulty based on floor level
		difficulty := models.GetRandomDifficulty(level)

		// Create mob instance
		mobPosition := models.Position{X: x, Y: y}
		mobInstance := models.CreateMobInstance(mobType, difficulty, level, mobPosition)

		// Convert to Entity for the floor
		entity := models.Entity{
			ID:        mobInstance.ID,
			Type:      string(mobInstance.Type),
			Name:      mobInstance.Name,
			Position:  mobInstance.Position,
			Health:    mobInstance.Health,
			MaxHealth: mobInstance.MaxHealth,
			Status:    mobInstance.Status,
		}

		entities = append(entities, entity)
	}

	return entities
}

// SpawnMobsOnAllFloors spawns mobs on all floors
func (ms *MobSpawner) SpawnMobsOnAllFloors() {
	for i := range ms.game.Dungeon.Floors {
		ms.SpawnMobsOnFloor(i)
	}
}

// SpawnAdditionalMobs spawns additional mobs on the current floor
// This can be used for events, reinforcements, etc.
func (ms *MobSpawner) SpawnAdditionalMobs(count int) {
	currentFloor := ms.game.Dungeon.CurrentFloor
	if currentFloor < 0 || currentFloor >= len(ms.game.Dungeon.Floors) {
		return
	}

	floor := ms.game.Dungeon.Floors[currentFloor]
	floorLevel := currentFloor + 1

	// Get available mob types for this floor level
	availableMobTypes := models.GetMobsForFloorLevel(floorLevel)
	if len(availableMobTypes) == 0 {
		// Fallback to basic mobs if none are available
		availableMobTypes = []models.MobType{models.MobRatman, models.MobGoblin, models.MobSkeleton}
	}

	// Find rooms that are not the player's current room
	var availableRooms []models.Room
	playerRoom := ms.getPlayerRoom()

	for _, room := range floor.Rooms {
		if playerRoom == nil || !roomsOverlap(&room, playerRoom, 0) {
			availableRooms = append(availableRooms, room)
		}
	}

	if len(availableRooms) == 0 {
		return // No available rooms to spawn mobs
	}

	// Generate additional entities
	for i := 0; i < count; i++ {
		// Pick a random room
		roomIndex := rand.Intn(len(availableRooms))
		room := availableRooms[roomIndex]

		// Pick a random position within the room
		x := room.X + rand.Intn(room.Width)
		y := room.Y + rand.Intn(room.Height)

		// Pick a random mob type from available types
		mobTypeIndex := rand.Intn(len(availableMobTypes))
		mobType := availableMobTypes[mobTypeIndex]

		// Determine difficulty based on floor level
		difficulty := models.GetRandomDifficulty(floorLevel)

		// Create mob instance
		mobPosition := models.Position{X: x, Y: y}
		mobInstance := models.CreateMobInstance(mobType, difficulty, floorLevel, mobPosition)

		// Convert to Entity for the floor
		entity := models.Entity{
			ID:        mobInstance.ID,
			Type:      string(mobInstance.Type),
			Name:      mobInstance.Name,
			Position:  mobInstance.Position,
			Health:    mobInstance.Health,
			MaxHealth: mobInstance.MaxHealth,
			Status:    mobInstance.Status,
		}

		floor.Entities = append(floor.Entities, entity)
	}

	log.Printf("Spawned %d additional mobs on floor %d", count, floorLevel)
}

// getPlayerRoom returns the room that contains the player
func (ms *MobSpawner) getPlayerRoom() *models.Room {
	currentFloor := ms.game.Dungeon.CurrentFloor
	if currentFloor < 0 || currentFloor >= len(ms.game.Dungeon.Floors) {
		return nil
	}

	floor := ms.game.Dungeon.Floors[currentFloor]
	playerPos := ms.game.Dungeon.PlayerPosition

	for _, room := range floor.Rooms {
		if room.Contains(playerPos.X, playerPos.Y) {
			return &room
		}
	}

	return nil
}

// roomsOverlap checks if two rooms overlap
func roomsOverlap(r1, r2 *models.Room, minDistance int) bool {
	return !(r1.X+r1.Width+minDistance <= r2.X || r2.X+r2.Width+minDistance <= r1.X ||
		r1.Y+r1.Height+minDistance <= r2.Y || r2.Y+r2.Height+minDistance <= r1.Y)
}
