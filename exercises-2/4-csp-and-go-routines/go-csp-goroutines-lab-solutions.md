

# Go


## CSP and goroutines

```go
package main

import "fmt"
import "math/rand"
import "time"

func hits(probability int) bool {
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    rand := r1.Intn(100)
    fmt.Println(rand)
    if rand > probability {
        return false
    }
    return true
}

func play(player1prob, player2prob int) {
    table := make(chan string)
    playing := true
    i := 0
    fmt.Println("Press enter to begin")
    var input string
    var player int
    fmt.Scanln(&input)
    for playing {
        player = 0
        fmt.Println("Rally Count:", i, "\n")
        playing = hits(player1prob)
        if !playing {
            fmt.Println("Oops! Player", player, "missed")
            break
        }
        go func() { table <- "ping" }()
        player1 := <- table
        fmt.Println("Player0:", player1)
        fmt.Scanln(&input)
        player = 1
        playing = hits(player2prob)
        if !playing {
            fmt.Println("Oops! Player", player, "missed")
            break
        }
        go func() { table <- "pong"}()
        player2 := <- table
        fmt.Println("Player1:",player2)
        fmt.Scanln(&input)
        i ++
    }
    player = (player + 1) % 2
    fmt.Println("Player", player, "won")
}
func main() {
    play(99, 76)
}
```
