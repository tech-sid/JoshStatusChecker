package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var server string = "localhost:3000"

func main() {
	fmt.Println("Starting server at", server)
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(server, nil); err != nil {
		log.Printf("Error encountered while starting server: %v \n", err)
	}
}

var input = make(map[string]string)

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		website := r.URL.Query().Get("website")
		if website != "" {
			response, _ := input[website]
			fmt.Fprintln(w, website, response)

		} else {
			for key, val := range input {
				fmt.Fprintln(w, key, " ", val)
			}
		}
	} else if r.Method == "POST" {
		var user_input []string
		err := json.NewDecoder(r.Body).Decode(&user_input)
		if err != nil {
			fmt.Fprintf(w, "Error encountered while running post request: %v", err)
		}

		var channel chan string
		for _, domain := range user_input {
			go website_stat(domain, channel)
		}
		for itr := range channel {
			go func(query string) {
				time.Sleep(60 * time.Second)
				website_stat(query, channel)
			}(itr)
		}
	} else {
		fmt.Fprint(w, "Error encountered while processing request")
	}
}

func website_stat(query string, channel chan string) {
	_, err := http.Get(query)
	if err != nil {
		fmt.Println("404", query)
		input[query] = "DOWN"
		channel <- query
	} else {
		fmt.Println("200 OK", query)
		input[query] = "UP"
		channel <- query
	}
}
