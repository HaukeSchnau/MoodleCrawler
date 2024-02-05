package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
)

func login(client *http.Client, username string, password string) string {
	res, err := client.Get("https://lernen.min.uni-hamburg.de/login/index.php")
	if err != nil {
		log.Fatalf("Error while getting login page %s", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error while parsing login page %s", err)
	}

	identityProviderUrl, exists := doc.Find(".login-identityprovider-btn").Attr("href")
	if !exists {
		log.Fatalf("Could not find identity provider button")
	}

	res, err = client.Get(identityProviderUrl)
	if err != nil {
		log.Fatalf("Error while getting identity provider page %s", err)
	}
	defer res.Body.Close()

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error while parsing identity provider page %s", err)
	}

	csrfToken, exists := doc.Find("input[name='csrf_token']").Attr("value")
	if !exists {
		log.Fatalf("Could not find csrf token")
	}

	res, err = client.PostForm("https://login.uni-hamburg.de/idp/profile/SAML2/Redirect/SSO?execution=e1s1", url.Values{
		"j_username":       {username},
		"j_password":       {password},
		"csrf_token":       {csrfToken},
		"_eventId_proceed": {""},
	})
	if err != nil {
		log.Fatalf("Error while logging in %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Error while logging in, status code %d", res.StatusCode)
	}

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error while parsing continue page %s", err)
	}

	form := doc.Find("form")
	action, exists := form.Attr("action")
	if !exists {
		log.Fatalf("Could not find action")
	}

	formValues := makeFormValues(form)

	res, err = client.PostForm(action, formValues)
	if err != nil {
		log.Fatalf("Error while continuing %s", err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("Error while continuing, status code %d", res.StatusCode)
	}

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error while parsing final page %s", err)
	}

	sessionKey, exists := doc.Find("input[name='sesskey']").Attr("value")
	if !exists {
		log.Fatalf("Could not find session key")
	}

	return sessionKey
}
