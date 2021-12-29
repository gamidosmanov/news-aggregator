package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"news-aggregator/pkg/api"
	"news-aggregator/pkg/rss"
	"news-aggregator/pkg/storage"
	"time"
)

// server encapsulates web-server internals
type server struct {
	db  *storage.DB
	api *api.API
}

// config encapsulates server configuration
type config struct {
	Links  []string `json:"links"`
	Period int      `json:"requestPeriod"`
}

// shutdown stops source parsing
func shutdown(done chan<- struct{}) {
	done <- struct{}{}
}

func main() {
	// Read config file
	confRaw, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	var conf config
	err = json.Unmarshal(confRaw, &conf)
	if err != nil {
		log.Fatal(err)
	}

	// Create database connection
	db, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	// Set up server
	var srv server
	srv.db = db
	srv.api = api.New(srv.db)

	// Launching routines
	errChan := make(chan error, 10)
	postChan := make(chan storage.Post, 10)
	done := make(chan struct{})
	// Graceful shutdown of goroutines
	defer shutdown(done)

	// Creating separate routine for each news source
	for _, link := range conf.Links {
		go func(link string) {
			for {
				select {
				case <-done:
					return
				default:
					log.Printf("Pasring feed: %v", link)
					posts, err := rss.Parse(link)
					if err != nil {
						errChan <- fmt.Errorf("error on %v: %v", link, err)
					}
					for _, p := range posts {
						postChan <- p
					}
					time.Sleep(time.Duration(conf.Period) * time.Minute)
				}
			}
		}(link)
	}

	// main routine
	go func() {
		for {
			select {
			case <-done:
				return
			case p := <-postChan:
				id, err := db.SavePost(p)
				if err != nil {
					errChan <- err
				}
				log.Printf("Post saved with id %d", id)
			case e := <-errChan:
				log.Print(e)
			}
		}
	}()

	// server startup
	err = http.ListenAndServe(":80", srv.api.Router())

	// on error shutdown
	if err != nil {
		shutdown(done)
		log.Fatal(err)
	}

}
