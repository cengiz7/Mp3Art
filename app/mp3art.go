package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"../jobs"
)

var (
	baseURL    	= "https://open.spotify.com/search"
	searchURL  	= "https://api.spotify.com/v1/search"
	accessToken	= ""
)

// "api.spotify.com/v1/search?type=album%2Cartist%2Cplaylist%2Ctrack%2Cshow_audio%2Cepisode_audio&q=dualipa*&decorate_restrictions=false&best_match=true&include_external=audio&limit=50&userless=true&market=TR"

func init() {
	// create musics folder if not exists
	root, err := jobs.GetRootPath(); if err != nil {
		log.Println()
	} else {
		root += "/musics"
		if _, err := os.Stat(root); os.IsNotExist(err) {
			err = os.Mkdir(root,  os.ModePerm);  if err != nil {
				log.Println("Failed to create /musics folder.", err)
			}
		}
	}
}

func setAlbumArt( setArtCh <- chan map[string][]byte ){
 	// set album arts
	for {
		time.Sleep(time.Second * 100)
	}
}

func getAlbumArt( getArtCh <- chan map[string]string, setArtCh chan <- map[string][]byte ){
	// get album art and send it to the setalbumart ch
	for {
		time.Sleep(time.Second * 100)
	}
}

func parseJSON ( getQueryResultCh <- chan map[string][]byte , getArtCh chan <- map[string]string ) {
	// TODO : parse json and send to the getart ch
	for data := range getQueryResultCh {
		var jsonMap map[string]interface{}
		for path, jsonData := range data {
			err := json.Unmarshal(jsonData, &jsonMap); if err != nil {
				log.Println("Error while processing json bytes to map:",err.Error())
			} else {
				fmt.Printf("%#v\n", jsonMap)
				fmt.Println(path)
			}
		}
		time.Sleep(time.Second * 100)
	}
}

func main() {
	accessToken = jobs.GetAccessToken( baseURL )
	localZone 	:= "TR"
	save  		:= true

	// carries mp3 file path and GET response data
	getQueryResultChannel 	:= make(chan map[string][]byte)
	// carries mp3 file path and image url
	getArtChannel 			:= make(chan map[string]string)
	// carries mp3 file path and image data
	setArtChannel 			:= make(chan map[string][]byte)
	// carries parameters for GET query request
	paramsChannel 			:= make(chan []string)

	go jobs.GetQueryResult( paramsChannel,getQueryResultChannel )
	go parseJSON  ( getQueryResultChannel,getArtChannel )
	go getAlbumArt( getArtChannel, setArtChannel )
	go setAlbumArt( setArtChannel )

	musicListMap, path, err := jobs.FixFileNames("",save); if err != nil {
		log.Println("Couldn't fix file names at "+ path + "\n" + err.Error())
	} else {
		for path, name := range musicListMap {
			paramsChannel <- []string{accessToken,searchURL,localZone,name,path}
			time.Sleep(time.Second * 100)
			break
		}
	}

	for {
		/* fmt.Printf("Give me the search string: ")
		_, err := fmt.Scanln(&name); if err != nil {
			log.Println("Couldn't get the string:",err)
		}
		jobs.GetMusicList( searchURL, accessToken, market, name) */

		break
	}
}