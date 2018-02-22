package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/nlopes/slack"
)

type card struct {
	Name      string            `json:"name"`
	ImageURIs map[string]string `json:"image_uris"`
}

func main() {

	token := "<SLACK-TOKEN>"
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				log.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				log.Printf("Message: %v\n", ev)

				// TODO: allow this to hit multiple cards
				re, err := regexp.Compile(`\[\[(.+)\]\]`)
				if err != nil {
					log.Fatal(err)
				}
				// Change this once better regex is made
				result := re.FindAllStringSubmatch(ev.Text, 1)
				if result != nil {
					c := result[0][1]
					url := "https://api.scryfall.com/cards/named?fuzzy=" + url.PathEscape(c)
					resp, err := http.Get(url)
					if err != nil {
						log.Println(err)
					}
					defer resp.Body.Close()

					var cd card
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Println(err)
					}
					err = json.Unmarshal(body, &cd)
					if err != nil {
						log.Println(err)
					}
					rtm.SendMessage(rtm.NewOutgoingMessage(string(cd.ImageURIs["normal"]), ev.Channel))

				}

			case *slack.RTMError:
				log.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop
			}
		}
	}
}
