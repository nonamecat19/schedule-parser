package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	geziyorClient "github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"os"
)

const url = "https://rozklad.ztu.edu.ua"

func parseGroups() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: quotesParse,
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()
}

func quotesParse(g *geziyor.Geziyor, r *geziyorClient.Response) {
	r.HTMLDoc.Find(".row.auto-clear").Each(func(i int, faculty *goquery.Selection) {
		var courses []map[string]any
		faculty.Find(".col.l2.s6.m4").Each(func(j int, course *goquery.Selection) {
			var groups = []map[string]any{}
			course.Find(".collection-item").Each(func(k int, group *goquery.Selection) {
				groups = append(groups, map[string]any{
					"group": group.Text(),
					"link":  group.AttrOr("href", ""),
				})
			})

			courses = append(courses, map[string]any{
				"name":   course.Find(".blue-text").Text(),
				"groups": groups,
			})
		})
		g.Exports <- map[string]any{
			"faculty": faculty.Find("h4").Text(),
			"courses": courses,
		}
	})
}

func main() {
	os.Remove("./out.json")
	parseGroups()
}
