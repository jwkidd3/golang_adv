package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/someuser/gameserver/internal/games"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var gameManager games.GameManager

//ActiveGames type
func init() {

	gameManager = games.CreateGameManager()
	go gameManager.Run()

}

//open a web socket for the request and sends it back
func openWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

//HandleUserJoinedGame type
func HandleUserJoinedGame(w http.ResponseWriter, r *http.Request, SessionID string) error {

	//1. get the user from the request context

	//2. call the openWebSocket to get the connection

	//3 . get the game session

	//4. create the player

	//5. start the players

	//6. call the SendCurrentGameStateToPlayer to send the inital game session data to the newly joined player

}

func validatGame(w http.ResponseWriter, r *http.Request) (games.Game, error) {

	if keys, ok := r.URL.Query()["gameid"]; ok {
		id := keys[0]
		return gameManager.GetGame(id)
	}

	return games.Game{}, errors.New("couldnt find gameid in query")
}

//HandleStartGame type
func HandleStartGame(w http.ResponseWriter, r *http.Request) error {

	//try to get the user from the request context
	//...

	//1. get the user from the request context

	//2. make sure the gameid is valid using the helper validatGame

	//3. call the openWebSocket to get the connection

	//3 . create a new game session

	//4. start the game session run loop

	//5. create the player

	//5. start the player

	//6. call the SendCurrentGameStateToPlayer

	// 7. send back the OnNewGameSessionCreated

	// 8.send back a message to the player updating him that the game sesion is created
	// and that he can send invitation to players games.ON_GAME_SESSION_CREATED
	//you can use the games.WrapCommand for easier construction of the message

	return nil

}

//JoinGame called for joining a user to a game session
func JoinGame(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var id = params["gametoken"]
	if id == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := HandleUserJoinedGame(w, r, id); err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

}

//StartNewGame called for creating a new game session
func StartNewGame(w http.ResponseWriter, r *http.Request) {
	if err := HandleStartGame(w, r); err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
}

//StartNewGame called for creating a new game session
func GetGameInfo(w http.ResponseWriter, r *http.Request) {

	g, err := validatGame(w, r)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var resp = map[string]interface{}{"status": true, "message": g}
	json.NewEncoder(w).Encode(resp)
	return

}
