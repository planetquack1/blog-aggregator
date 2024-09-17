package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

type RSS struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchRSSFeed(url string) (RSS, error) {

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return RSS{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error") // TODOD
		return RSS{}, err
	}

	// Parse the RSS XML
	var rss RSS
	err = xml.Unmarshal(dat, &rss)
	if err != nil {
		log.Println("error") // TODOD
		return RSS{}, err
	}

	return rss, nil

}
