package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

func TraverseChapters(cMap map[string]interface{}, cAll *[]string) {
	for _, val := range cMap {
		if c, ok := val.(map[string]interface{}); ok {
			TraverseChapters(c, cAll)
			if nc, ok := c["chapter"].(string); ok {
				*cAll = append(*cAll, nc)
			}
		}
	}
}

func (t *TitleDownloader) GetChapterList() ([]string, map[string][]string) {

	var cList []string
	pMap := make(map[string][]string)

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

	var cAll []string

	for _, val := range data {
		if cMap, ok := val.(map[string]interface{}); ok {
			TraverseChapters(cMap, &cAll)
		}
	}

	sort.Slice(cAll, func(i, j int) bool {
		n1, _ := strconv.ParseFloat(cAll[i], 64)
		n2, _ := strconv.ParseFloat(cAll[j], 64)
		return n1 < n2
	})

	for _, str := range t.Query.TitleQuery.Chapter {

		var cpart, ppart string
		var cl, cr float64

		if strings.ContainsRune(str, '[') {
			arr := strings.Split(str, "[")
			if len(arr) > 2 {
				fmt.Println("error")
				os.Exit(1)
			}
			cpart = arr[0]
			fmt.Println(cpart)
			if cpart == "" {
				cpart = "-"
			}
			ppart = arr[1]
			ppart = ppart[:len(ppart)-1]
		} else {
			cpart = str
		}

		if strings.ContainsRune(cpart, '-') {
			arr := strings.Split(cpart, "-")
			if arr[0] == "" {
				cl = float64(0)
			} else {
				cl, _ = strconv.ParseFloat(arr[0], 64)
			}

			if arr[1] == "" {
				cr, _ = strconv.ParseFloat(cAll[len(cAll)-1], 64)
			} else {
				cr, _ = strconv.ParseFloat(arr[1], 64)
			}

			lb := sort.Search(len(cAll), func(i int) bool {
				n, _ := strconv.ParseFloat(cAll[i], 64)
				return n >= cl
			})

			ub := sort.Search(len(cAll), func(i int) bool {
				n, _ := strconv.ParseFloat(cAll[i], 64)
				return n > cr
			})

			for _, c := range cAll[lb:ub] {
				if _, ok := pMap[c]; !ok {
					cList = append(cList, c)
				}
				pMap[c] = append(pMap[c], ppart)
			}

		} else if _, ok := pMap[cpart]; !ok {
			cList = append(cList, cpart)
		}
	}

	return cList, pMap

}
