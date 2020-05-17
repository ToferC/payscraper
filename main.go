package main

import (
	"strings"
	"time"
)

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

	groups := []Group{}
	errorGroups := []Group{}

	for _, url := range urls {

		if strings.Contains(url, "https://www.tbs-sct.gc.ca/agreements-conventions/view-visualiser-eng.aspx?") {

			today := time.Now().String()

			g := Group{
				Identifier:      strings.ToUpper(url[len(url)-2:]),
				URL:             url,
				ScrapedDate:     today,
				IrregularFormat: false,
			}

			GetPayScales(url, &g)

			if g.PayScales == nil {
				// something is wrong
				g.IrregularFormat = true
			}

			if g.PayScales != nil {
				if len(g.PayScales) == 0 {
					g.IrregularFormat = true
				} else {
					if g.PayScales[0].Steps == 0 {
						g.IrregularFormat = true
					}
				}
			}

			if g.IrregularFormat == true {
				errorGroups = append(errorGroups, g)
				g.save()
			} else {
				// should be okay
				groups = append(groups, g)
				g.save()
			}
		}
	}

	saveGroupData(groups, "groups_data.json")
	saveGroupData(errorGroups, "error_groups_data.json")
}
