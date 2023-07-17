package main

import "fmt"

func (t *TitleDownloader) StartDownloading() {
	fmt.Println("title downloader")
	fmt.Println(t.Url)
	fmt.Println(t.Query)
}
