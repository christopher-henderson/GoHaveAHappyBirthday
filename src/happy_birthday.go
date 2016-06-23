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
	&Acrostic{"i", "Intergalactic hyper star.", "Probably"},
	&Acrostic{"l", "Last time I heard, really damned nice.", "Of fucking course."},
	&Acrostic{"l", "...really good at parsing Logs? Damned second L...", "Sorry for that one, there's more."},
	&Acrostic{"i", "Intergalactic hyper star.", "Of fucking course."},
	&Acrostic{"a", "Always there.", "Of fucking course."},
	&Acrostic{"n", "Last time I heard, really damned nice.", "Of fucking course."},
	
	&Acrostic{"K", "By far the Kindest and most considerate person I know.", "Of fucking course."},
	&Acrostic{"a", "So amazingly Accomplished.", "Of fucking course."},
	&Acrostic{"r", "Last time I heard, really damned nice.", "Of fucking course."},
	&Acrostic{"n", "Last time I heard, really damned nice.", "Of fucking course."},
	&Acrostic{"e", "Last time I heard, really damned nice.", "Of fucking course."},
	&Acrostic{"r", "Last time I heard, really damned nice.", "Of fucking course."},
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
	http.HandleFunc("/next", next)
	http.HandleFunc("/restart", restart)
	fmt.Println(http.ListenAndServe(":8080", nil))
}
