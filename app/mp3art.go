package main

import (
	"log"

	"mp3art/Mp3Art/jobs"
)

var (
	baseURL    	= "https://open.spotify.com/search"
	searchURL  	= "https://api.spotify.com/v1/search"
	accessToken	= ""
)

// "api.spotify.com/v1/search?type=album%2Cartist%2Cplaylist%2Ctrack%2Cshow_audio%2Cepisode_audio&q=dualipa*&decorate_restrictions=false&best_match=true&include_external=audio&limit=50&userless=true&market=TR"

func main() {
	//accessToken = jobs.GetAccessToken( baseURL )
	//market := "TR"
	//var name string
	var save = true
	for {
		/* fmt.Printf("Give me the search string: ")
		_, err := fmt.Scanln(&name); if err != nil {
			log.Println("Couldn't get the string:",err)
		}
		jobs.GetMusicList( searchURL, accessToken, market, name) */
		path, err := jobs.FixFileNames("",save); if err != nil {
			log.Println("Couldn't fix file names at "+ path + "\n" + err.Error())
		}
		break
	}

}