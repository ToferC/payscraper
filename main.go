package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var groupList = []string{
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=4#rates-ec",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=1#rates-cs",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=10#rates-fb",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=11#rates-fi",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=15#rates-as",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=15#rates-pm",
}

type PayScale struct {
	Name       string
	Identifier string
	PayScale   map[string][]string
	URL        string
}

func main() {
	for _, group := range groupList {

		g := PayScale{
			URL: group,
		}

		GetPayScales(group, &g)
	}
}

// StringToLines - Convert HTML table strings into text lines
func StringToLines(s string) []string {
	var lines []string

	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return lines
}

func processTable(tableObject *goquery.Selection) {
	fmt.Println("Processing table")

	// map of level(int) x column(string) x value(string)
	payRates := make(map[string][]string)

	tableObject.Each(func(i int, table *goquery.Selection) {

		table.Find("tr").Each(func(rowIndex int, tr *goquery.Selection) {

			date := "2020-01-01"

			tr.Find("time").Each(func(indexOfTd int, th *goquery.Selection) {
				date, _ = th.Attr("datetime")
				payRates[date] = []string{}

			})

			if date != "2020-01-01" {
				tr.Find("td").Each(func(indexOfTd int, td *goquery.Selection) {

					pay := strings.Replace(td.Text(), ",", "", -1)

					payRates[date] = append(payRates[date], pay)

				})
			}

		})
	})
	fmt.Println(payRates)
}

func GetPayScales(groupURL string, g *PayScale) map[string][]string {

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

		goquerySelection.Find("tbody").Each(func(index int, tablehtml *goquery.Selection) {
			if index == 0 {
			} else {
				fmt.Println("Found Pay Table", index)
				processTable(tablehtml)
			}
		})
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(path)

	//fmt.Println(abilities)

	g.URL = groupURL
	g.PayScale = payRates

	fmt.Println(g)

	return payRates

}
