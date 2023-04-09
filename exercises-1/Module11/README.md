# Preparing our app for deployment

### Overview

In this lab you will prepare the gameserver app for deployment, by containeraizing your app, and creating the necessary deployment.yaml files for being able to deploy to Kubernetes cluster.

├── Dockerfile

├── cmd

│  └── main.go

├── configs

│  └── app.env

├── deployments

│  ├── config.yaml

│  ├── database.yaml

│  └── gameserver.yaml

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

│  │  └── service

│  │    └── game.service.go

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

### creating a docker file 

1. create a docker file for the game server , support an environment variable names GAME_SERVER_HOMEDIR

   that will store the path to the configuration files, and secured keys.

2. change your code to use the GAME_SERVER_HOMEDIR if it exists to load configuration from that path(db,game info, keys)

### create the deployment for the gameserver app

1. create a service , deployment and a configmap for the gameserver app.

   the configmap should store they Game information as environment vars that the app will be able to load from

   (GAME_DESCRIPTION,GAME_ID,GAME_NAME)

   **Please note that the mysql.yaml deployment and service is provided for you to save time**

2. Run kind cluster with the following configuration "kind create cluster --config kind.yaml"

```bash
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 31654
    hostPort: 31654
  - containerPort: 30685
    hostPort: 30685
- role: worker
- role: worker
```

3. deploy the gameserver.yaml , configmap.yaml and mysql.yaml to the cluster

```bash
	kubectl apply -f yourdir/deployments
```

### Testing 

1. run the following docker container that implement the game server protocol , create 2 players and start playing

   ```bash
   docker run --env GAMES_SERVER_URL=localhost:30685 -p 3002:3002 motisoffer/pokemongame:1.0
   ```


