package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

type Lecture2GoResource struct {
	Url string `json:"file"`
}

func downloadLecture2Go(client *http.Client, pageUrl string, outputPath string) {
	regex := regexp.MustCompile(`\Qconst uris = [].concat(\E(.*)\Q);\E`)
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

		resourcesJson := regex.FindStringSubmatch(content)[1]
		log.Printf("Found uris: %s", resourcesJson)

		// Extract the JSON array from the script tag
		var resources []Lecture2GoResource
		err := json.Unmarshal([]byte(resourcesJson), &resources)
		if err != nil {
			log.Fatalf("Error while parsing JSON %s", err)
			return
		}

		for _, resource := range resources {
			fileUrl := resource.Url
			log.Printf("Downloading %s", fileUrl)
			cmd := exec.Command("youtube-dl", "-o", outputPath, fileUrl)
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				fmt.Println("could not run command: ", err)
			}

			log.Printf("Downloaded %s", fileUrl)

			// One successful download is enough
			return
		}
	})
}
