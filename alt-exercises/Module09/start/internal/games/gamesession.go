package games

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxGameStartTime = (30 * time.Minute)
)

type GameSession struct {
	SendToGame      chan *GameMsg
	Register        chan *Player
	UnRegister      chan *Player
	Players         map[string]*Player
	ID              string
	InitialGameData string
	gameManager     *GameManager
}

func (gameSession *GameSession) CreateNewPlayer(conn *websocket.Conn, id uint, name string, email string) *Player {

	player := &Player{
		ID:          id,
		Name:        name,
		Email:       email,
		Conn:        nil,
		RecvMsgChan: nil,
		GameSession: gameSession,
	}
	return player
}

func (gameSession *GameSession) sendMsgToPlayers(gameMsg *GameMsg) {
	//iterate over the players
	//only send for connected users
	//skip sending to the one who send the message (compare email) , unless it is ment for all which means empty Player{}

}

//this method should be called to add the players to the session
// when the session recieves the STAER_GAME message from the player
func (gameSession *GameSession) addUsersToSession(players []Player) {

}

//also should be called to se the initial game data in the game session
// when the session recieves the STAER_GAME message from the player
func (gameSession *GameSession) setInitData(data string) {

}

//remove the users
//unregister the session
func (gameSession *GameSession) cleanGameSession() {

}

func (gameSession *GameSession) removeUser(player *Player) {
	player.Stop()
	delete(gameSession.Players, player.Email)
}

//the main game session loop
//register , unregister player , forward messages to all users , and dropping the sesson after maxGameStartTime
//if not all users are connectes
func (gameSession *GameSession) Run() {

	timer := time.NewTimer(maxGameStartTime)
	defer func() {
		timer.Stop()
		gameSession.cleanGameSession()
	}()
	for {
		select {}
	}
}

//helper function to verify if all players in the game session are connected
func (gameSession *GameSession) allplayersAreConnected() bool {
	for _, player := range gameSession.Players {
		if !player.IsConnected() {
			return false
		}
	}
	return true
}
