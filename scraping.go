package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// GetGroupURLs scrapes the main TBS collective agreements page and returns specific agreement URLs.
func GetGroupURLs(url string) []string {

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
	path := url

	// set empty array for urls
	urls := []string{}

	// Test scraping function rates of pay
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		fmt.Println(link)

		if strings.Contains(link, "rates") {
			urls = append(urls, link)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(path)

	return urls
}

// GetPayScales parses a specific TBS Website looking for pay tables.
func GetPayScales(groupURL string, g *Group) {

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

	// Test scraping function rates of pay
	c.OnHTML("body", func(e *colly.HTMLElement) {

		goquerySelection := e.DOM

		g.Name = strings.TrimSpace(goquerySelection.Find("h1").Text())
		g.PayScales = []PayScale{}

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
}

func processTable(tableObject *goquery.Selection, g *Group) {
	fmt.Println("Processing table and generating payscale")

	tableObject.Each(func(i int, table *goquery.Selection) {

		rawCaption := strings.TrimSpace(table.Find("caption").Text())

		// different groups format their captions differently.
		// figure out which separator they use ":" or " - " and split on that
		captionArray := []string{}
		caption2 := ""

		if strings.Contains(rawCaption, ":") {
			// caption is split by :
			captionArray = strings.Split(rawCaption, ":")
		} else if strings.Contains(rawCaption, " - ") {
			// caption is split by " - "
			captionArray = strings.Split(rawCaption, " - ")
		} else {
			// caption isn't split
			captionArray = append(captionArray, rawCaption)
		}

		if len(captionArray) > 1 {
			caption2 = captionArray[1]
		}

		// extract level from name
		level := 0

		levelString := strings.Split(captionArray[0], "-")
		if len(levelString) > 1 {
			level, _ = strconv.Atoi(levelString[1])
		} else {
			level = 0
		}

		// Isn't empty
		if captionArray[0] != "" &&
			// Is under 12 characters
			len(captionArray[0]) <= 12 &&
			// is at least 3 characters
			len(captionArray[0]) > 2 &&
			// refers to annual pay
			(strings.Contains(strings.ToLower(caption2), "annual") ||
				caption2 == "") &&
			// contains the identifer we are looking for
			strings.Contains(strings.ToUpper(captionArray[0]), g.Identifier) &&
			// not already duplicated in Payscales
			!g.existsInPayScaleNames(captionArray[0]) &&
			// able to extract numerical level from caption
			level != 0 {

			// Create PayScale struct
			p := PayScale{
				Name:  captionArray[0],
				Level: level,
			}

			// Find table body
			tb := table.Find("tbody")

			// Find each table row
			tb.Find("tr").Each(func(rowIndex int, tr *goquery.Selection) {

				// initialize row for salary
				row := []int{}

				date := "1980-01-01T11:45:26.371Z"
				dateString := ""

				// find the datetime value for the row
				tr.Find("time").Each(func(indexOfTd int, th *goquery.Selection) {
					dateString, _ = th.Attr("datetime")
					date = dateString + "T11:45:26.371Z"
				})

				// Check if data is valid and, if so, scan rows
				if date != "1980-01-01T11:45:26.371Z" {

					// find each salary cell
					tr.Find("td").Each(func(indexOfTd int, td *goquery.Selection) {

						if strings.Contains(td.Text(), "to") {
							payRange := strings.Split(td.Text(), " to ")
							pay1, _ := strconv.Atoi(strings.TrimSpace(payRange[0]))
							pay2, _ := strconv.Atoi(strings.TrimSpace(payRange[1]))
							row = append(row, pay1, pay2)
						} else {
							pay := strings.Replace(td.Text(), ",", "", -1)

							payAsNum, err := strconv.Atoi(pay)
							if err != nil {
								payAsNum = 0
							}

							row = append(row, payAsNum)
						}

					})
					// if row totals zero or has no steps, disregard
				}

				// Check that there are non 0 values here
				if len(row) > 0 &&
					sum(row) > 0 {

					inc := RateOfPay{
						DateTime: dateString,
						Salary:   row,
					}

					p.RatesOfPay = append(p.RatesOfPay, inc)
					p.Steps = len(inc.Salary)
				}

			})
			// table row iterator complete, add payscale to group payscales
			g.PayScales = append(g.PayScales, p)
		}
	})
}
