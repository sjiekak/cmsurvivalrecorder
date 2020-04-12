package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" //init postgres driver
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var urlDonation = "https://cameroonsurvival.org/fr/dons/"
var frequency = 10 * time.Second

func main() {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)
	signal.Notify(osSignals, syscall.SIGTERM)

	log.Info("connecting to the database...")
	db, err := setupDB("postgres://test:test@:65432/cmsurvival?sslmode=disable")
	if err != nil {
		log.Fatal("error setting db : %w", err)
	}
	log.Info("connected")
	log.Info("scrapping website... ")
	defer db.Close()
	recordErrors := 0
	run := true
	for run && recordErrors < 10 {
		select {
		case s := <-osSignals:
			log.Warnf("Received signal %s", s.String())
			run = false
			break
		case <-time.Tick(frequency):
			v, err := getLatestValue()
			if err != nil {
				log.Error("couldn't get latest value : %w", err)
				recordErrors++
				continue
			}
			if err = writeValue(db, time.Now(), v); err != nil {
				log.Error("couldn't get latest value : %w", err)
				recordErrors++
				continue
			}
			recordErrors = 0
		}
	}
}

func getWebPageData(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("could not get page from server. status %d %s", resp.StatusCode, resp.Status)
	}
	return resp.Body, nil
}

func getLatestValue() (float64, error) {
	body, err := getWebPageData(urlDonation)
	if err != nil {
		if body != nil {
			body.Close()
		}
		return 0.0, err
	}
	defer body.Close()
	var doc *html.Node
	var raisedStr string
	var value float64
	if doc, err = html.Parse(body); err != nil {
		return 0.0, err
	}
	if raisedStr, err = crawl(doc); err != nil {
		return 0.0, err
	}
	value, err = parseAmount(raisedStr)
	// the website is in french, and numbers use the dot '.' for thousands
	return value * 1000, err
}
