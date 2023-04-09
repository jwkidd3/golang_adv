package games

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Conn        *websocket.Conn `json:"-"`
	RecvMsgChan chan GameMsg    `json:"-"`
	GameSession *GameSession    `json:"-"`
}

func (player *Player) IsConnected() bool {
	if player.Conn != nil {
		return true
	}
	return false
}
func (player *Player) Start(conn *websocket.Conn) {

	//init  connection and channel

	//register to session

	//open go routines for recieving and sending

}

func (player *Player) Stop() {

	if player.RecvMsgChan != nil {
		close(player.RecvMsgChan) //close the channe and signal the listening select loop finish
	}
}

//RecieveMessages from the players
func (player *Player) recieveMessages() {
	defer func() {
		player.GameSession.UnRegister <- player
	}()

	for {
		//recieve message on the websocket
		//break if an error
		//else send the message to the game session
	}
}

//send message to the palyer using the players player.RecvMsgChan
func (player *Player) SendMessage(msg *GameMsg) {
}

//HandleMessageToPlayer sends the message from the game to the player using the players player.RecvMsgChan
func (player *Player) handleMessageToPlayer() {

	for {
		select {}
	}
}

//this method is called right after a player joined an existing game session,
// and send the intial game data to him
func (player *Player) SendCurrentGameStateToPlayer() {

	//send the initial data for the user
	gameData, err := WrapCommand(ON_GAME_INIT, player.GameSession.InitialGameData, *player)

	if err != nil {
		return
	}
	player.SendMessage(&gameData)
}
