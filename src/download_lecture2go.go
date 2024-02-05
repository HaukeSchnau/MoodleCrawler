package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

func downloadLecture2Go(client *http.Client, pageUrl string, outputPath string) {
	regex := regexp.MustCompile(`\Qconst uris = [].concat([{"file":"\E(.+)\Q"}]);\E`)
	res, err := client.Get(pageUrl)
	if err != nil {
		log.Fatalf("Error while getting text page %s", err)
	}
	defer res.Body.Close()

	log.Printf("Downloading %s", pageUrl)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error while parsing text page %s", err)
	}

	scriptTags := doc.Find("script")
	scriptTags.Each(func(i int, selection *goquery.Selection) {
		content := selection.Text()
		if !regex.MatchString(content) {
			return
		}

		matches := regex.FindStringSubmatch(content)
		if len(matches) != 2 {
			log.Fatalf("Could not find file in %s", content)
		}

		fileUrl := matches[1]

		cmd := exec.Command("youtube-dl", "-o", outputPath, fileUrl)
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Println("could not run command: ", err)
		}

		log.Printf("Downloaded %s", fileUrl)
	})
}
