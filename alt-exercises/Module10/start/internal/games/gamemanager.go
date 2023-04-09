package games

import (
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type GetSession struct {
	sessionId   string
	gameSession chan *GameSession
}

type GameManager struct {
	register chan *GameSession

	unRegister chan *GameSession

	getSessionChannel chan *GetSession

	activeGames map[string]*GameSession

	game Game
}

func CreateGameManager() GameManager {

	manager := GameManager{
		register:          make(chan *GameSession),
		unRegister:        make(chan *GameSession),
		getSessionChannel: make(chan *GetSession),
		activeGames:       make(map[string]*GameSession),
	}

	manager.loadGameConfig()

	return manager
}

func (manager *GameManager) loadGameConfig() (id string, name string, description string) {
	dir, _ := os.Getwd()
	viper.SetConfigName("app")
	// Set the path to look for the configurations file
	viper.AddConfigPath(dir + "/../configs")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	manager.game = Game{
		ID:          viper.GetString("GAME_ID"),
		Name:        viper.GetString("GAME_NAME"),
		Description: viper.GetString("GAME_DESCRIPTION"),
	}

	return
}

//CreateNewGameSession type
func (manager *GameManager) CreateNewGameSession() *GameSession {
	game := &GameSession{
		SendToGame:  make(chan *GameMsg),
		Register:    make(chan *Player),
		UnRegister:  make(chan *Player),
		Players:     make(map[string]*Player),
		ID:          uuid.New().String(),
		gameManager: manager,
	}

	manager.register <- game

	return game
}

func (manager *GameManager) GetSessionByID(session string) *GameSession {
	if manager.getSessionChannel != nil {
		ch := make(chan *GameSession)
		manager.getSessionChannel <- &GetSession{
			sessionId:   session,
			gameSession: ch,
		}
		return <-ch
	}
	return nil
}

func (manager *GameManager) Run() {

	for {
		select {
		case session := <-manager.register:
			manager.activeGames[session.ID] = session
		case session := <-manager.unRegister:
			delete(manager.activeGames, session.ID)
		case getsession := <-manager.getSessionChannel:
			getsession.gameSession <- manager.activeGames[getsession.sessionId]
		}
	}
}

func (manager *GameManager) GetGame(gameId string) (Game, error) {

	if manager.game.ID == gameId {
		return manager.game, nil
	}
	return Game{}, errors.New("Game Is Not Supported")
}
