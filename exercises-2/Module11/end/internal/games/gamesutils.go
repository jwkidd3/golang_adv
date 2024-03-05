package games

import (
	"encoding/json"
	"errors"
	"log"
)

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

func UnWrapGameMsg(gameMsg GameMsg) interface{} {
	switch gameMsg.GameAction {
	case START_GAME:
		var startGame StartGameMsg
		json.Unmarshal([]byte(gameMsg.Data), &startGame)
		return startGame
	}
	return nil
}
