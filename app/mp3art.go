package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/bogem/id3v2"
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

func main() {
	accessToken = jobs.GetAccessToken( baseURL )
	localZone 	:= "TR"
	limit 		:= "1" // list object count for search query
	save  		:= true
	// ----------- Channel declarations ------------ //
	// carries struct for params slice and GET response data
	getQueryResultChannel 	:= make(chan jobs.GetQueryResultStruct,  1000)
	// carries mp3 file path and image url
	getArtChannel 			:= make(chan map[string]string, 500)
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

func setAlbumArt( setArtCh <- chan map[string][]byte ){
 	// set album arts
	for image := range setArtCh{
		for path, img := range image {
			if mp3File, err := id3v2.Open(path, id3v2.Options{Parse: true}); err != nil {
				log.Println("setAlbumArt err while opening ",path)
			} else {
				pic := id3v2.PictureFrame{
					Encoding:    id3v2.EncodingUTF8,
					MimeType:    "image/jpeg",
					//PictureType: id3v2.PTFrontCover,
					PictureType: 0x04,
					Description: "Front cover",
					Picture:     img,
				}
				mp3File.AddAttachedPicture(pic)
				if err = mp3File.Save(); err != nil {
					log.Println("Couldn't set the album art for ",path)
				} else {
					log.Println("Album art updated for ", path)
				}
				mp3File.Close()
			}
		}
	}
}

func getAlbumArt( getArtCh <- chan map[string]string, setArtCh chan <- map[string][]byte ){
	// get album art and send it to the setalbumart ch
	for image := range getArtCh {
		for path, url := range image {
			resp, err := http.Get(url)
			content, err := ioutil.ReadAll(resp.Body)
			if err == nil && len(content) > 0 {
				setArtCh <- map[string][]byte{path:content}
			}
			resp.Body.Close()
		}
	}
}

func getImageUrlFromMap( items []interface{} ) ( string, bool ){
	if bestItem , ok := items[0].(map[string]interface{}); ok {
		if album , ok := bestItem["album"].(map[string]interface{}); ok {
			for _, images := range album["images"].([]interface{}) {
				if img, ok := images.(map[string]interface{}); ok {
					if img["width"] == float64(640){
						return img["url"].(string), true
					}
				}
			}
		}
	}
	return "", false
}


func sendBackToSearc(paramsCh chan <- []string, resultStruct jobs.GetQueryResultStruct) {
	fmt.Println("Couldn't found anything with name => ",resultStruct.Params[3])
	if name, ok := jobs.TrimFromEnd(resultStruct.Params[3]); ok {
		resultStruct.Params[3] = name
		paramsCh <- resultStruct.Params
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
				bestMatches := jsonMap["best_match"].(map[string]interface{})
				items := bestMatches["items"].([]interface{})
				if len(items) > 0 {
					if url , ok := getImageUrlFromMap(items); ok {
						getArtCh <- map[string]string {path:url}
					} else {
						log.Println("Couldn't find any image for ", resultStruct.Params[3])
					}
				} else {
					sendBackToSearc(paramsCh,resultStruct)
				}
			}
		}
	}
}