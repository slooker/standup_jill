package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nopes/slack"
	"github.com/nopes/slack/slackevents"
)

type SlackChallenge struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	EventType string `json:"type"`
}

type Event struct {
	EventType   string `json:"type"`
	User        string `json:"user"`
	Text        string `json:"text"`
	ClientMsgID string `json:"client_msg_id"`
	TS          string `json:"ts"`
	Channel     string `json:"channel"`
	EventTS     string `json:"event_ts"`
}

type SlackEvent struct {
	Token       string   `json:"token"`
	TeamID      string   `json:"team_id"`
	APIAppID    string   `json:"api_app_id"`
	Event       Event    `json:"event"`
	EventType   string   `json:"type"`
	EventID     string   `json:"event_id"`
	EventTime   int64    `json:"event_time"`
	AuthedUsers []string `json:"authed_users"`
}

/*
{
        "token": "yMR5q67W3B67whB1ABGxBjku",
        "team_id": "TC1NDN45C",
        "api_app_id": "AC1E3Q415",
        "event": {
                "type": "app_mention",
                "user": "UC093A5U3",
                "text": "<@UC0UR0L1G> Jello!",
                "client_msg_id": "8587e8c5-4ced-412b-b11d-8eb6360b4737",
                "ts": "1533139648.000201",
                "channel": "CC1NDNCG6",
                "event_ts": "1533139648.000201"
        },
        "type": "event_callback",
        "event_id": "EvC0CD9KR7",
        "event_time": 1533139648,
        "authed_users": ["UC0UR0L1G"]
}

*/

func handler(w http.ResponseWriter, r *http.Request) {
	videofile := "test"
	bucket := "test"

	if videofile == "" || bucket == "" {
		errorHandler(w, r, 400, "videofile and bucket (and optionally region) must be supplied in the query string")
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{"TOKEN"}))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		postParams := slack.PostMessageParameters{}
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			api.PostMessage(ev.Channel, "Yes, hello.", postParams)
		}
	}
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
