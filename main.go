package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Config .
type Config struct {
	Webhook          string
	Minutes          int
	SendIfEmptyCache bool
}

// Cache .
type Cache struct {
	LastVersion   string
	LastChangelog string
}

var config Config
var cache Cache

func interval(what func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func log(msg string) {
	fmt.Println("[" + formatTime(time.Now()) + "] [DiUW] " + msg)
}

func check() {
	log("Checking for new update..")
	doc, err := goquery.NewDocument("https://apps.apple.com/us/app/discord/id985746746")
	if err != nil {
		log("[ERROR] " + err.Error())
	} else {
		appVer := strings.Split(doc.Find(".whats-new__latest__version").Text(), " ")[1]
		appChangelog, _ := doc.Find(".whats-new__content p").Last().Html()
		appChangelog = strings.Replace(html.UnescapeString(appChangelog), "<br/>", "\n", -1)

		if cache.LastVersion != "" || config.SendIfEmptyCache {
			if cache.LastVersion != appVer || cache.LastChangelog != appChangelog {
				icon, _ := doc.Find(".we-artwork__source").First().Attr("srcset")
				icon = strings.Split(icon, " ")[0]
				bd := map[string]interface{}{"embeds": []map[string]interface{}{map[string]interface{}{
					"author": map[string]string{
						"name":     "Discord",
						"icon_url": icon,
						"url":      "https://apps.apple.com/us/app/discord/id985746746",
					},
					"title":       "New version: **" + appVer + "**",
					"description": appChangelog,
					"footer": map[string]string{
						"text": "Updated " + doc.Find(".whats-new__latest time").Text(),
					},
					"color": 7506394,
				}}}
				body, _ := json.Marshal(bd)
				log("Found new version " + appVer + ", Sending")
				http.Post(config.Webhook, "application/json", bytes.NewReader(body))
			}
		}

		if cache.LastVersion != appVer || cache.LastChangelog != appChangelog {
			cache.LastVersion = appVer
			cache.LastChangelog = appChangelog
			b, _ := json.Marshal(cache)
			ioutil.WriteFile("cache.json", b, 0644)
		}
	}
}

func main() {
	jfile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		if os.IsNotExist(err) {
			return
		}
	}
	defer jfile.Close()
	byteV, _ := ioutil.ReadAll(jfile)
	json.Unmarshal(byteV, &config)

	jfile, err = os.Open("cache.json")
	if err != nil {
		if os.IsNotExist(err) == false {
			fmt.Println(err)
		}
	}
	defer jfile.Close()
	byteV, _ = ioutil.ReadAll(jfile)
	json.Unmarshal(byteV, &cache)

	log("Started | Interval: " + fmt.Sprint(config.Minutes) + " minute/s")

	<-interval(check, time.Duration(config.Minutes)*time.Minute)
}
