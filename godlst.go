package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./godlst <streamable url>")
		return
	}

	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	title := time.Now().Format("2006-01-02 15:04:05")
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		parts := strings.Split(s.Text(), "Watch ")
		if len(parts) > 1 {
			parts = strings.Split(parts[1], " | Streamable")
		}
		title = strings.Join(parts, " ")
	})

	doc.Find("meta[property='og:video:url']").Each(func(i int, s *goquery.Selection) {
		source, exists := s.Attr("content")
		if exists {
			downloadVideo(source, title+".mp4")
		} else {
			log.Println("No video source found")
		}
	})
}

func downloadVideo(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download video: %v", err)
	}
	defer resp.Body.Close()

	output, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	_, err = io.Copy(output, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save video: %v", err)
	}

	return nil
}
