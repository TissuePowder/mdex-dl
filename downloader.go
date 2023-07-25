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
	Client *MdClient
	Query Query
}

type ChapterDownloader struct {
	Url   string
	Client *MdClient
	Query Query
}

func NewDownloader(client *MdClient, query Query) Downloader {
	url := strings.TrimPrefix(query.Url, "https://")
	arr := strings.Split(url, "/")
	if arr[1] == "chapter" {
		query.TitleQuery.Ids = []string{arr[2]}
		return NewTitleDownloader("", client, query)
	}
	return NewTitleDownloader(arr[2], client, query)
}

func NewTitleDownloader(id string, client *MdClient, query Query) Downloader {
	query.TitleQuery.Manga = id
	return &TitleDownloader{
		Url:   BaseUrl + "/chapter",
		Client: client,
		Query: query,
	}
}

func NewChapterDownloader(id string, client *MdClient, query Query) Downloader {
	return &ChapterDownloader{
		Url:   BaseUrl + "/at-home/server/" + id,
		Client: client,
		Query: query,
	}
}
