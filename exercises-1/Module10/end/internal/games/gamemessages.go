package games

type GameAction string

const (
	START_GAME              GameAction = "START_GAME"
	GAME_PLAY                          = "GAME_PLAY"
	UPDATE_GAME_STATE                  = "UPDATE_GAME_STATE"
	ON_GAME_SESSION_CREATED            = "ON_GAME_SESSION_CREATED"
	ON_GAME_OVER                       = "ON_GAME_OVER"
	ON_GAME_INIT                       = "ON_GAME_INIT"
	ON_USER_CONNECTED                  = "ON_USER_CONNECTED"
	ON_USER_DISCONNECTED               = "ON_USER_DISCONNECTED"
)

type GameMsg struct {
	GameAction `json:"action"`
	Data       string `json:"data"`
	Player     Player `json:"player"`
}

//the create game happens through http

type StartGameMsg struct {
	Players  []Player `json:"players"`
	GameData string   `json:"gamedata"`
}

type OnNewGameSessionCreated struct {
	Game      `json:"game"`
	SessionID string `json:"id"`
}
