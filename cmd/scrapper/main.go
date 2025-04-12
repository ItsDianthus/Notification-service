// scrapper/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Command struct {
	Command string `json:"command"`
	URL     string `json:"url"`
}

func handleBotCommand(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request...")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var cmd Command
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cmd); err != nil {
		http.Error(w, "Failed to decode command", http.StatusBadRequest)
		return
	}

	if cmd.Command == "get" {

		fmt.Println("Command is get. ")
		fmt.Println(cmd.Command, cmd.URL)
		if cmd.URL != "" {
			fmt.Fprintf(w, "Fetching data from %s...\n", cmd.URL)
		} else {
			http.Error(w, "URL is required for the 'get' command", http.StatusBadRequest)
		}
	} else {
		http.Error(w, "Unknown command", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/command", handleBotCommand)

	port := ":8080"
	fmt.Printf("Scrapper service is running on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
