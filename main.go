package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type PiItemsJson struct {
	TwoGb []Item `json:"2gb"`
}

type Item struct {
	Store string `json:"store"`
	Link  string `json:"link"`
	Ram   int    `json:"ram"`
	Misc  string `json:"misc"`
}

func readJsonFile() PiItemsJson {
	content, err := os.ReadFile("./store.json")
	if err != nil {
		log.Fatalln(err)
	}

	var data PiItemsJson

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

func isVilrosAvailable(i Item) bool {
	resp, err := http.Get(i.Misc)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	if result["available"] == true {
		return true
	} else {
		return false
	}
}

func testLoop(stores PiItemsJson) {
	// chicacgo electonic distributors
	// adafruit
	// pishop.us
	// sparkfun
	// vilros
	// canakit
	// microcenter

	for {
		for _, item := range stores.TwoGb {
			if item.Store == "adafruit" {
				// html := getStoreHtmlBody(item.Link)
				// log.Printf(html)
				fmt.Println(item.Store)
			} else if item.Store == "vilros" {
				isVilrosAvailable(item)
			} else if item.Store == "pishop.us" {
				fmt.Println(item.Store)
			}
		}

		time.Sleep(60 * time.Second)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func main() {
	token := "NjYzMTcwNTIyMDc4NTExMTA0.GEfq7G._azpPjUB_fajKlZi6VDgK7r7_pvRF3mrwdNj88"
	// data := readJsonFile()
	// testLoop(data)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
}
