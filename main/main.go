package main

import (
	"battletester/netCaller"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Results struct {
	Url          string
	StatusCode   int
	ResponseTime time.Duration
}

type UrlDetail struct {
	StatusCode   map[int]int
	ResponseTime float64
	Count        int
}

func (detail *UrlDetail) AddItem(item float64) {
	detail.ResponseTime += item
}

func startWorkers(ch <-chan string, rc chan<- Results) {
	url := <-ch
	urlRespC := map[int]netCaller.UrlResponseCode{}
	_, err := netCaller.UrlCall(url, urlRespC)
	if err != nil {
		log.Printf("%s \n", err)
		res := Results{Url: url, StatusCode: 700, ResponseTime: 0}
		rc <- res
		return
	}
	var finalres Results
	for stcode, res := range urlRespC {
		for k, v := range res.Response {
			finalres = Results{Url: k, StatusCode: stcode, ResponseTime: v.ResponseTime}
		}
	}
	rc <- finalres

}

func createWorkers(ch chan string, rc chan Results, workers int) {
	for i := 0; i < workers; i++ {
		go startWorkers(ch, rc)
	}
}

func pushTo(url string, ch chan<- string, workers int) {
	for i := 0; i < workers; i++ {
		ch <- url
	}
}
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	file_url := flag.String("f", "", "select the host to attack")
	workers := flag.Int("w", 10, "select numberof workers to attack with")
	flag.Parse()
	if *file_url == "" {
		fmt.Println("Host url missing")
		return
	}
	lines, err := readLines(*file_url)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	lineSlice := []string{}
	for _, line := range lines {
		lineSlice = append(lineSlice, line)
	}

	urlDict := map[string]*UrlDetail{}

	for _, urlSt := range lineSlice {
		fmt.Println(urlSt)
		ch := make(chan string)
		rc := make(chan Results)
		createWorkers(ch, rc, *workers)
		pushTo(urlSt, ch, *workers)
		for j := 0; j < *workers; j++ {
			result := <-rc
			_, ok := urlDict[result.Url]
			if !ok {
				child := UrlDetail{StatusCode: map[int]int{result.StatusCode: 1}, Count: 1, ResponseTime: result.ResponseTime.Seconds()}
				urlDict[result.Url] = &child
			} else {
				urlDict[result.Url].StatusCode[result.StatusCode]++
				urlDict[result.Url].ResponseTime += result.ResponseTime.Seconds()
				urlDict[result.Url].Count++

			}

		}
	}

	for k, values := range urlDict {
		fmt.Println(k)
		fmt.Printf("\tTotal Request:%d \n\tTotal Response:%f  \n\tAverage Response Time:%f\n", values.Count, values.ResponseTime, values.ResponseTime/float64(values.Count))
		for key, ele := range values.StatusCode {
			fmt.Printf("\t\t%d : %d\n", key, ele)
		}
		fmt.Println()

	}

}
