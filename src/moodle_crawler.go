package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"time"
)

func main() {
	basePath := flag.String("base-path", "", "Base path to download files to")
	moodleUsername := flag.String("moodle-username", "", "Moodle username")
	moodlePassword := flag.String("moodle-password", "", "Moodle password")
	flag.Parse()

	if *basePath == "" || *moodleUsername == "" || *moodlePassword == "" {
		fmt.Println("Missing required argument. Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cookieJar, _ := cookiejar.New(nil)
	c := http.Client{Timeout: time.Duration(30) * time.Second, Jar: cookieJar}

	sessionKey := login(&c, *moodleUsername, *moodlePassword)

	courses := getCourses(&c, sessionKey).Courses

	for _, course := range courses {
		log.Printf("Course: %s", course.FullName)
		splitRegex := regexp.MustCompile(`[/ ]`)
		outDir := *basePath + "/" + splitRegex.Split(course.Shortname, -1)[0]

		resources := getCourseResources(&c, sessionKey, course.Id)
		for _, resource := range resources.Cm {
			outPath := makePath(outDir, resource.Name)
			switch resource.ModName {
			case "Datei":
				if !resource.UserVisible {
					continue
				}
				downloadFile(&c, resource.Url, outPath)
			case "Textseite":
				getTextPage(&c, resource.Url, outPath)
			case "Aufgabe":
				getAufgabe(&c, resource.Url, outPath)
			case "Verzeichnis":
				getVerzeichnis(&c, resource.Url, outPath)
			case "Forum", "Freie Gruppeneinteilung", "Moodleoverflow", "Feedback", "Text- und Medienfeld", "Abstimmung", "Planer", "Kollaboratives Dokument", "Anwesenheit", "Link/URL", "Test":
				// Do nothing for now. Maybe add support for these later.
			default:
				log.Fatalf("Unknown resource type %s: %s", resource.ModName, resource.Name)
			}
		}
	}
}
