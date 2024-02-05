package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
)

func makeFormValues(form *goquery.Selection) url.Values {
	formValues := url.Values{}
	form.Find("input").Each(func(i int, selection *goquery.Selection) {
		name, exists := selection.Attr("name")
		if !exists {
			log.Printf("Could not find name for input %d", i)
			return
		}

		value, exists := selection.Attr("value")
		if !exists {
			log.Printf("Could not find value for input %d", i)
			value = ""
		}

		formValues.Set(name, value)
	})

	return formValues
}
