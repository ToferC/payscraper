package main

import "strings"

// Group represents a pay group as defined by a collective agreement
type Group struct {
	Name            string     `json:"name"`
	Identifier      string     `json:"identifier"`
	URL             string     `json:"url"`
	PayScales       []PayScale `json:"pay_scales"`
	ScrapedDate     string     `json:"date_scraped"`
	IrregularFormat bool       `json:"irregular_format"`
}

func (g Group) existsInPayScaleNames(target string) bool {

	namesStrings := ""
	for _, p := range g.PayScales {
		namesStrings += p.Name + " "
	}

	return strings.Contains(namesStrings, target)
}

// PayScale contains a specific level and agreed pay rates for a period of time and pay steps.
type PayScale struct {
	Name       string      `json:"name"`
	Level      int         `json:"level"`
	Steps      int         `json:"steps"`
	RatesOfPay []RateOfPay `json:"rates_of_pay"`
}

// RateOfPay for a collective agreement at a point in time across several pay steps. Includes a date_time for when the rate of pay comes into force and an array of salary steps.
type RateOfPay struct {
	DateTime string `json:"date_time"`
	Salary   []int  `json:"salary"`
}
