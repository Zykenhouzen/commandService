package main

import (
	"encoding/json"
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
	slackCommandInputText := r.FormValue("text")

	textSpaceDelineated := strings.Split(slackCommandInputText, " ")
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

	returnJSON := ReturnCommand{
		Text:         groups,
		ResponseType: "in_channel",
	}
	//b, err := json.Marshal(returnJSON)

	//	if err != nil {
	//		panic(err)
	//}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(returnJSON)
}

type ReturnCommand struct {
	Text         string `json:"text"`
	ResponseType string `json:"response_type"`
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
