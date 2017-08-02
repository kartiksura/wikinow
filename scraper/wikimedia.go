package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Response is a struct used for unmarshaling the MediaWiki JSON response.
type Response struct {
	Query struct {
		// The JSON response for this part of the struct is dumb.
		// It will return something like { '23': { 'pageid': 23 ...
		//
		// As a workaround you can use PageSlice which will create
		// a list of pages from the map.
		Pages map[string]Page
	}
}

// A Page represents a MediaWiki page and its metadata.
type Page struct {
	Pageid    int
	Ns        int
	Title     string
	Touched   string
	Lastrevid int
	// Mediawiki will return '' for zero, this makes me sad.
	// If for some reason you need this value you'll have to
	// do some type assertion sillyness.
	Counter   interface{}
	Length    int
	Edittoken string
	Revisions []struct {
		Revid         int       `json:"revid"`
		Parentid      int       `json:"parentid"`
		Minor         string    `json:"minor"`
		User          string    `json:"user"`
		Userid        int       `json:"userid"`
		Timestamp     time.Time `json:"timestamp"`
		Size          int       `json:"size"`
		Sha1          string    `json:"sha1"`
		ContentModel  string    `json:"contentmodel"`
		Comment       string    `json:"comment"`
		ParsedComment string    `json:"parsedcomment"`
		ContentFormat string    `json:"contentformat"`
		Body          string    `json:"*"` // Take note, MediaWiki literally returns { '*':
	}
	Imageinfo []struct {
		Url            string
		Descriptionurl string
	}
}

func mediaWikiCall(title string) (op string, err error) {

	param := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&titles=%s&prop=revisions&rvprop=content&format=json", title)

	u, err := url.Parse(param)

	q := u.Query()
	u.RawQuery = q.Encode()

	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return
	}
	var response Response
	err = json.Unmarshal(htmlData, &response)
	if err != nil {
		return "", err
	}

	if len(response.Query.Pages) != 1 {
		return "", errors.New("received unexpected number of pages")
	}

	// we use a hacky way of extracting the map's lone value
	var page *Page
	for _, pg := range response.Query.Pages {
		page = &pg
	}
	if len(page.Revisions) > 0 {
		return page.Revisions[0].Body, nil

	}
	return "", nil

}
