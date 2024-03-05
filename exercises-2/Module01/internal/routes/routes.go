package routes

import "net/http"

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to go!!!"))
}

func greet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from go!!!"))
}
func Handlers() {
	http.HandleFunc("/welcome", welcome)
	http.HandleFunc("/greet", greet)
//http.HandleFunc("/greet", http.HandlerFunc(greet))
}
