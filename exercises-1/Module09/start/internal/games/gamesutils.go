package games

import (
	"encoding/json"
	"errors"
	"log"
)

//helper function for creating a GameMsg object
//this function gets the GameAction which entails what message it is,
// the cmd data, and the player who send the message, and creates a GameMsg
func WrapCommand(action GameAction, cmd interface{}, player Player) (GameMsg, error) {

	var dataStr string
	if t, ok := cmd.(string); !ok {
		msg, err := json.Marshal(cmd)
		if err != nil {
			log.Println("couldn't join player")
			//TBD exception should be thrown here to the client by sending a json error message
			return GameMsg{}, errors.New("couldn't Marshal Object")
		}
		dataStr = string(msg)
	} else {
		dataStr = string(t)
	}

	var gameMsg = GameMsg{
		GameAction: action,
		Data:       dataStr,
		Player:     player,
	}
	return gameMsg, nil
}

//helper function for stripping the data from GameMsg object, currently just being used for convinience
//for stripping the data of a START_GAME action
func UnWrapGameMsg(gameMsg GameMsg) interface{} {
	switch gameMsg.GameAction {
	case START_GAME:
		var startGame StartGameMsg
		json.Unmarshal([]byte(gameMsg.Data), &startGame)
		return startGame
	}
	return nil
}
