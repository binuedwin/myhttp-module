package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func main() {
	// read parallel flag if present
	parallel := flag.Int("parallel", 10, "allows requests to be perfomed in parallel, cannot be set > 10")

	flag.Parse()

	if *parallel > 10 {
		flag.PrintDefaults()
		return
	}

	// get all other non flags args; urls
	urls := flag.Args()

	// schedule worker goroutines and aggregate result
	results := scheduleWorkers(urls, *parallel)

	// print results to stdout
	fmt.Println(strings.Join(results, "\n"))
}

// scheduleWorkers allows to schedule worker goroutines provided the worker count and the
// request urls.
func scheduleWorkers(urls []string, parallelCount int) []string {
	// create buffered channel for the request channel and fill with request urls
	// this allows for non-blocking reads
	reqChan := make(chan string, len(urls))
	for _, url := range urls {
		reqChan <- normalizeToURL(url)
	}

	// since no other values will ever be sent on reqChan, close it.
	close(reqChan)

	// create buffered response channel, this allows non-blocking writes
	// and capacity equal to the number of requests
	respChan := make(chan string, len(urls))

	// create WaitGroup for workers
	var wg sync.WaitGroup
	for i := 0; i < parallelCount; i++ {
		wg.Add(1)
		go worker(&wg, getBodyContent, reqChan, respChan)
	}

	// wait for parallel workers to complete execution
	wg.Wait()

	// No other data will ever be sent on the channel
	close(respChan)

	var respPairs []string
	for {
		resp, more := <-respChan
		if !more {
			// if no more data from channel, break read loop
			break
		}
		respPairs = append(respPairs, resp)
	}

	return respPairs
}

// worker listens to the request channel and generates the MD5 hash string
// and sends it to the response channel.
func worker(wg *sync.WaitGroup, bodyContent func(string) []byte, reqChan <-chan string, respChan chan<- string) {
	defer wg.Done()
	for {
		reqURL, more := <-reqChan
		if !more {
			return
		}

		respChan <- fmt.Sprintf("%s %s", reqURL, generateMD5HashString(bodyContent(reqURL)))
	}
}

// getBodyContent performs a HTTP GET request on the url provided
// returning the body content.
func getBodyContent(reqURL string) []byte {
	// send request and retrieve response
	resp, err := http.Get(reqURL)
	if err != nil {
		fmt.Println(err)
	}

	// client must close the response body when finished with it
	defer resp.Body.Close()

	// create buffer with content
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// return content buffer
	return buff
}

// generateMD5HashString generates the MD5 hash string of the given content.
func generateMD5HashString(content []byte) string {
	// generate the MD5 hash
	hash := md5.Sum(content)
	// convert the hash to a string
	return hex.EncodeToString(hash[:])
}

// mormalizeToURL helps normalite the host names by making sure they have HTTP
// protocol declared at the prefix.
func normalizeToURL(urlStr string) string {
	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		return urlStr
	}
	return "http://" + urlStr
}
