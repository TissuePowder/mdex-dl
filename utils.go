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
	"sync"
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

	// fmt.Println(url)

	res, err := t.Client.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

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

		// fmt.Println(str)

		if strings.ContainsRune(str, '[') {
			arr := strings.Split(str, "[")
			if len(arr) > 2 {
				fmt.Println("Invalid chapter string format")
				os.Exit(1)
			}
			cpart = arr[0]
			// fmt.Println(cpart)
			if cpart == "" {
				cpart = "-"
			}
			ppart = arr[1]
			ppart = ppart[:len(ppart)-1]
		} else {
			cpart = str
		}

		// fmt.Println(cpart, ppart)

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

			if lb > ub {
				lb, ub = ub, lb
			}

			for _, c := range cAll[lb:ub] {
				if _, ok := pMap[c]; !ok {
					cList = append(cList, c)
				}
				pMap[c] = append(pMap[c], strings.Split(ppart, ",")...)
			}

		} else if _, ok := pMap[cpart]; !ok {
			cList = append(cList, cpart)
			pMap[cpart] = append(pMap[cpart], strings.Split(ppart, ",")...)
		} else {
			pMap[cpart] = append(pMap[cpart], strings.Split(ppart, ",")...)
		}
	}

	// fmt.Println(cList, pMap)

	return cList, pMap

}

type ScanGroup struct {
	Data struct {
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
	} `json:"data"`
}

func GetScanGroupName(id string, scanGroup *map[string]string, groups *[]string, wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()
	var d ScanGroup
	res, _ := http.Get(BaseUrl + "/group/" + id)
	json.NewDecoder(res.Body).Decode(&d)
	res.Body.Close()
	m.Lock()
	name := strings.Replace(d.Data.Attributes.Name, "/", "_", -1)
	(*scanGroup)[id] = name
	*groups = append(*groups, name)
	m.Unlock()
}

type TitleInfo struct {
	Data struct {
		Attributes struct {
			Title struct {
				En string `json:"en"`
			} `json:"title"`
		} `json:"attributes"`
	} `json:"data"`
}

func GetTitleName(id string, titleName *string, wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()
	var d TitleInfo
	res, _ := http.Get(BaseUrl + "/manga/" + id)
	json.NewDecoder(res.Body).Decode(&d)
	res.Body.Close()
	m.Lock()
	*titleName = strings.Replace(d.Data.Attributes.Title.En, "/", "_", -1)
	m.Unlock()
}
