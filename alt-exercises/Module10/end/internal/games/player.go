package games

import (
	"log"

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
	player.Conn = conn
	player.RecvMsgChan = make(chan GameMsg)

	//register to session
	player.GameSession.Register <- player

	//open for recieving and sending
	go player.handleMessageToPlayer()
	go player.recieveMessages()

}

func (player *Player) Stop() {

	if player.RecvMsgChan != nil {
		close(player.RecvMsgChan)
	}
}

//RecieveMessages from the players
func (player *Player) recieveMessages() {
	defer func() {
		player.GameSession.UnRegister <- player
	}()

	for {
		var gameMsg GameMsg
		if err := player.Conn.ReadJSON(&gameMsg); err != nil {
			log.Println(err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v %s %s", err, player.Email, player.GameSession.ID)
			}
			break
		}
		//add the current user who sends the message
		gameMsg.Player = *player
		player.GameSession.SendToGame <- &gameMsg
	}
}
func (player *Player) SendMessage(msg *GameMsg) {
	player.RecvMsgChan <- *msg
}

//HandleMessageToPlayer sends the message from the game to the player through a channel
func (player *Player) handleMessageToPlayer() {

	for {
		select {
		case msg, ok := <-player.RecvMsgChan:
			if ok {
				player.Conn.WriteJSON(msg)
			} else {
				if player.Conn != nil {
					player.Conn.Close()
					player.Conn = nil
				}
				return
			}
		}
	}
}
func (player *Player) SendCurrentGameStateToPlayer() {

	//send the initial data for the user
	gameData, err := WrapCommand(ON_GAME_INIT, player.GameSession.InitialGameData, *player)

	if err != nil {
		return
	}
	player.SendMessage(&gameData)
}
