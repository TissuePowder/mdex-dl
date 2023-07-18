package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

func TraverseChapters(cMap map[string]interface{}, cList *[]string) {
	for _, val := range cMap {
		if c, ok := val.(map[string]interface{}); ok {
			TraverseChapters(c, cList)
			if nc, ok := c["chapter"].(string); ok {
				*cList = append(*cList, nc)
			}
		}
	}
}

func (t *TitleDownloader) GetChapterList() []string {

	params := url.Values{}

	for _, value := range t.Query.TitleQuery.TranslatedLanguage {
		params.Add("translatedLanguage[]", value)
	}

	for _, value := range t.Query.TitleQuery.Groups {
		params.Add("groups[]", value)
	}

	url := fmt.Sprintf("%s/manga/%s/aggregate?%s", BaseUrl, t.Query.TitleQuery.Manga, params.Encode())

	fmt.Println(url)

	res, _ := http.Get(url)

	var data map[string]interface{}

	json.NewDecoder(res.Body).Decode(&data)

	// p, _ := json.MarshalIndent(data, "", " ")

	// fmt.Println(string(p))

	var cList []string

	for _, val := range data {
		if cMap, ok := val.(map[string]interface{}); ok {
			TraverseChapters(cMap, &cList)
		}
	}

	sort.Slice(cList, func(i, j int) bool {
		n1, _ := strconv.ParseFloat(cList[i], 64)
		n2, _ := strconv.ParseFloat(cList[j], 64)
		return n1 < n2
	})

	// fmt.Println(cList)

	return cList

}
