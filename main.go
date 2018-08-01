package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SlackChallenge struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	EventType string `json:"type"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	videofile := "test"
	bucket := "test"

	if videofile == "" || bucket == "" {
		errorHandler(w, r, 400, "videofile and bucket (and optionally region) must be supplied in the query string")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, r, 400, fmt.Sprintf("Could not read body: %s", err.Error))
		return
	}
	fmt.Printf("Body: %s\n", body)

	challenge := new(SlackChallenge)
	err = json.Unmarshal(body, &challenge)
	if err != nil {
		errorHandler(w, r, 400, fmt.Sprintf("Could not read body: %s", err.Error))
		return
	}
	fmt.Printf("Challenge: %v\n", challenge)

	err = setCredentials()
	if err != nil {
		http.Error(w, "Unable to set credentials", 401)
		log.Fatal(err)
	}
	//fmt.Fprintf(w, "Hit the endpoints we wanted to hit.")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	fmt.Fprintf(w, challenge.Challenge)

}

// initiates the webapp and downloads dependencies
func main() {
	fmt.Println("v0.40")

	err := setCredentials()
	if err != nil {
		log.Fatal("Error setting credentials: ", err)
	}

	http.HandleFunc("/", handler)
	fmt.Println("Starting server on port 1313")
	log.Fatal(http.ListenAndServe(":1313", nil))
}

func setCredentials() (err error) {
	return nil
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprint(w, message)
}

