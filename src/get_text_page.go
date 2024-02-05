package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
)

func getTextPage(client *http.Client, pageUrl string, outPath string) {
	res, err := client.Get(pageUrl)
	if err != nil {
		log.Fatalf("Error while getting text page %s", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error while parsing text page %s", err)
	}

	iframes := doc.Find("iframe")
	if iframes.Length() == 0 {
		log.Printf("No iframes found on text page %s", pageUrl)
	}

	iframes.Each(func(i int, selection *goquery.Selection) {
		iframeUrlStr, exists := selection.Attr("src")
		var finalPath = outPath
		if iframes.Length() > 1 {
			finalPath += fmt.Sprintf("_%d.mp4", i)
		} else {
			finalPath += ".mp4"
		}

		if !exists {
			log.Fatalf("Could not find src for iframe %d", i)
		}

		iframeUrl, err := url.Parse(iframeUrlStr)
		if err != nil {
			log.Fatalf("Error while parsing iframe url %s", err)
		}

		switch iframeUrl.Host {
		case "lecture2go.uni-hamburg.de":
			downloadLecture2Go(client, iframeUrlStr, finalPath)
		default:
			log.Printf("Unknown iframe host %s", iframeUrl.Host)
		}
	})
}
