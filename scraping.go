package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

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

	g.save()

	fmt.Println(g)
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
			strings.Contains(strings.ToLower(captionArray[0]), g.Identifier) {

			p := PayScale{
				Name: captionArray[0],
			}

			tb := table.Find("tbody")

			tb.Find("tr").Each(func(rowIndex int, tr *goquery.Selection) {

				inc := Increment{}

				date := "1980-01-01T11:45:26.371Z"

				tr.Find("time").Each(func(indexOfTd int, th *goquery.Selection) {
					dateString, _ := th.Attr("datetime")
					date = dateString + "T11:45:26.371Z"
					inc.DateTime = dateString
				})

				if date != "1980-01-01T11:45:26.371Z" {

					tr.Find("td").Each(func(indexOfTd int, td *goquery.Selection) {

						if strings.Contains(td.Text(), "to") {
							payRange := strings.Split(td.Text(), " to ")
							pay1, _ := strconv.Atoi(strings.TrimSpace(payRange[0]))
							pay2, _ := strconv.Atoi(strings.TrimSpace(payRange[1]))
							inc.Salary = append(inc.Salary, pay1, pay2)
						} else {
							pay := strings.Replace(td.Text(), ",", "", -1)

							payAsNum, err := strconv.Atoi(pay)
							if err != nil {
								payAsNum = 0
							}

							inc.Salary = append(inc.Salary, payAsNum)
						}

					})
				}
				if len(inc.Salary) > 0 {
					p.Increments = append(p.Increments, inc)
				}
				p.Steps = len(inc.Salary)

				// match current increment here
				inForce, _ := time.Parse(time.RFC3339, date)
				if afterTimeSpan(inForce, time.Now()) {
					// today is after the in_force date for the pay agreement increments
					// should return false if inForce is in the future
					p.CurrentPayScale = inc.Salary
				}
			})
			g.PayScales = append(g.PayScales, p)
		}
	})
}
