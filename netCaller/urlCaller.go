package netCaller

import (
	//"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type UrlResponseCode struct {
	Response map[string]UrlResponse
}
type UrlResponse struct {
	Status       string
	ResponseTime time.Duration
}

func UrlCall(url string, urlRespC map[int]UrlResponseCode) ([]byte, error) {
	start := time.Now()
	var netClient = &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		//fmt.Printf("Error:in url call %v \n", err)
		return nil, err
	}
	total := time.Since(start)
	//fmt.Printf("Status: %s: %s: %v Time taken : %v \n", url, resp.Status, resp.StatusCode, total.Seconds())
	respStruct := UrlResponse{Status: resp.Status, ResponseTime: total}
	_, ok := urlRespC[resp.StatusCode]
	if !ok {
		response := UrlResponseCode{Response: map[string]UrlResponse{}}
		urlRespC[resp.StatusCode] = response
	}
	urlRespC[resp.StatusCode].Response[url] = respStruct
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
