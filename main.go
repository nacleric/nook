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

func isVilrosAvailable(i Item) {
	resp, err := http.Get(i.Misc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.Body)
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

func main() {
	data := readJsonFile()
	testLoop(data)
}
