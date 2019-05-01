
package main

import (
	"golang.org/x/exp/errors/fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sclevine/agouti"
)


const (
	baseURL   = "https://open.spotify.com/search"
	searchURL = "https://api.spotify.com/v1/search"
)


// "api.spotify.com/v1/search?type=album%2Cartist%2Cplaylist%2Ctrack%2Cshow_audio%2Cepisode_audio&q=dualipa*&decorate_restrictions=false&best_match=true&include_external=audio&limit=50&userless=true&market=TR"


func main() {
	driver := agouti.ChromeDriver(agouti.ChromeOptions("args", []string{"--headless", "--disable-gpu", "--no-sandbox"}), )
	client := &http.Client{}

	if err := driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}

	page, err := driver.NewPage(); if err != nil {
		log.Fatal("Failed to open page:", err)
	}

	if err := page.Navigate(baseURL); err != nil {
		log.Fatal("Failed to navigate:", err)
	}

	cookies, err  := page.GetCookies(); if err != nil {
		log.Fatal("Failed to get cookies: ", err)
	}

	if err := driver.Stop(); err != nil {
		log.Fatal("Failed to close pages and stop WebDriver:", err)
	}

	req, err := http.NewRequest("GET", searchURL ,nil); if err != nil {
		log.Fatal("Failed to close pages and stop WebDriver:", err)
	}

	q := req.URL.Query()
	q.Add("decorate_restrictions", "false")
	q.Add("include_external", "audio")
	q.Add("best_match", "true")
	q.Add("userless", "true")
	q.Add("market", "TR")
	q.Add("limit", "5")
	q.Add("type", "artist,playlist,track,show_audio,episode_audio")
	q.Add("q", "dualipa*")

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	for _, cookie := range cookies {
		if strings.Contains(cookie.String(), "wp_access_token") {
			token := strings.Split(cookie.String(),";")[0]
			token = strings.Split(token,"=")[1]
			fmt.Println(token)
			req.Header.Set("Authorization","Bearer " +token)
		}
	}


	resp, err := client.Do(req)
	// Read response
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))


}