package jobs

import (
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
)


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
	if err != nil {
		log.Fatal("Mayday, mayday. We got a problem!:", err)
	}
	// Read response
	fmt.Println(resp.StatusCode)
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

func GetMusicList(searchURL,accessToken,market,name string) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", searchURL ,nil); if err != nil {
		log.Fatal("Couldn't initialize new request: ", err)
	}

	q := req.URL.Query()
	q.Add("decorate_restrictions", "false")
	q.Add("include_external", "audio")
	q.Add("best_match", "true")
	q.Add("userless", "true")
	q.Add("market", market)
	q.Add("limit", "50")
	q.Add("type", "track,show_audio,episode_audio")
	q.Add("q", name)

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	req.Header.Set("Authorization", "Bearer " + accessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Mayday, mayday. We got a problem!:", err)
	}
	// Read json response data
	fmt.Println(resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))
}