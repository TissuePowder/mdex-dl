package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

type Image struct {
	Url  string
	Path string
	Idx  int
}

func DownloadPage(image Image, client *http.Client) error {

	res, err := client.Get(image.Url)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	ext := filepath.Ext(image.Url)

	filename := fmt.Sprintf("%s_p%04d%s", image.Path, image.Idx+1, ext)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(filename), 0775)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)

	res.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(filename)
	return nil
}

func Worker(id int, jobs <-chan Image, results chan<- error, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		err := DownloadPage(job, client)
		results <- err
	}

}

func (c *ChapterDownloader) StartDownloading() {

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
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
	}

	var chapter Chapter

	var res *http.Response
	var err error

	for {
		res, err = client.Get(c.Url)
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(time.Duration(3) * time.Second)
		}
		if res.StatusCode == 200 {
			break
		} else if res.StatusCode == 429 {
			retryAfterStr := res.Header.Get("X-RateLimit-Retry-After")
			retryAfterUnix, _ := strconv.ParseInt(retryAfterStr, 10, 64)
			currentUnixTime := time.Now().Unix()
			durationSeconds := retryAfterUnix - currentUnixTime
			if durationSeconds > 0 {
				fmt.Printf("Hit rate limit. Waiting %ds before continuing.\n", durationSeconds)
				time.Sleep(time.Duration(durationSeconds) * time.Second)
			}
		} else {
			fmt.Println(res.StatusCode, res.Status, err)
		}
	}

	// fmt.Println(res.StatusCode)

	if err != nil {
		fmt.Println(err.Error())
	}

	json.NewDecoder(res.Body).Decode(&chapter)
	res.Body.Close()
	// fmt.Println(chapter)

	jobs := make(chan Image)
	results := make(chan error, len(chapter.Chapter.Data))

	for i := 1; i <= maxWorkers; i++ {
		go Worker(i, jobs, results, client, &wg)
	}

	var coll []string
	var ds string

	if c.Query.ChapterQuery.DataSaver {
		coll = chapter.Chapter.DataSaver
		ds = "data-saver"
	} else {
		coll = chapter.Chapter.Data
		ds = "data"
	}

	var pages []string

	for _, p := range c.Query.ChapterQuery.Pages {
		if p != "" {
			pages = append(pages, p)
		}
	}

	// fmt.Println(pages)

	if len(pages) > 0 {

		// fmt.Println("this is page array", pages)

		for _, p := range pages {

			if strings.Contains(p, "-") {
				arr := strings.Split(p, "-")
				var lb, ub int
				if arr[0] == "" {
					lb = 0
				} else {
					lb, _ = strconv.Atoi(arr[0])
					lb -= 1
				}

				if arr[1] == "" {
					ub = len(coll)
				} else {
					ub, _ = strconv.Atoi(arr[1])
					// ub += 1
					if ub > len(coll) {
						ub = len(coll)
					}
				}

				// fmt.Println(lb, ub)

				for i := lb; i < ub; i++ {
					img := coll[i]
					fullUrl := fmt.Sprintf("%s/%s/%s/%s", chapter.BaseUrl, ds, chapter.Chapter.Hash, img)
					// fmt.Println(fullUrl, i)
					jobs <- Image{fullUrl, c.Query.ChapterQuery.Path, i}
				}

			} else {
				i, _ := strconv.Atoi(p)
				if i > len(coll) {
					// fmt.Println("in continue block", i, p)
					// fmt.Println(chapter.Chapter.Data, c.Url)
					continue
				}
				// fmt.Println(p)
				img := coll[i-1]
				fullUrl := fmt.Sprintf("%s/%s/%s/%s", chapter.BaseUrl, ds, chapter.Chapter.Hash, img)
				jobs <- Image{fullUrl, c.Query.ChapterQuery.Path, i - 1}
			}

		}

	} else {

		for i, img := range coll {
			fullUrl := fmt.Sprintf("%s/%s/%s/%s", chapter.BaseUrl, ds, chapter.Chapter.Hash, img)
			jobs <- Image{fullUrl, c.Query.ChapterQuery.Path, i}
		}

	}

	close(jobs)
	wg.Wait()
	close(results)

}
