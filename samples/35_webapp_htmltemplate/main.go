package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	tmplView = "template/view.html"
	dataFile = "data/todos.txt"
)

type TodoList struct {
	Amount int
	Todos  []string
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writeResponse(w http.ResponseWriter, msg string) {
	// Converse msg to bytes
	_, err := w.Write([]byte(msg))
	handleError(err)
}

func greetingHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "Hello there!")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	// Get input data from form
	todo := r.FormValue("todo")

	// Options for working with file
	options := os.O_CREATE | os.O_WRONLY | os.O_APPEND

	file, err := os.OpenFile(dataFile, options, os.FileMode(0600))
	handleError(err)

	// Append new record to file
	_, err = fmt.Fprintln(file, todo)
	handleError(err)

	err = file.Close()
	handleError(err)

	http.Redirect(w, r, "/todos", http.StatusFound)
}

// Get lines of text from file
func getStrings(filePath string) (lines []string) {
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		_, err := os.Create(filePath)
		handleError(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	handleError(scanner.Err())

	return
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	// Get text from file
	vals := getStrings(dataFile)
	fmt.Printf("%#v\n", vals)

	tmpl, err := template.ParseFiles(tmplView)
	handleError(err)

	todos := TodoList{
		Amount: len(vals),
		Todos:  vals,
	}

	// Write the template to ResponseWriter
	// Pass data to todo parameter in template
	err = tmpl.Execute(w, todos)
	handleError(err)
}

func setupRoutes() {
	// Receive the requests from urls
	http.HandleFunc("/", greetingHandler)
	http.HandleFunc("/todos", listHandler)
	http.HandleFunc("/todos/create", createHandler)
}

func main() {
	setupRoutes()
	fmt.Println("Start server at 8080")

	// Run server at localhost:8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
