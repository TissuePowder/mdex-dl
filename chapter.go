package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (c *ChapterDownloader) StartDownloading() {

	maxWorkers := c.Query.Threads
	var wg sync.WaitGroup
	wg.Add(maxWorkers)

	var chapter Chapter

	var res *http.Response
	var err error

	for {
		res, err = c.Client.Get(c.Url)
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
		go Worker(i, jobs, results, c.Client, &wg)
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
		var pagesSerialized, pagesUniq []int

		for _, p := range pages {

			if strings.Contains(p, "-") {
				arr := strings.Split(p, "-")
				var lb, ub int
				if arr[0] == "" {
					lb = 1
				} else {
					lb, _ = strconv.Atoi(arr[0])
				}

				if arr[1] == "" {
					ub = len(coll)
				} else {
					ub, _ = strconv.Atoi(arr[1])
					if ub > len(coll) {
						ub = len(coll)
					}
				}

				// fmt.Println(lb, ub)

				for i := lb; i <= ub; i++ {
					pagesSerialized = append(pagesSerialized, i)
				}

			} else {
				i, _ := strconv.Atoi(p)
				if i > len(coll) {
					// fmt.Println("in continue block", i, p)
					// fmt.Println(chapter.Chapter.Data, c.Url)
					continue
				}
				pagesSerialized = append(pagesSerialized, i)
			}

		}

		mp := make(map[int]bool)
		for _, v := range pagesSerialized {
			if _, ok := mp[v]; !ok {
				mp[v] = true
				pagesUniq = append(pagesUniq, v)
			}
		}

		// fmt.Println(pagesSerialized, pagesUniq)

		for _, i := range pagesUniq {
			img := coll[i-1]
			fullUrl := fmt.Sprintf("%s/%s/%s/%s", chapter.BaseUrl, ds, chapter.Chapter.Hash, img)
			jobs <- Image{fullUrl, c.Query.ChapterQuery.Path, i - 1}
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
