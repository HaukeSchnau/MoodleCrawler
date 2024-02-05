package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

func getVerzeichnis(client *http.Client, pageUrl string, outPath string) {
	res, err := client.Get(pageUrl)
	if err != nil {
		log.Fatalf("Error while getting text page %s", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error while parsing text page %s", err)
	}

	anchors := doc.Find(".fp-filename a")
	if anchors.Length() == 0 {
		log.Fatalf("Could not find any download links")
	}

	anchors.Each(func(i int, selection *goquery.Selection) {
		downloadUrl, exists := selection.Attr("href")
		if !exists {
			log.Fatalf("Could not find href for anchor %d", i)
		}

		fullPath := outPath + "/" + selection.Text()

		downloadFile(client, downloadUrl, fullPath)
	})
}
