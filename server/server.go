package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/handlers"
	"github.com/jchauncey/TheDeeps/server/log"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// Server represents the game server
type Server struct {
	router           *mux.Router
	characterRepo    *repositories.CharacterRepository
	dungeonRepo      *repositories.DungeonRepository
	inventoryRepo    *repositories.InventoryRepository
	gameManager      *game.GameManager
	characterHandler *handlers.CharacterHandler
	dungeonHandler   *handlers.DungeonHandler
	combatHandler    *handlers.CombatHandler
	inventoryHandler *handlers.InventoryHandler
}

// NewServer creates a new server instance
func NewServer() *Server {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()
	inventoryRepo := repositories.NewInventoryRepository()

	// Create game manager
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create handlers
	characterHandler := handlers.NewCharacterHandler(characterRepo)
	dungeonHandler := handlers.NewDungeonHandler(dungeonRepo, characterRepo)
	combatHandler := handlers.NewCombatHandler(characterRepo, dungeonRepo, gameManager)
	inventoryHandler := handlers.NewInventoryHandler(characterRepo, inventoryRepo)

	// Create server
	server := &Server{
		router:           mux.NewRouter(),
		characterRepo:    characterRepo,
		dungeonRepo:      dungeonRepo,
		inventoryRepo:    inventoryRepo,
		gameManager:      gameManager,
		characterHandler: characterHandler,
		dungeonHandler:   dungeonHandler,
		combatHandler:    combatHandler,
		inventoryHandler: inventoryHandler,
	}

	// Setup routes
	server.SetupRoutes()

	return server
}

// SetupRoutes configures the server routes
func (s *Server) SetupRoutes() {
	// Character routes
	s.router.HandleFunc("/characters", s.characterHandler.GetCharacters).Methods("GET")
	s.router.HandleFunc("/characters", s.characterHandler.CreateCharacter).Methods("POST")
	s.router.HandleFunc("/characters/{id}", s.characterHandler.GetCharacter).Methods("GET")
	s.router.HandleFunc("/characters/{id}", s.characterHandler.DeleteCharacter).Methods("DELETE")
	s.router.HandleFunc("/characters/{id}/save", s.characterHandler.SaveCharacter).Methods("POST")
	s.router.HandleFunc("/characters/{id}/floor", s.characterHandler.GetCharacterFloor).Methods("GET")

	// Dungeon routes
	s.router.HandleFunc("/dungeons", s.dungeonHandler.GetDungeons).Methods("GET")
	s.router.HandleFunc("/dungeons", s.dungeonHandler.CreateDungeon).Methods("POST")
	s.router.HandleFunc("/dungeons/{id}/join", s.dungeonHandler.JoinDungeon).Methods("POST")
	s.router.HandleFunc("/dungeons/{id}/floor/{level}", s.dungeonHandler.GetFloor).Methods("GET")
	s.router.HandleFunc("/api/dungeons/{id}/floors/{floorNumber}", s.dungeonHandler.GetFloorByNumber).Methods("GET")
	s.router.HandleFunc("/test/room", s.dungeonHandler.GenerateTestRoom).Methods("GET")

	// Combat routes
	s.router.HandleFunc("/characters/{id}/combat", s.combatHandler.GetCombatState).Methods("GET")
	s.router.HandleFunc("/ws/combat", s.combatHandler.HandleCombat).Methods("GET")

	// Inventory routes
	s.inventoryHandler.RegisterRoutes(s.router)

	// WebSocket route for real-time game updates
	s.router.HandleFunc("/ws/game", s.gameManager.HandleConnection)
}

// Start starts the server on the specified address
func (s *Server) Start(addr string) error {
	log.Info("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
