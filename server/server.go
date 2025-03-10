package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jchauncey/TheDeeps/server/game"
	"github.com/jchauncey/TheDeeps/server/handlers"
	"github.com/jchauncey/TheDeeps/server/repositories"
)

// Server represents the game server
type Server struct {
	router           *mux.Router
	characterRepo    *repositories.CharacterRepository
	dungeonRepo      *repositories.DungeonRepository
	gameManager      *game.GameManager
	characterHandler *handlers.CharacterHandler
	dungeonHandler   *handlers.DungeonHandler
	combatHandler    *handlers.CombatHandler
}

// NewServer creates a new game server
func NewServer() *Server {
	// Create repositories
	characterRepo := repositories.NewCharacterRepository()
	dungeonRepo := repositories.NewDungeonRepository()

	// Create game manager
	gameManager := game.NewGameManager(characterRepo, dungeonRepo)

	// Create handlers
	characterHandler := handlers.NewCharacterHandler()
	dungeonHandler := handlers.NewDungeonHandler()
	combatHandler := handlers.NewCombatHandler(characterRepo, dungeonRepo, gameManager)

	// Create router
	router := mux.NewRouter()

	return &Server{
		router:           router,
		characterRepo:    characterRepo,
		dungeonRepo:      dungeonRepo,
		gameManager:      gameManager,
		characterHandler: characterHandler,
		dungeonHandler:   dungeonHandler,
		combatHandler:    combatHandler,
	}
}

// SetupRoutes sets up the server routes
func (s *Server) SetupRoutes() {
	// Character routes
	s.router.HandleFunc("/characters", s.characterHandler.GetCharacters).Methods("GET")
	s.router.HandleFunc("/characters/{id}", s.characterHandler.GetCharacter).Methods("GET")
	s.router.HandleFunc("/characters", s.characterHandler.CreateCharacter).Methods("POST")
	s.router.HandleFunc("/characters/{id}", s.characterHandler.DeleteCharacter).Methods("DELETE")
	s.router.HandleFunc("/characters/{id}/save", s.characterHandler.SaveCharacter).Methods("POST")
	s.router.HandleFunc("/characters/{id}/floor", s.characterHandler.GetCharacterFloor).Methods("GET")
	s.router.HandleFunc("/characters/{id}/combat", s.combatHandler.GetCombatState).Methods("GET")

	// Dungeon routes
	s.router.HandleFunc("/dungeons", s.dungeonHandler.GetDungeons).Methods("GET")
	s.router.HandleFunc("/dungeons", s.dungeonHandler.CreateDungeon).Methods("POST")
	s.router.HandleFunc("/dungeons/{id}/join", s.dungeonHandler.JoinDungeon).Methods("POST")
	s.router.HandleFunc("/dungeons/{id}/floor/{level}", s.dungeonHandler.GetFloor).Methods("GET")

	// WebSocket routes
	s.router.HandleFunc("/ws/game", s.gameManager.HandleConnection)
	s.router.HandleFunc("/ws/combat", s.combatHandler.HandleCombat)
}

// Start starts the server
func (s *Server) Start(addr string) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
