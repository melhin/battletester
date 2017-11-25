package netCaller

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseParent(body []byte, urlList *map[string]interface{}) {
	var dat map[string]interface{}
	err := json.Unmarshal(body, &dat)
	if err != nil {
		fmt.Println(err)
	}
	resp_data := dat["data"].(map[string]interface{})
	*urlList = resp_data["endpoints"].(map[string]interface{})

}

func ProcessUrlList(urlList map[string]interface{}, urlResponseC map[int]UrlResponseCode, failed_url *[]string) {
	for k, v := range urlList {
		prfx := strings.HasPrefix(k, "get_")
		if prfx {
			_, err := UrlCall(v.(string), urlResponseC)
			if err != nil {
				fmt.Println("Error: Failed %v with error %v", k, err)
				*failed_url = append(*failed_url, v.(string))
			}
		}
	}

}
