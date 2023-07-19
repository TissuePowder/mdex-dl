package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Order struct {
	CreatedAt  string `url:"createdAt,omitempty"`
	UpdatedAt  string `url:"updatedAt,omitempty"`
	PublishAt  string `url:"publishAt,omitempty"`
	ReadableAt string `url:"readableAt,omitempty"`
	Volume     string `url:"volume,omitempty"`
	Chapter    string `url:"chapter,omitempty"`
}

type TitleQuery struct {
	Limit              int      `url:"limit"`
	Offset             int      `url:"offset"`
	Ids                []string `url:"ids[],omitempty"`
	Title              string   `url:"title,omitempty"`
	Groups             []string `url:"groups[],omitempty"`
	Uploader           []string `url:"uploader[],omitempty"`
	Manga              string   `url:"manga,omitempty"`
	Volume             []string `url:"volume[],omitempty"`
	Chapter            []string `url:"chapter[],omitempty"`
	TranslatedLanguage []string `url:"translatedLanguage[],omitempty"`
	ContentRating      []string `url:"contentRating[]"`
	ExcludedGroups     []string `url:"excludedGroups[],omitempty"`
	ExcludedUploaders  []string `url:"excludedUploaders[],omitempty"`
	CreatedAtSince     string   `url:"createdAtSince,omitempty"`
	UpdatedAtSince     string   `url:"updatedAtSince,omitempty"`
	PublishAtSince     string   `url:"publishAtSince,omitempty"`
	Order              Order    `url:"order,omitempty"`
	Includes           []string `url:"includes[],omitempty"`
}

type Query struct {
	DataSaver  bool
	Threads    int
	Url        string
	TitleQuery TitleQuery
}

type QueryBuilder struct {
	Query Query
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		Query: Query{},
	}
}

func (q *QueryBuilder) Build() Query {
	return q.Query
}

func SplitChapterString(s string) ([]string, error) {
	var arr []string
	var builder strings.Builder
	b := 0

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch c {
		case ',':
			if b == 0 && builder.String() != "" {
				arr = append(arr, builder.String())
				builder.Reset()
			} else {
				builder.WriteByte(c)
			}
		case '[':
			if b > 0 {
				return nil, fmt.Errorf("not a valid filter string")
			}
			b++
			builder.WriteByte(c)

		case ']':
			if b < 0 {
				return nil, fmt.Errorf("invalid filter string")
			}
			b--
			builder.WriteByte(c)

		default:
			builder.WriteByte(c)
		}
	}

	if b != 0 {
		return nil, fmt.Errorf("invalid filter string")
	}

	arr = append(arr, builder.String())
	return arr, nil
}

func (q *QueryBuilder) SetFlags() *QueryBuilder {

	title := flag.String("title", "", "")
	group := flag.String("group", "", "")
	uploader := flag.String("uploader", "", "")
	volume := flag.String("volume", "", "")
	chapter := flag.String("chapter", "", "")
	lang := flag.String("lang", "", "")
	egroup := flag.String("excluded-group", "", "")
	euploader := flag.String("excluded-uploader", "", "")
	created := flag.String("created-since", "", "")
	updated := flag.String("updated-since", "", "")
	published := flag.String("published-since", "", "")

	datasaver := flag.Bool("data-saver", false, "")
	thread := flag.Int("thread", 5, "")

	flag.Parse()

	if *title != "" {
		q.Query.TitleQuery.Title = *title
	}

	if *group != "" {
		q.Query.TitleQuery.Groups = strings.Split(*group, ",")
	}

	if *uploader != "" {
		q.Query.TitleQuery.Uploader = strings.Split(*uploader, ",")
	}

	if *volume != "" {
		q.Query.TitleQuery.Volume = strings.Split(*volume, ",")
	}

	if *chapter != "" {
		if cList, err := SplitChapterString(*chapter); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			q.Query.TitleQuery.Chapter = cList
			fmt.Println(cList)
		}

	}

	if *lang != "" {
		q.Query.TitleQuery.TranslatedLanguage = strings.Split(*lang, ",")
	}

	if *egroup != "" {
		q.Query.TitleQuery.ExcludedGroups = strings.Split(*egroup, ",")
	}

	if *euploader != "" {
		q.Query.TitleQuery.ExcludedUploaders = strings.Split(*euploader, ",")
	}

	if *created != "" {
		q.Query.TitleQuery.CreatedAtSince = *created
	}

	if *updated != "" {
		q.Query.TitleQuery.UpdatedAtSince = *updated
	}

	if *published != "" {
		q.Query.TitleQuery.PublishAtSince = *published
	}

	q.Query.TitleQuery.Limit = 100
	q.Query.TitleQuery.Offset = 0
	q.Query.TitleQuery.Order = Order{
		Chapter: "asc",
	}

	q.Query.DataSaver = *datasaver
	q.Query.Threads = *thread

	url := flag.Args()

	// fmt.Println(url)

	if len(url) != 1 {
		fmt.Println("Pleas use proper format: mdex-dl [options] url")
		os.Exit(1)
	}

	q.Query.Url = url[0]

	return q

}
