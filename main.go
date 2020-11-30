package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseWiki(url string) []string {
        // Make Get request to a url
        resp, _ := http.Get(url)

        // Close the request
        defer resp.Body.Close()

        // Create a goquery doc
        doc, _ := goquery.NewDocumentFromReader(resp.Body)

        // Create array for tracks
        var tracks []string
        var track string

        tracks = append(tracks, "artist,song\n")

        // Find all songs on page and parse string into artist and song
        doc.Find(".div-col").Each(func(_ int, s *goquery.Selection) {
                s.Find("li").Each(func(_ int, t *goquery.Selection) {
                        // Split string by "-" to separate artist and track
                        text := strings.Split(t.Text(), " –")

                        // Grab the artist and song
                        artist := text[0]
                        song := strings.Trim(text[1], " \"")

                        // Create track
                        track = artist + "," + song + "\n"

                        tracks = append(tracks, track)
                })
        })

        return tracks
}


func main() {

	// The Pitchfork 500 wiki page url
	url := "https://en.wikipedia.org/wiki/The_Pitchfork_500"

	// Output filename
	out_filename := "tracks.csv"

	// Grab the tracker from the parser
	tracks := ParseWiki(url)

	// Create the CSV to save the tracks
	file, _ := os.Create(out_filename)

	defer file.Close()

	var err error 	

	//Iterate through tracks in slice, writing them to the CSV
	for _, track := range tracks {
		_, err = io.WriteString(file, track)

		if err != nil {
			log.Fatal(err)
		}

		file.Sync()
	}
}