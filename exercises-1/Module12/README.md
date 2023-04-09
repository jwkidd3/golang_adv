# Deploying the Game Service as a kubernetes operator

### Overview

in this lab you will create a kubernetes operator to manage you game service , the operator can support any number of games provided they have different gameid.

each gameserver is deployed with a configmap object that holds the gameid, description and name for it.

the configmap is created automatically by the operator and does not need to be handled manually by the administrator

also one instance of mysql is deployed for all game servers , and it will be created on first game server CR deployment

## Getting started 

### installing the operator sdk

```bash
Download from https://github.com/operator-framework/operator-sdk/releases 
```

install:

```bash
$ wget https://github.com/operator-framework/operator-sdk/releases/download/v0.../...
$ sudo mv operator-sdk-v... /usr/local/bin/operator-sdk
$ sudo chmod +x /usr/local/bin/operator-sdk
```

### bootstrapping the Project 

1. create a new project (add your own repo location)

```bash
operator-sdk new gameserver --repo github.com/<yourname>/gameserver
```

### Creating the api object(CRD) 

1. run the following docker container that implement the game server protocol , create 2 players and start playing

```bash
operator-sdk add api --api-version=gameserver.com/v1 --kind=Gameserver
```

2. define the spec for the gameserver CRD, under pkg/apis/gameserver/v1/gameserver_types.go, add the following spec

```go
	type GameserverSpec struct {
    GameID      string `json:"gameid"`
    Name        string `json:"name"`
    Description string `json:"description"`
    ServerPort  int32  `json:"port"`
}
```

3. ask operator-sdk to regenerate the the crds

```bash
$ operator-sdk generate crds
```

### Adding the controller 

### Creating the api object(CRD) 

1. add the conrtoller to control the Gameserver Api object

```bash
operator-sdk add controller --api-version=gameserver.com/v1 --kind=Gameserver
```

2. ask operator-sdk to regenerate the the ck8s objects

```bash
$ operator-sdk generate k8s
```

### Add the reconciliation and watch logic  

1. in the start directory , you are provided with 3 files to get you up to speed common.go,mysql.co , and backend.go
2. use the functions in those files for creating the required objects in the reconcile function at the pkg/controller/gameserver/gameserver_controller.go file

```go
func (r *ReconcileGameserver) Reconcile(request reconcile.Request) (reconcile.Result, error) {...}
```

3. also add the relevant object to "watche" in the add function 

```go
func add(mgr manager.Manager, r reconcile.Reconciler) error {..}
```

### Deploy to the kind cluster  

1. from your root folder - deploy the CRD'd 

```bash
kubectl apply -f deploy/crds/*_crd.yaml 
```

2. under deploy/crds/gameserver.com_v1_gameserver_cr.yaml change the default cr to:

```bash
apiVersion: gameserver.com/v1
kind: Gameserver
metadata:
  name: gameserver
spec:
  # Add fields here
  gameid: "pokemoncards"
  name: "Pokemon memory game"
  description: "a memory game in whch you have to match two exact cards"
  port: 30685
```

3. from your root folder - deploy the CR 

```bash
kubectl apply -f deploy/crds/*_cr.yaml
```

### Test your operator 

1. run the following docker container that implement the game server protocol , create 2 players and start playing

   ```bash
   docker run --env GAMES_SERVER_URL=localhost:30685 -p 3002:3002 motisoffer/pokemongame:1.0
   ```

2. from your root directory on your project and in a different terminal run 

   ```bash
   operator-sdk run local --namespace default
   ```

### Debug your operator 

1. you will need to run your operator with delve support 

   ```bash
   operator-sdk run local --namespace default --enable-delve
   ```

2. You will need a launch json for Vscode to interact with this headless mode of delve

```bash
{
    "version": "0.2.0",
    "configurations": [
      {
        "name": "gameserver operator",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "program": "${workspaceFolder}/cmd/manager/main.go",
        "env": {
          "WATCH_NAMESPACE": "default"
        },
        "args": []
      }
    ]
  }
```

3. add break point and run vs code

### Deploy  your operator 

1. build the image

   ```bash
   operator-sdk build motisoffer/gogameserver:1.0
   ```

2. configure the deployment that the operator-sdk generated for you deploy/operator.yaml with the right container image

```bash
containers:
    - name: gameserver
          # Replace this with the built image name
          image: REPLACE_IMAGE
          command:
```

3. Deploy the crds to the cluster

```bash
$ kubectl apply -f deploy/crds/*_crd.yaml
```

4. Deploy the service account and role* that was auto generated as well under the deploy folder (note you should change based on your required permissions)

```bash
$ kubectl apply -f deploy/service_account.yaml 
$ kubectl apply -f deploy/role.yaml
$ kubectl apply -f deploy/role_binding.yaml
```

5. deploy the operator to the cluster

```bash
$ kubectl apply -f deploy/operator.yaml
```

