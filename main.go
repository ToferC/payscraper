package main

import "time"

var groupURLs = []string{
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=4#rates-ec",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=1#rates-cs",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=10#rates-fb",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=11#rates-fi",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=15#rates-as",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=15#rates-pm",
	"https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?id=15#rates-is",
}

func main() {

	urls := GetGroupURLs("https://www.tbs-sct.gc.ca/pubs_pol/hrpubs/coll_agre/rates-taux-eng.asp")

	for _, url := range urls {

		today := time.Now().String()

		g := Group{
			Identifier:  url[len(url)-2:],
			URL:         url,
			ScrapedDate: today,
		}

		GetPayScales(url, &g)
	}
}
