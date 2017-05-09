package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func SplitGroup(w http.ResponseWriter, r *http.Request) {
	var slackCommandInput SlackCommandInput
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &slackCommandInput); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	textSpaceDelineated := strings.Split(slackCommandInput.Text, " ")
	groupSizes := strings.Split(textSpaceDelineated[0], ":")
	names := append(textSpaceDelineated[:0], textSpaceDelineated[1:]...)
	namesLeft := names
	groups := ""

	rand.Seed(time.Now().UnixNano())
	currentGroup := 0
	namesInCurrentGroup := 0
	groupCount := 0

	for nameCount := 0; nameCount < len(names); nameCount++ {
		if namesInCurrentGroup == 0 {
			groupCount++
			groups = groups + "Group " + strconv.Itoa(groupCount) + ": "
		}
		indexChosen := rand.Intn(len(namesLeft))
		groups = groups + " " + namesLeft[indexChosen]
		namesLeft = append(namesLeft[:indexChosen], namesLeft[1+indexChosen:]...)

		namesInCurrentGroup++
		if strconv.Itoa(namesInCurrentGroup) >= groupSizes[currentGroup] {
			namesInCurrentGroup = 0
			currentGroup++
			if currentGroup > (len(groupSizes) - 1) {
				currentGroup = 0
			}
			groups = groups + "\n"
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		panic(err)
	}
}

type SlackCommandInput struct {
	Token        string `json:"token"`
	Team_id      string `json:"team_id"`
	Team_domain  string `json:"team_domain"`
	Channel_id   string `json:"channel_id"`
	Channel_name string `json:"channel_name"`
	User_id      string `json:"user_id"`
	User_name    string `json:"user_name"`
	Command      string `json:"command"`
	Text         string `json:"text"`
	Response_url string `json:"response_url"`
}
