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
	for _, player := range gameSession.Players {
		//only if we have a conn ready
		if !player.IsConnected() {
			continue
		}
		//skip sending to the one who send the message , unless it is ment for all
		if gameMsg.Player == (Player{}) { //if ment for all
			player.SendMessage(gameMsg)
		} else if gameMsg.Player.Email != player.Email {
			player.SendMessage(gameMsg)
		}
	}
}

//add the intvited users to session and wait for them to join the game
// when a user joins the game he become a player
func (gameSession *GameSession) addUsersToSession(players []Player) {
	for _, player := range players {
		gameSession.Players[player.Email] = &player
	}
}
func (gameSession *GameSession) setInitData(data string) {
	gameSession.InitialGameData = data
}

func (gameSession *GameSession) cleanGameSession() {
	for _, player := range gameSession.Players {
		gameSession.removeUser(player)
	}
	gameSession.gameManager.unRegister <- gameSession
}

func (gameSession *GameSession) removeUser(player *Player) {
	player.Stop()
	delete(gameSession.Players, player.Email)
}
func (gameSession *GameSession) Run() {

	timer := time.NewTimer(maxGameStartTime)
	defer func() {
		timer.Stop()
		gameSession.cleanGameSession()
	}()
	for {
		select {
		case player := <-gameSession.Register:
			//if the user exists in the invite list but not yet active override it with the new one
			//if the user is trying to reconnect drop the old connection in favour of the new one
			if val, ok := gameSession.Players[player.Email]; ok {
				gameSession.removeUser(val)
			}
			gameSession.Players[player.Email] = player
			gameData, _ := WrapCommand(ON_USER_CONNECTED, *player, *player)
			gameSession.sendMsgToPlayers(&gameData)

		case player := <-gameSession.UnRegister:
			if val, ok := gameSession.Players[player.Email]; ok {
				gameData, _ := WrapCommand(ON_USER_DISCONNECTED, *player, *player)
				gameSession.sendMsgToPlayers(&gameData)
				gameSession.removeUser(val)
			}

		case gameMsg := <-gameSession.SendToGame:
			msg := UnWrapGameMsg(*gameMsg)
			if t, ok := msg.(StartGameMsg); ok == true {
				gameSession.addUsersToSession(t.Players)
				gameSession.setInitData(t.GameData)
			} else if gameMsg.GameAction == UPDATE_GAME_STATE {
				gameSession.setInitData(gameMsg.Data)
			} else {
				timer.Reset(maxGameStartTime)
				gameSession.sendMsgToPlayers(gameMsg)
			}

		case <-timer.C:
			//check if there is no one on the session then delete the session
			if !gameSession.allplayersAreConnected() {
				exception := struct {
					Message string
				}{Message: "time out : not all participants have joined"}
				msg, _ := WrapCommand(ON_GAME_OVER, exception, Player{})
				gameSession.sendMsgToPlayers(&msg)
				return
			} else if len(gameSession.Players) == 0 {
				return
			} else {
				timer.Reset(maxGameStartTime)
			}

		}
	}
}

func (gameSession *GameSession) allplayersAreConnected() bool {
	for _, player := range gameSession.Players {
		if !player.IsConnected() {
			return false
		}
	}
	return true
}
