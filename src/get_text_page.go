package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
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

	log.Printf("Downloading text page %s", pageUrl)

	links := doc.Find("#page-content a")
	links.Each(func(i int, selection *goquery.Selection) {
		linkUrl, exists := selection.Attr("href")
		linkText := selection.Text()
		if !exists {
			log.Fatalf("Could not find href for link %d", i)
		}

		linkUrlParsed, err := url.Parse(linkUrl)
		if err != nil {
			log.Fatalf("Error while parsing link url %s", err)
		}

		switch linkUrlParsed.Host {
		case "lecture2go.uni-hamburg.de":
			downloadLecture2Go(client, linkUrl, outPath)
		case "lernen.min.uni-hamburg.de":
			fmt.Println(linkUrlParsed.Path)
			if strings.HasPrefix(linkUrlParsed.Path, "/pluginfile.php") {
				fileName := linkUrlParsed.Path[strings.LastIndex(linkUrlParsed.Path, "/")+1:]
				downloadFile(client, linkUrl, makePath(outPath, fileName))
			} else {
				log.Printf("Unknown link path %s: %s; Title: %s", linkUrlParsed.Path, linkUrl, linkText)
			}
		default:
			log.Printf("Unknown link host %s", linkUrlParsed.Host)
		}
	})

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
