package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BirthdayResponse struct {
	Status       string `json:"status"`
	Letter       string `json:"letter"`
	What_you_are string `json:"what_you_are"`
	Has_more     string `json:"has_more"`
}

type Acrostic struct {
	Letter  string
	Message string
	Last    string
}

var ACROSTICS []*Acrostic = []*Acrostic{
	&Acrostic{"J", "Just. As in 'Justice'. NOT as in, like, 'eh'.", "tRuE"},
	&Acrostic{"i", "Intent on what you want.", "Oh, Lord, if I can spell right then yeah. Probably."},
	&Acrostic{"l", "Lady Reagent of the Princessipality of Math.", "There' always more for the Math Princess."},
	&Acrostic{"l", "...really good at parsing Logs? Damned second L...", "Sorry for that one, there's more."},
	&Acrostic{"i", "Intergalactic hyper star.", "No points for guessing the next letter."},
	&Acrostic{"a", "Always there.", "There is also always more."},
	&Acrostic{"n", "Much Nice.", "Is he going to try for the last name, too!?"},

	&Acrostic{"K", "By far the Kindest and most considerate person I know.", "Of fucking course."},
	&Acrostic{"a", "So amazingly Accomplished.", "Ooof this is getting hard. I'm not good at these."},
	&Acrostic{"r", "400lbs by 25. Yeah, I gave up on this one. Too many letteRs.", "Aaaaalmost there!"},
	&Acrostic{"n", "N is relatively close to S so I'm just going to point out how smart you are.", "He's really running out of gas, now..."},
	&Acrostic{"e", "So Exuberant and full of life (when Phoenix isn't acively trying to suck it out of you).", "THIS IS HARDER THAN I THOUGHT IT WAS GOING TO BE."},
	&Acrostic{"r", "Have a Really awesome birthday, Jillian. You deserve it.", "Well, no. Not really. But it does loop around anyways."},
}

func getResponse(resetting bool) BirthdayResponse {
	next := getNext(resetting)
	return BirthdayResponse{
		http.StatusText(http.StatusOK),
		next.Letter,
		next.Message,
		next.Last,
	}
}

func _next(in chan bool, out chan *Acrostic) {
	current := 0
	for resetting := range in {
		if resetting {
			current = 0
		}
		out <- ACROSTICS[current]
		current = (current + 1) % len(ACROSTICS)
	}
}

var getNext func(resetting bool) *Acrostic = func() func(bool) *Acrostic {
	in := make(chan bool)
	out := make(chan *Acrostic)
	go _next(in, out)
	return func(resetting bool) *Acrostic {
		in <- resetting
		return <-out
	}
}()

func goodMorning(w http.ResponseWriter, r *http.Request) {
	returnJson(w, r, struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{http.StatusText(http.StatusOK), "Happy birthday, Jillian!"})
}

func next(w http.ResponseWriter, r *http.Request) {
	returnJson(w, r, getResponse(false))
}

func restart(w http.ResponseWriter, r *http.Request) {
	returnJson(w, r, getResponse(true))
}

func returnJson(w http.ResponseWriter, r *http.Request, obj interface{}) {
	value, err := json.MarshalIndent(obj, "", "   ")
	if err != nil {
		io.WriteString(w, "Oh man, I fucked up.\n")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(append(value, '\n'))
}

func main() {
	http.HandleFunc("/goodmorning", goodMorning)
	http.HandleFunc("/poem", next)
	http.HandleFunc("/restart", restart)
	fmt.Println(http.ListenAndServe(":8080", nil))
}
