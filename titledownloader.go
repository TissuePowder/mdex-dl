package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	v, _ := query.Values(t.Query.TitleQuery)
	fullUrl := fmt.Sprintf("%s?%s", t.Url, v.Encode())

	fmt.Println(fullUrl)

	res, _ := http.Get(fullUrl)

	var title Title

	json.NewDecoder(res.Body).Decode(&title)

	fmt.Printf("%+v\n", title)

	c := t.GetChapterList()

	fmt.Println(c)
}
