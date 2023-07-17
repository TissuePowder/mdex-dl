package main

import (
	"strings"
)

const (
	BaseUrl = "https://api.mangadex.org"
)

type Downloader interface {
	StartDownloading()
}

type TitleDownloader struct {
	Url   string
	Query Query
}

type ChapterDownloader struct {
	Url   string
	Query Query
}

func NewDownloader(query Query) Downloader {
	url := strings.TrimPrefix(query.Url, "https://")
	arr := strings.Split(url, "/")
	if arr[1] == "chapter" {
		return NewChapterDownloader(arr[2], query)
	} else {
		return NewTitleDownloader(arr[2], query)
	}
}

func NewTitleDownloader(id string, query Query) Downloader {
	return &TitleDownloader{
		Url:   BaseUrl + "/chapter",
		Query: query,
	}
}

func NewChapterDownloader(id string, query Query) Downloader {
	return &ChapterDownloader{
		Url:   BaseUrl + "/at-home/server/" + id,
		Query: query,
	}
}
