package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

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

func processTable(tableObject *goquery.Selection, g *Group) {
	fmt.Println("Processing table and generating payscale")

	// Generate payscale name and map[datetime][]int
	payRates := make(map[string][]int)

	tableObject.Each(func(i int, table *goquery.Selection) {

		rawCaption := strings.TrimSpace(table.Find("caption").Text())

		// different groups format their captions differently.
		// figure out which separator they use ":" or " - " and split on that
		captionArray := []string{}

		if strings.Contains(rawCaption, ":") {
			captionArray = strings.Split(rawCaption, ":")
		} else {
			captionArray = strings.Split(rawCaption, " - ")
		}

		// Isn't empty
		if captionArray[0] != "" &&
			// Is under 12 characters
			len(captionArray[0]) <= 12 &&
			// is at least 3 characters
			len(captionArray[0]) > 2 &&
			// refers to annual pay
			strings.Contains(strings.ToLower(captionArray[1]), "annual") &&
			// contains the identifer we are looking for
			strings.Contains(strings.ToLower(captionArray[0]), g.Identifier) {

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

						if strings.Contains(td.Text(), "to") {
							payRange := strings.Split(td.Text(), " to ")
							pay1, _ := strconv.Atoi(strings.TrimSpace(payRange[0]))
							pay2, _ := strconv.Atoi(strings.TrimSpace(payRange[1]))
							payRates[date] = append(payRates[date], pay1, pay2)
						} else {
							pay := strings.Replace(td.Text(), ",", "", -1)

							payAsNum, err := strconv.Atoi(pay)
							if err != nil {
								payAsNum = 0
							}

							payRates[date] = append(payRates[date], payAsNum)
						}

					})
				}
				p.PayScale = payRates
				p.Steps = len(payRates[date])
			})
			g.PayScales = append(g.PayScales, p)
		}
	})
}
