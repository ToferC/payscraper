package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var groupURLs = []string{
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=4#rates-ec",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=1#rates-cs",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=10#rates-fb",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=11#rates-fi",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=15#rates-as",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=15#rates-pm",
}

type Group struct {
	Name       string
	Identifier string
	URL        string
	PayScales  []PayScale
}

type PayScale struct {
	Name     string
	PayScale map[string][]int
}

func main() {

	groups := []Group{}

	for _, url := range groupURLs {

		g := Group{
			Identifier: url[len(url)-2:],
			URL:        url,
		}

		GetPayScales(url, &g)

		g.save()

		groups = append(groups, g)
	}
	// fmt.Println(groups) or do something else with them
}

func processTable(tableObject *goquery.Selection, g *Group) {
	fmt.Println("Processing table and generating payscale")

	// Generate payscale name and map[datetime][]int
	payRates := make(map[string][]int)

	tableObject.Each(func(i int, table *goquery.Selection) {

		rawCaption := strings.TrimSpace(table.Find("caption").Text())

		captionArray := strings.Split(rawCaption, " - ")

		if captionArray[0] != "" && len(captionArray[0]) <= 6 {

			p := PayScale{
				Name: captionArray[0],
			}

			tb := table.Find("tbody")

			tb.Find("tr").Each(func(rowIndex int, tr *goquery.Selection) {

				date := "2020-01-01"

				tr.Find("time").Each(func(indexOfTd int, th *goquery.Selection) {
					date, _ = th.Attr("datetime")
					payRates[date] = []int{}

				})

				if date != "2020-01-01" {
					tr.Find("td").Each(func(indexOfTd int, td *goquery.Selection) {

						pay := strings.Replace(td.Text(), ",", "", -1)

						payAsNum, err := strconv.Atoi(pay)
						if err != nil {
							payAsNum = 0
						}

						payRates[date] = append(payRates[date], payAsNum)

					})
				}
				p.PayScale = payRates
			})
			g.PayScales = append(g.PayScales, p)
		}
	})
}

func GetPayScales(groupURL string, g *Group) map[string][]string {

	// Initialize Colly Collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.tbs-sct.gc.ca"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	// set URLs for scraping

	path := groupURL

	payRates := make(map[string][]string)

	// Test scraping function rates of pay
	c.OnHTML("body", func(e *colly.HTMLElement) {

		goquerySelection := e.DOM

		g.Name = strings.TrimSpace(goquerySelection.Find("h1").Text())

		goquerySelection.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
			if index == 0 {
			} else {
				fmt.Println("Found Pay Table", index)
				processTable(tablehtml, g)
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(path)

	g.URL = groupURL

	fmt.Println(g)

	return payRates

}
