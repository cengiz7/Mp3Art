package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"../jobs"
)

var (
	baseURL    	= "https://open.spotify.com/search"
	searchURL  	= "https://api.spotify.com/v1/search"
	accessToken	= ""
)

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

func dumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			dumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}



func parseJSON ( getQueryResultCh <- chan jobs.GetQueryResultStruct ,paramsCh chan <- []string, getArtCh chan <- map[string]string ) {
	// TODO : parse json and send to the getart ch
	for resultStruct := range getQueryResultCh {
		var jsonMap map[string]interface{}
		for path, jsonData := range resultStruct.Result {
			err := json.Unmarshal(jsonData, &jsonMap); if err != nil {
				log.Println("Error while processing json bytes to map:",err.Error())
			} else {
				//fmt.Printf("%#v\n", jsonMap)
				//fmt.Println(jsonMap["best_match"])
				//_ , _ = dumpMap( jsonMap , true)
				items := jsonMap["best_match"].(map[string]interface{})
				for _ , a := range items["items"].([]interface{}) {
					dumpMap("", a.(map[string]interface{}))
				}
				// we will continue
				if _, ok := items["items"].(map[string]interface{}); ok || 1 == 1 {
					path = path
				} else {
					fmt.Println("Couldn't found anything with name => ",resultStruct.Params[3])
					if name, ok := jobs.TrimFromEnd(resultStruct.Params[3]); ok {
						resultStruct.Params[3] = name
						time.Sleep(10 * time.Second)
						paramsCh <- resultStruct.Params
					}
				}
			}
		}
	}
}

func main() {
	accessToken = jobs.GetAccessToken( baseURL )
	localZone 	:= "TR"
	limit 		:= "1" // list object count for search query
	save  		:= true
	// ----------- Channel declarations ------------ //
	// carries struct for params slice and GET response data
	getQueryResultChannel 	:= make(chan jobs.GetQueryResultStruct,  1000)
	// carries mp3 file path and image url
	getArtChannel 			:= make(chan map[string]string, 100)
	// carries mp3 file path and image data
	setArtChannel 			:= make(chan map[string][]byte, 500)
	// carries parameters for GET query request
	paramsChannel 			:= make(chan []string, 5000)

	// --------------- Goroutines ------------------ //
	go jobs.GetQueryResult( paramsChannel,getQueryResultChannel )
	go parseJSON  ( getQueryResultChannel,paramsChannel, getArtChannel )
	go getAlbumArt( getArtChannel, setArtChannel )
	go setAlbumArt( setArtChannel )

	musicListMap, path, err := jobs.FixFileNames("",save); if err != nil {
		log.Println("Couldn't fix file names at "+ path + "\n" + err.Error())
	} else {
		for path, name := range musicListMap {
			name = strings.Replace(name,".mp3","",-1)
			paramsChannel <- []string{accessToken,searchURL,localZone,name,limit,path}
		}
		time.Sleep(time.Second * 100)
	}
}