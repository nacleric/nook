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
	"strings"
	"syscall"
	"time"

	"golang.org/x/net/html"

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

func getStoreResponseBody(link string) string {
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
		return false, errors.New("json mapping for Vilros did not return boolean value")
	}
	return isAvailable, nil
}

func isAdaFruitAvailable(i Item) (bool, error) {
	var isAvailable bool = true

	res_body := getStoreResponseBody(i.Link)
	doc, err := html.Parse(strings.NewReader(res_body))

	// https://go.dev/play/p/sJqlctpSGQA
	// Traversing dom-tree recursively
	// Looking for <div itemprop="availability" class="oos-header">Out of stock</div>
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Val == "availability" {
					if strings.ToLower(n.FirstChild.Data) == "out of stock" {
						isAvailable = false
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return isAvailable, err
}

func pollStores(stores ItemsJson, dg *discordgo.Session) {
	for {
		for _, item := range stores.TwoGb {
			switch item.Store {
			case "adafruit":
				isAvailable, err := isAdaFruitAvailable(item)
				if err != nil {
					fmt.Println("adafruit error:", err)
				}

				if isAvailable {
					notifyEric(dg, item.Link)
				}
			case "pishop.us":
				fmt.Println(item.Store)
			case "vilros":
				isAvailable, err := isVilrosAvailable(item)
				if err != nil {
					fmt.Println("vilros error", err)
				}

				if isAvailable {
					notifyEric(dg, item.Link)
				}
			}
		}

		time.Sleep(60 * time.Second)
	}
}

func initBot() *discordgo.Session {
	token := "NjYzMTcwNTIyMDc4NTExMTA0.GEfq7G._azpPjUB_fajKlZi6VDgK7r7_pvRF3mrwdNj88"

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}

	return dg
}

func initBotCommands(dg *discordgo.Session) {
	dg.AddHandler(pingMe)
	dg.AddHandler(dmUser)
}

func discordServerListener(dg *discordgo.Session) {
	// Receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func main() {
	data := readJsonFile()
	dg := initBot()
	go pollStores(data, dg)
	initBotCommands(dg)
	discordServerListener(dg)
}

func notifyEric(s *discordgo.Session, msg string) {
	userId := "115272309870297090"
	channelId := "1090462882317221960"
	_, err := s.UserChannelCreate(userId)
	if err != nil {
		fmt.Println("error creating channel:", err)
		return
	}

	_, err = s.ChannelMessageSend(channelId, msg)
	if err != nil {
		fmt.Println("error sending DM message:", err)
	}
}

// Example code
func pingMe(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

// Example code
func dmUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "ping" {
		return
	}

	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		// Failed to create channel error
		// Possible reasons
		// 1. We don't share a server with the user
		// 2. Ratelimit
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, "Pong!")
	if err != nil {
		// Failed to send messages
		// Possible reasons
		// 1. User may have disabled DM's
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}
