package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ItemsJson struct {
	TwoGb []Item `json:"2gb"`
}

type Item struct {
	Store string `json:"store"`
	Link  string `json:"link"`
	Ram   int    `json:"ram"`
	Misc  string `json:"misc"`
}

func readJsonFile() ItemsJson {
	content, err := os.ReadFile("./store.json")
	if err != nil {
		log.Fatalln(err)
	}

	var data ItemsJson

	if err := json.Unmarshal(content, &data); err != nil {
		log.Fatalln(err)
	}

	return data
}

func getStoreHtmlBody(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	html := string(body)
	return html
}

func isVilrosAvailable(i Item) (bool, error) {
	resp, err := http.Get(i.Misc)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var jsonMapping map[string]interface{}
	json.Unmarshal([]byte(body), &jsonMapping)

	isAvailable, ok := jsonMapping["available"].(bool)
	if !ok {
		return false, errors.New("Json mapping for Vilros did not return boolean value")
	}
	return isAvailable, nil
}

func pollStores(stores ItemsJson) {
	for {
		for _, item := range stores.TwoGb {
            switch item.Store {
                case "adafruit":
                    fmt.Println(item.Store)
                case "pishop.us":
                    fmt.Println(item.Store)
                case "vilros":
                    // isVilrosAvailable(item)
                    fmt.Println(item.Store)
            }
		}

		time.Sleep(5 * time.Second)
	}
}

func initBot() {
    token := "NjYzMTcwNTIyMDc4NTExMTA0.GEfq7G._azpPjUB_fajKlZi6VDgK7r7_pvRF3mrwdNj88"

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(pingMe)

	// Receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func main() {
    data := readJsonFile()
	go pollStores(data)
    initBot()
}

func pingMe(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
    
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
