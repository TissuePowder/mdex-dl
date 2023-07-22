package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

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

	var titleName string
	scanGroup := make(map[string]string)

	c, p := t.GetChapterList()

	t.Query.TitleQuery.Chapter = c

	for {

		v, _ := query.Values(t.Query.TitleQuery)
		fullUrl := fmt.Sprintf("%s?%s", t.Url, v.Encode())

		// fmt.Println(fullUrl)

		res, _ := t.Client.Get(fullUrl)

		var title Title

		json.NewDecoder(res.Body).Decode(&title)

		res.Body.Close()

		for _, d := range title.Data {

			var wg sync.WaitGroup
			var m sync.Mutex

			var groups []string

			for _, r := range d.Relationships {
				if r.Type == "manga" {
					if titleName == "" {
						wg.Add(1)
						go GetTitleName(r.Id, &titleName, &wg, &m)
					}
				} else if r.Type == "scanlation_group" {
					if _, ok := scanGroup[r.Id]; !ok {
						wg.Add(1)
						go GetScanGroupName(r.Id, &scanGroup, &groups, &wg, &m)
					} else {
						groups = append(groups, scanGroup[r.Id])
					}
				}
			}

			wg.Wait()

			gstring := strings.Join(groups, " + ")

			d.Attributes.Chapter = strings.Replace(d.Attributes.Chapter, ",", ".", -1)

			carr := strings.Split(d.Attributes.Chapter, ".")

			cstr := fmt.Sprintf("%04s", carr[0])

			if len(carr) > 1 {
				cstr = cstr + "." + carr[1]
			}

			if gstring == "" {
				gstring = "No Group"
			}

			var path string
			if t.Query.NoChapterDir {
				path = fmt.Sprintf("%s/c%s_[%s]", titleName, cstr, gstring)
			} else {
				path = fmt.Sprintf("%s/Ch.%s [%s]/c%s", titleName, cstr, gstring, cstr)
			}
			// fmt.Println(path)

			t.Query.ChapterQuery.Path = path
			t.Query.ChapterQuery.Pages = p[d.Attributes.Chapter]
			cDownloader := NewChapterDownloader(d.Id, t.Client, t.Query)
			t := time.Now()
			cDownloader.StartDownloading()
			fmt.Printf("chapter %s download time: %s\n", d.Attributes.Chapter, time.Since(t).String())
		}

		// fmt.Printf("%+v\n", title)

		t.Query.TitleQuery.Offset += 100

		if t.Query.TitleQuery.Offset > title.Total {
			break
		}

	}

}
