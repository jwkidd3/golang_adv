/*
	WAITGROUP
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var urls = []string{
	"https://google.com",
	"https://youtube.com",
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	// Declare synchronizing
	var wg sync.WaitGroup
	for _, url := range urls {
		// Sets the number of goroutines to wait for
		// Or increase the WaitGroup counter / number of goroutines by 1
		wg.Add(1)
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprintf(w, "%v\n", err)
			}
			fmt.Fprintf(w, "%v\n", resp.Status)

			// Decrements the WaitGroup counter by 1.
			// So this is called by the goroutines to indicate that it's finished
			wg.Done()
		}(url)
	}

	// Wait until the WaitGroup counter is 0
	wg.Wait()
}

func main() {
	fmt.Println("WaitGroup tutorial")
	http.HandleFunc("/health", checkHealth)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
