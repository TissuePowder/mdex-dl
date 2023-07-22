package main

import (
	"fmt"
	"time"
)

func main() {

	start := time.Now()

	builder := NewQueryBuilder()

	builder.SetFlags()

	query := builder.Build()

	client := NewMdClient()

	downloader := NewDownloader(client, query)

	downloader.StartDownloading()

	fmt.Println("Total time:", time.Since(start))
}
