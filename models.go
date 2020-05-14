package main

type Group struct {
	Name        string     `json:"name"`
	Identifier  string     `json:"identifier"`
	URL         string     `json:"url"`
	PayScales   []PayScale `json:"pay_scales"`
	ScrapedDate string     `json:"date_scraped"`
}

type PayScale struct {
	Name            string      `json:"name"`
	Steps           int         `json:"steps"`
	CurrentPayScale []int       `json:"current_pay"`
	Increments      []Increment `json:"increments"`
}

type Increment struct {
	DateTime string `json:"date_time"`
	Salary   []int  `json:"salary"`
}
