package main

import (
	"encoding/json"
	"fmt"
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
}

// testing out the json library
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

func getStoreHtmlBody(Link string) string {
	resp, err := http.Get(Link)
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

func testLoop(pi PiItemsJson) {
	// chicacgo electonic distributors
	// adafruit
	// pishop.us
	// sparkfun
	// vilros
	// canakit
	// microcenter

	for {
		for _, value := range pi.TwoGb {
			if value.Store == "adafruit" {
				// html := getStoreHtmlBody(value.Link)
				// log.Printf(html)
				fmt.Println(value.Store)
			} else if value.Store == "vilros" {
				fmt.Println(value.Store)
			} else if value.Store == "pishop.us" {
				fmt.Println(value.Store)
			}
		}

		time.Sleep(60 * time.Second)
	}
}

func main() {
	data := readJsonFile()
    testLoop(data)
}
