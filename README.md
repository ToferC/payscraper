# PayScraper - Web Scraper for Government of Canada rates of pay
This go project is a Web scraper designed to feed into an API.

The Web scraper that crawls the [TBS Rates of pay for public service employees page](https://www.tbs-sct.gc.ca/pubs_pol/hrpubs/coll_agre/rates-taux-eng.asp), pulls down pay-date and transforms it into a machine readable format. It can currently parse 40 collective agreements and extract pay-related data.

You can find the scraped data in the /rates_of_pay_groups folder.

This data is then used to power the gc-payscales API, available here: [https://gc-payscales.herokuapp.com/playground](https://gc-payscales.herokuapp.com/playground)

### Sample
![Rates of Pay sample scrape](https://github.com/ToferC/payscraper/raw/master/payscraper.png)

Initial testing and improvements to the scraper are still ongoing.
