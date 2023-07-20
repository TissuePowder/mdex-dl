package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/go-querystring/query"
)

type Title struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     []struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Title              string `json:"title"`
			Volume             string `json:"volume"`
			Chapter            string `json:"chapter"`
			Pages              int    `json:"pages"`
			TranslatedLanguage string `json:"translatedLanguage"`
			Uploader           string `json:"uploader"`
			ExternalUrl        string `json:"externalUrl"`
			Version            int    `json:"version"`
			CreatedAt          string `json:"createdAt"`
			UpdatedAt          string `json:"updatedAt"`
			PublishAt          string `json:"publishAt"`
			ReadableAt         string `json:"readableAt"`
		} `json:"attributes"`
		Relationships []struct {
			Id      string `json:"id"`
			Type    string `json:"type"`
			Related string `json:"related"`
		} `json:"relationships"`
	} `json:"data"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func (t *TitleDownloader) StartDownloading() {
	// fmt.Println("title downloader")
	// fmt.Println(t.Url)
	// fmt.Println(t.Query)
	var titleName string
	var scanGroup map[string]string

	c, p := t.GetChapterList()

	// fmt.Println(c)
	// fmt.Println(p)

	t.Query.TitleQuery.Chapter = c

	for {

		v, _ := query.Values(t.Query.TitleQuery)
		fullUrl := fmt.Sprintf("%s?%s", t.Url, v.Encode())

		fmt.Println(fullUrl)

		res, _ := http.Get(fullUrl)

		var title Title

		json.NewDecoder(res.Body).Decode(&title)

		res.Body.Close()

		for _, d := range title.Data {

			var wg sync.WaitGroup

			for _, r := range d.Relationships {
				if r.Type == "manga" {
					if titleName == "" {
						wg.Add(1)
						go GetTitleName(r.Id, &titleName, &wg)
					}
				} else if r.Type == "scanlation_group" {
					if _, ok := scanGroup[r.Id]; !ok {
						wg.Add(1)
						go GetScanGroupName(r.Id, &scanGroup, &wg)
					}
				}
			}

			wg.Wait()



			t.Query.ChapterQuery.Pages = p[d.Attributes.Chapter]
			cDownloader := NewChapterDownloader(d.Id, t.Query)
			cDownloader.StartDownloading()
		}

		fmt.Printf("%+v\n", title)

		t.Query.TitleQuery.Offset += 100

		if t.Query.TitleQuery.Offset > title.Total {
			break
		}

	}

}
