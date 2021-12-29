package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"news-aggregator/pkg/storage"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// RSS is a root element of RSS feed
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel is a channel element of RSS feed
type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"string"`
	PubTime     string   `xml:"pubDate"`
	Items       []Item   `xml:"item"`
}

// Item represents single blog post in RSS feed
type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
}

// Parse fetches posts from URL of RSS feed
func Parse(url string) ([]storage.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// It would be easier to write tests for xml parser
	posts, err := parse(raw)
	return posts, err
}

// parse receives raw xml bytes and returns slice of posts
func parse(raw []byte) ([]storage.Post, error) {
	var rss RSS
	err := xml.Unmarshal(raw, &rss)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	for _, i := range rss.Channel.Items {
		var post storage.Post
		post.Title = i.Title
		// Sanitizing html tags from text
		sanitizer := bluemonday.StripTagsPolicy()
		post.Content = sanitizer.Sanitize(i.Description)
		// Removing unnecessary information
		post.Content = strings.ReplaceAll(post.Content, " Читать далее", "")
		post.Content = strings.ReplaceAll(post.Content, " Читать дальше", "")
		post.Content = strings.ReplaceAll(post.Content, "&#34;", "")
		post.Content = strings.ReplaceAll(post.Content, "  ", " ")
		// Parsing time
		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", i.PubDate)
		if err != nil {
			t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 GMT", i.PubDate)
		}
		if err != nil {
			return nil, err
		}
		post.PubTime = t.Unix()
		post.Link = i.Link
		posts = append(posts, post)
	}
	return posts, nil
}
