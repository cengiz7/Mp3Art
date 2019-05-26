package jobs

import (
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"net/url"
	"time"
)

type GetQueryResultStruct struct {
	Params []string
	Result map[string][]byte
}

func GetAccessToken( baseURL string ) string {
	return findAccessToken( makeTokenRequest( baseURL ) )
}

func makeTokenRequest( baseURL string ) *http.Response {
	client := &http.Client{}

	req, err := http.NewRequest("GET", baseURL,nil); if err != nil {
		log.Fatal("Couldn't initialize new request: ", err)
	}
	// set the necessary headers
	req.Header.Set("Host", "open.spotify.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36")
	req.Header.Set("DNT", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")
	// perform request
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		log.Fatal("Mayday, mayday. We got a problem with getting access token!:", err)
	}
	// Read response
	//data, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))
	return resp
}

func findAccessToken( response *http.Response ) string {
	token := ""
	for _, cookie := range response.Cookies() {
		if cookie.Name == "wp_access_token" {
			token = cookie.Value
			fmt.Println("Access token has successfully taken: ",token)
			break
		}
	}
	if token == "" {
		log.Fatal("Access token unreachable. \nExitting...")
	}
	return token
}

func GetQueryResult( paramsCh <- chan []string , getQueryResultCh chan <- GetQueryResultStruct ) {
	client := &http.Client{}
	for paramList := range paramsCh {
		accessToken	:= paramList[0]
		searchURL 	:= paramList[1]
		market		:= paramList[2]
		name		:= paramList[3]
		limit		:= paramList[4]
		filePath    := paramList[5]

		req, err := http.NewRequest("GET", searchURL ,nil); if err != nil {
			log.Fatal("Couldn't initialize new request: ", err)
		}
		// set params
		q := setQueryParams( req, market, limit, name)
		req.URL.RawQuery = q.Encode()
		req.Header.Set("Authorization", "Bearer " + accessToken)
		// fmt.Println("\n"+req.URL.String()+"\n")

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			log.Println(resp.StatusCode)
			log.Fatal("Mayday, mayday. We got a problem with making search!:", err)
		} else {
			// Read json response data
			data, err := ioutil.ReadAll(resp.Body); if err != nil {
				log.Println("Error while getting search query result: ",err)
			} else {
				// initialize map and send to the queue
				m 		:= map[string][]byte{filePath: data}
				strct 	:= GetQueryResultStruct{paramList, m}
				getQueryResultCh <- strct
			}
		}
		// wait half second for the nex request
		time.Sleep(time.Second / 2)
	}
}

func setQueryParams(req *http.Request, market, limit, name string) url.Values {
	// "api.spotify.com/v1/search?type=album%2Cartist%2Cplaylist%2Ctrack%2Cshow_audio%2Cepisode_audio&q=dualipa*&decorate_restrictions=false&best_match=true&include_external=audio&limit=50&userless=true&market=TR"
	q := req.URL.Query()
	q.Add("decorate_restrictions", "false")
	q.Add("include_external", "audio")
	q.Add("best_match", "true")
	q.Add("userless", "true")
	q.Add("market", market)
	q.Add("limit", limit)
	q.Add("type", "track,album")
	q.Add("q", name)
	return q
}