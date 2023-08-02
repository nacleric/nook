package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
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
		fmt.Println(err)
	}

	var data ItemsJson

	if err := json.Unmarshal(content, &data); err != nil {
		fmt.Println(err)
	}

	return data
}

func getStoreResponseBody(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	html := string(body)
	return html
}

func isVilrosAvailable(i Item) (bool, error) {
	resp, err := http.Get(i.Misc)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
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

func isPiShopAvailable(i Item) (bool, error) {
	var isAvailable bool = true

	res_body := getStoreResponseBody(i.Link)
	doc, err := html.Parse(strings.NewReader(res_body))

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, a := range n.Attr {
				if strings.ToLower(a.Val) == "out of stock" {
					isAvailable = false
					break
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
			var isAvailable bool
			var err error
			switch item.Store {
			case "adafruit":
				isAvailable, err = isAdaFruitAvailable(item)
				if err != nil {
					fmt.Println("adafruit error:", err)
				}
			case "pishop.us":
				isAvailable, err = isPiShopAvailable(item)
				if err != nil {
					fmt.Println("pishop.us error:", err)
				}
			case "vilros":
				isAvailable, err = isVilrosAvailable(item)
				if err != nil {
					fmt.Println("vilros error", err)
				}
			}

			if isAvailable {
				notifyOwner(dg, item.Link)
			}
		}

		time.Sleep(60 * time.Second)
	}
}

func initBot() *discordgo.Session {
	token := os.Getenv("TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
	}

	return dg
}

func initBotCommands(dg *discordgo.Session) {
	dg.AddHandler(pingMe)
	dg.AddHandler(youtubeDownloadMp3)
	dg.AddHandler(whenPartyAnimals)
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
	// stores := readJsonFile()
	dg := initBot()
	// go pollStores(stores, dg)
	initBotCommands(dg)
	discordServerListener(dg)
}

func notifyOwner(s *discordgo.Session, msg string) {
	userId := os.Getenv("USERID")
	channelId := os.Getenv("CHANNELID")
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

	if m.Content == "$ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func whenPartyAnimals(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "$wen" {
		currentTime := time.Now()
		partyAnimalsRelease := time.Date(2023, 9, 20, 0, 0, 0, 0, time.UTC)
		diff := int(math.Round(partyAnimalsRelease.Sub(currentTime).Hours() / 24))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d days till Party Animals", diff))
	}
}

func youtubeDownloadMp3(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	removeExtraWhitespace := regexp.MustCompile(`\s+`).ReplaceAllString(m.Content, " ")
	params := strings.Split(removeExtraWhitespace, " ")

	if params[0] == "$dl" && len(params) == 1 {
		s.ChannelMessageSend(m.ChannelID, "$dl requires link to youtube video, expecting `$dl <yt link>`")
		return
	}

	if params[0] == "$dl" && len(params) > 2 {
		s.ChannelMessageSend(m.ChannelID, "Too many arguments, expecting `$dl <yt link>`")
		return
	}

	if params[0] == "$dl" && len(params) == 2 {
		// TODO check if it's a valid youtube link
		_, err := url.ParseRequestURI(params[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Not a valid URL")
			return
		}

		// Note: requestId also acts as a folder
		requestId := uuid.NewString()
		if err := os.Mkdir(requestId, os.ModePerm); err != nil {
			log.Printf("ReqID %s | err: %s", requestId, err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Preparing to Download Music... ReqID: %s", requestId))
		if err := downloadMusic(params[1], requestId); err != nil {
			log.Printf("Trouble downloading Music | ReqID %s | err: %s", requestId, err)
			s.ChannelMessageSend(m.ChannelID, "Trouble downloading Music")
			if err := os.RemoveAll(requestId); err != nil {
				log.Println(err)
			}
		}

		files, err := ioutil.ReadDir(requestId)
		if err != nil {
			log.Printf("Directory might not exist | ReqID %s | err: %s", requestId, err)
			return
		}

		if len(files) == 1 {
			fileName := files[0].Name()
			fileLocation := filepath.Join(requestId, fileName)
			mp3File, err := os.Open(fileLocation)
			if err != nil {
				log.Printf("mp3File might not exist | ReqID %s | err: %s", requestId, err)
			}
			defer mp3File.Close()

			data := discordgo.MessageSend{Content: "Here you go", File: &discordgo.File{Name: fileName, ContentType: "mp3", Reader: mp3File}}
			s.ChannelMessageSendComplex(m.ChannelID, &data)
		} else {
			zipName := fmt.Sprintf("%s.zip", requestId)
			zipFile, err := os.Create(zipName)
			if err != nil {
				log.Printf("Trouble creating zipFile | ReqID %s | err: %s", requestId, err)
			}
			defer zipFile.Close()

			zipw := zip.NewWriter(zipFile)
			defer zipw.Close()

			for _, file := range files {
				fileLocation := filepath.Join(requestId, file.Name())
				mp3File, err := os.Open(fileLocation)
				if err != nil {
					log.Printf("ReqID %s | err: %s", requestId, err)
				}
				defer mp3File.Close()

				wr, err := zipw.Create(file.Name())
				if err != nil {
					log.Printf("Failed to create zipWriter | ReqID %s | err: %s", requestId, err)
				}

				if _, err := io.Copy(wr, mp3File); err != nil {
					log.Printf("Failed to write %s to zip | ReqID %s | err: %s", file.Name(), requestId, err)
				}
			}

			zipData, err := os.Open(zipName)
			if err != nil {
				log.Println(err)
			}
			data := discordgo.MessageSend{File: &discordgo.File{Name: zipName, ContentType: "application/zip", Reader: zipData}}
			s.ChannelMessageSendComplex(m.ChannelID, &data)

			// Clean up zipfile
			if err := os.Remove(zipName); err != nil {
				log.Println(err)
			}
		}
		// Clean up
		if err := os.RemoveAll(requestId); err != nil {
			log.Println(err)
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Task Finished for %s", requestId))
	}
}

func downloadMusic(ytlink string, requestId string) error {
	// Notes: https://www.reddit.com/r/youtubedl/comments/rh893t/comment/hopd2b5/?utm_source=share&utm_medium=web2x&context=3
	dlDirectory := filepath.Join(requestId, "%(title)s.%(ext)s")
	cmd := exec.Command("yt-dlp", "-x", "--audio-format=mp3", ytlink, "-o", dlDirectory)
	_, err := cmd.Output()
	return err
}
