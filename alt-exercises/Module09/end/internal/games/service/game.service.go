package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/someuser/gameserver/internal/games"

	"github.com/someuser/gameserver/internal/users"
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

	//first lets validate that the user is authenticated
	var user *users.User
	if m := r.Context().Value("user"); m != nil {
		if val, ok := m.(*users.User); ok {
			user = val
		} else {
			return errors.New("not a valid user")
		}
	}

	conn, err := openWebSocket(w, r)
	if err != nil {
		return err
	}

	gameSession := gameManager.GetSessionByID(SessionID)

	if gameSession == nil {
		return errors.New("no such game session exists")
	}

	player := gameSession.CreateNewPlayer(conn, user.ID, user.Name, user.Email)

	if player == nil {
		return errors.New("Invalid user for game")
	}

	player.Start(conn)

	player.SendCurrentGameStateToPlayer()

	return nil

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
	//first lets validate that the user is authenticated
	var user *users.User
	if m := r.Context().Value("user"); m != nil {
		if val, ok := m.(*users.User); ok {
			user = val
		} else {
			return errors.New("not a valid user")
		}
	}

	g, err := validatGame(w, r)
	if err != nil {
		return err
	}
	conn, err := openWebSocket(w, r)
	if err != nil {
		return err
	}

	gameSession := gameManager.CreateNewGameSession()
	go gameSession.Run()

	player := gameSession.CreateNewPlayer(conn, user.ID, user.Name, user.Email)
	if player == nil {
		return errors.New("Invalid user for game")
	}

	player.Start(conn)

	//for simplicity we always returns the same game type,
	// and do not use the id potentially passed to see if the game is supported
	var msgPlay = games.OnNewGameSessionCreated{
		SessionID: gameSession.ID,
		Game:      g,
	}

	//send back a message to the host updating him that the game sesion is created
	// and that he can send invitation to players
	gameMsg, err := games.WrapCommand(games.ON_GAME_SESSION_CREATED, &msgPlay, *player)
	if err != nil {
		return err
	}
	player.SendMessage(&gameMsg)

	return nil
	//SendGameInvitation(gameSession.ID, user.Email, gameSession.Players)

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
