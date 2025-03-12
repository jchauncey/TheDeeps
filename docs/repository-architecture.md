# Repository Architecture

## Current Implementation

The current implementation uses in-memory repositories for storing game data. These repositories are defined in the `repositories` package and include:

- `CharacterRepository`: Stores character data
- `DungeonRepository`: Stores dungeon data
- `InventoryRepository`: Stores inventory data

## Issue: Repository Sharing

We identified an issue where handlers were creating their own instances of repositories instead of using shared instances. This caused problems when trying to join a dungeon, as the character would be found in one repository instance but not in another.

### Problem

In the `DungeonHandler`, the constructor was creating new repository instances:

```go
// Old implementation
func NewDungeonHandler() *DungeonHandler {
    return &DungeonHandler{
        dungeonRepo:   repositories.NewDungeonRepository(),
        characterRepo: repositories.NewCharacterRepository(),
        mapGenerator:  game.NewMapGenerator(time.Now().UnixNano()),
    }
}
```

This meant that when a character was created using one repository instance and then tried to join a dungeon using another repository instance, the character would not be found.

### Solution

We modified the `DungeonHandler` constructor to accept repositories as parameters:

```go
// New implementation
func NewDungeonHandler(dungeonRepo *repositories.DungeonRepository, characterRepo *repositories.CharacterRepository) *DungeonHandler {
    return &DungeonHandler{
        dungeonRepo:   dungeonRepo,
        characterRepo: characterRepo,
        mapGenerator:  game.NewMapGenerator(time.Now().UnixNano()),
    }
}
```

And updated the server initialization to pass the shared repositories to the handler:

```go
// In server.go
dungeonHandler := handlers.NewDungeonHandler(dungeonRepo, characterRepo)
```

This ensures that the same repository instances are used throughout the application, preventing the "character not found" error when trying to join a dungeon.

## Future SQL Database Implementation

For a future SQL database implementation, we plan to use the following architecture:

### 1. Repository Interfaces

Define interfaces for each repository that specify the methods each repository must implement:

```go
// CharacterRepository interface
type CharacterRepository interface {
    GetByID(id string) (*models.Character, error)
    GetAll() []*models.Character
    Save(character *models.Character) error
    Delete(id string) error
}

// DungeonRepository interface
type DungeonRepository interface {
    GetByID(id string) (*models.Dungeon, error)
    GetAll() []*models.Dungeon
    Save(dungeon *models.Dungeon) error
    AddCharacterToDungeon(dungeonID string, characterID string) error
    RemoveCharacterFromDungeon(dungeonID string, characterID string) error
    GetFloor(dungeonID string, level int) (*models.Floor, error)
}

// InventoryRepository interface
type InventoryRepository interface {
    GetInventory(characterID string) ([]*models.Item, error)
    AddItem(characterID string, item *models.Item) error
    RemoveItem(characterID string, itemID string) error
    GetItem(itemID string) (*models.Item, error)
}
```

### 2. In-Memory Implementations

Implement the interfaces with in-memory storage for development and testing:

```go
// InMemoryCharacterRepository implements CharacterRepository
type InMemoryCharacterRepository struct {
    characters map[string]*models.Character
    mutex      sync.RWMutex
}

// Implement all methods...
```

### 3. SQL Implementations

Implement the interfaces with SQL storage for production:

```go
// SQLCharacterRepository implements CharacterRepository
type SQLCharacterRepository struct {
    db *sql.DB
}

// GetByID retrieves a character from the database by ID
func (r *SQLCharacterRepository) GetByID(id string) (*models.Character, error) {
    // SQL implementation
    var character models.Character
    err := r.db.QueryRow("SELECT * FROM characters WHERE id = ?", id).Scan(
        &character.ID,
        &character.Name,
        // ... other fields
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("character not found")
        }
        return nil, err
    }
    return &character, nil
}

// Implement other methods...
```

### 4. Factory Functions

Create factory functions to instantiate the appropriate repository implementation:

```go
// NewCharacterRepository creates a new character repository
func NewCharacterRepository(db *sql.DB) CharacterRepository {
    if db == nil {
        return &InMemoryCharacterRepository{
            characters: make(map[string]*models.Character),
        }
    }
    return &SQLCharacterRepository{
        db: db,
    }
}
```

### 5. Dependency Injection

Use dependency injection to pass the repositories to the handlers:

```go
// In server.go
db := initDatabase() // Returns nil in development mode
characterRepo := repositories.NewCharacterRepository(db)
dungeonRepo := repositories.NewDungeonRepository(db)
inventoryRepo := repositories.NewInventoryRepository(db)

// Create handlers
characterHandler := handlers.NewCharacterHandler(characterRepo)
dungeonHandler := handlers.NewDungeonHandler(dungeonRepo, characterRepo)
// ...
```

### 6. Database Schema

Define the database schema for each model:

```sql
CREATE TABLE characters (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    class VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    experience INT NOT NULL,
    max_hp INT NOT NULL,
    current_hp INT NOT NULL,
    max_mana INT NOT NULL,
    current_mana INT NOT NULL,
    gold INT NOT NULL,
    current_floor INT NOT NULL,
    position_x INT NOT NULL,
    position_y INT NOT NULL,
    current_dungeon VARCHAR(36),
    FOREIGN KEY (current_dungeon) REFERENCES dungeons(id)
);

CREATE TABLE dungeons (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    floors INT NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    seed BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- ... other tables
```

### 7. Transaction Support

Add transaction support for operations that affect multiple tables:

```go
// JoinDungeon handles POST /dungeons/{id}/join
func (h *DungeonHandler) JoinDungeon(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...

    // Start a transaction
    tx, err := h.db.Begin()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Add character to dungeon
    if err := h.dungeonRepo.AddCharacterToDungeonTx(tx, dungeonID, request.CharacterID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Update character
    character.CurrentFloor = 1
    character.CurrentDungeon = dungeonID
    if err := h.characterRepo.SaveTx(tx, character); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // ... rest of the code ...
}
```

## Testing

We've added tests to verify that repositories are shared correctly:

1. `TestNewDungeonHandler`: Verifies that the handler is initialized with the correct repositories
2. `TestSharedRepositories`: Verifies that the repositories are shared correctly when joining a dungeon

These tests ensure that our fix works correctly and will continue to work as we evolve the codebase.

## Benefits

This architecture provides several benefits:

1. **Separation of Concerns**: The repository interfaces define a clear contract for data access
2. **Testability**: In-memory implementations make testing easier
3. **Flexibility**: We can switch between in-memory and SQL implementations without changing the handler code
4. **Scalability**: SQL implementations will allow the game to scale to more players and persist data across server restarts
5. **Maintainability**: The code is more modular and easier to maintain 