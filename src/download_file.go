package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func downloadFile(client *http.Client, url string, outputPath string) {
	log.Printf("Downloading %s", outputPath)

	contentTypeExtMap := map[string]string{
		"application/pdf": ".pdf",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	}

	parentPath := outputPath[:len(outputPath)-len(outputPath[strings.LastIndex(outputPath, "/"):])]
	err := os.MkdirAll(parentPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while creating directory %s", err)
	}

	res, err := client.Get(url)
	if err != nil {
		log.Fatalf("Error while getting file %s", err)
	}
	defer res.Body.Close()

	ext := contentTypeExtMap[res.Header.Get("Content-Type")]

	out, err := os.Create(outputPath + ext)
	if err != nil {
		log.Fatalf("Error while creating file %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Fatalf("Error while writing file %s", err)
	}
}
