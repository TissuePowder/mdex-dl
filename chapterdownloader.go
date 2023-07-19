package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Chapter struct {
	Result  string `json:"result"`
	BaseUrl string `json:"baseUrl"`
	Chapter struct {
		Hash      string   `json:"hash"`
		Data      []string `json:"data"`
		DataSaver []string `json:"dataSaver"`
	} `json:"chapter"`
}

func DownloadPage(url string, client *http.Client) error {

	res, err := client.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer res.Body.Close()

	arr := strings.Split(url, "/")
	filename := arr[len(arr)-1]

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(filename)
	return nil
}

func Worker(id int, jobs <-chan string, results chan<- error, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range jobs {
		err := DownloadPage(url, client)
		results <- err
		// fmt.Println(id, url)
	}

}

func (c *ChapterDownloader) StartDownloading() {
	start := time.Now()
	fmt.Println("chapter downloader")
	fmt.Println(c.Url)
	fmt.Println(c.Query)

	var chapter Chapter
	res, _ := http.Get(c.Url)
	json.NewDecoder(res.Body).Decode(&chapter)
	res.Body.Close()
	// fmt.Println(chapter)

	jobs := make(chan string)
	results := make(chan error, len(chapter.Chapter.Data))

	maxWorkers := c.Query.Threads
	var wg sync.WaitGroup
	wg.Add(maxWorkers)

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxConnsPerHost:       maxWorkers,
		MaxIdleConnsPerHost:   maxWorkers,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	for i := 1; i <= maxWorkers; i++ {
		go Worker(i, jobs, results, client, &wg)
	}

	for _, file := range chapter.Chapter.Data {
		fullUrl := fmt.Sprintf("%s/data/%s/%s", chapter.BaseUrl, chapter.Chapter.Hash, file)
		jobs <- fullUrl
	}

	close(jobs)
	wg.Wait()
	close(results)

	fmt.Println(time.Since(start))

}
