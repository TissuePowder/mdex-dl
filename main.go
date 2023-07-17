package main

func main() {

	builder := NewQueryBuilder()

	builder.SetFlags()

	query := builder.Build()

	downloader := NewDownloader(query)

	downloader.StartDownloading()
}
