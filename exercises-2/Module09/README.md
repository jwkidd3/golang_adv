# Building the game service

### Overview

In this lab you will learn how to use websocket for real time socket communication, you will implement a game service that 

support game sessions with multi player.

 in order to save you some time you will be provided with the skeleton files for the various game design components, and you wil implement the buisness logic with a focus on concurrency  and synchronization 

the game service is a simple service that just act as a mediator between game clients, the game service does not contain any game logic whatsoever . the game service basically just manage sessions of connected players and transefer messages between players at the game

at this point you can take the project from the start directory , it contains everything we did so far in the labs , and provides the neccessary files as outlined here:

├── README.md

├── cmd

│  └── main.go

├── configs

│  └── app.env

├── go.mod

├── go.sum

├── internal

│  ├── games

│  │  ├── game.go

│  │  ├── gamemanager.go

│  │  ├── gamemessages.go

│  │  ├── gamesession.go

│  │  ├── gamesutils.go

│  │  ├── player.go

│  │  ├── service

│  │  │  └── game.service.go

│  │  └── utils

│  ├── routes

│  │  └── routes.go

│  └── users

│    ├── auth

│    │  ├── auth.go

│    │  └── token.go

│    ├── db

│    │  └── db.go

│    ├── service

│    │  ├── user.service.go

│    │  ├── user.service.repo.go

│    │  └── user.service_test.go

│    └── user.go

└── keys

  ├── app.rsa

  └── app.rsa.pub

## Getting started 

### Understanding the componets 

all the game related components are located under the game package at different level , lets start by talking about the new files added the packages they are at and the components they hold:

1. game.go - in the file you will find the definition of the Game type that descibes a game having id, name,description
2. Gamemanager.go -  manage game sessions. 
3. Gamesession.go - the session components to manage the players, a session is really a running game
4. Player.go -  Player components is the entity representing a player in the game
5. Game.service.go  - this files mainely define the entry points handlers for the game
6. Gamemesages.go - in this file there is the definition of mainely GameMsg and the game message identifiers 

```bash
type GameMsg struct {
	GameAction `json:"action"`       // the id of the message that is being sent ie,GAME_PLAY,START_GAME,etc...
	Data       string `json:"data"`  // the data that is being sent
	Player     Player `json:"player"` // the player(user) that sent the message
}
```

just for clarity all messages from client to server and from server to client are sent with as GameMsg .

the "Data" field in the GameMsg  contain more information relevant for the message that is being sent.

there are two more types defined in this file **StartGameMsg{}** and **OnNewGameSessionCreated{}** (used for the handshake process)

they serve mainely for convinience of serializing and deserializing  from the **"Data"** of the GameMsg

### The Game protocol

the game protocol is very straightforward , 

1. the Server is up and  a GameManager object is created

2. The client authenticate with the game server using the register or login of the user service

3. the client send a get requests to get the list of all the users (so he will be able to choose who it wants to play with)

4. the client send get request to **/games/gameinfo** endpoint with the gameid to verify that the server supports it, if so

5. a client send a websocket request for a new game through the **/games/startgame** passing in the query the **gameid** and **access-token** endpoint 

6. the startgame handler validate the gameid  (user is validated in the jwt middleware) and upgrade the connection to websocket , and then on the websocket sends back to the client a message with action id ON_GAME_SESSION_CREATED and data containing OnNewGameSessionCreated struct with the **sessionid** (string) and the **gameinfo**(Game)

7. after the client get the ON_GAME_SESSION_CREATED indicating a session was created, it sends t othe server a START_GAME message , using the StartGameMsg (wrapped in GameMsg off course) , this message contains the list of players for the game and an initial data to send for each player that is connected to this game session

8. when an invited client user is connecting to the session , he will do it through another endpoint 

    **/games/joingame/{game token}** and will have to pass in the query string the **?x-access-token=exampletoken** the handler will  send the user ON_GAME_INIT message with the inital data for the game it got from the owner who created the game ,connect the user to the session and send all the other users in the session the ON_USER_CONNECTED message with the connected user in the data field

9. when a player makes a move it sends the GAME_PLAY messge (the server is sgnostic to what is in the data field) and potentially UPDATE_GAME_STATE to replace what's was given in the initial when the game was created 

10. ON_USER_DISCONNECTED is sent by the session to everyone when user is disconnected from the game

11. ON_GAME_OVER will be sent if not all were online for the last 30 minuates (it will clear the session)

