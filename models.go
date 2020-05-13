package main

type Group struct {
	Name       string     `json:"name"`
	Identifier string     `json:"identifier"`
	URL        string     `json:"url"`
	PayScales  []PayScale `json:"pay_scales"`
}

type PayScale struct {
	Name    string   `json:"name"`
	Steps   int      `json:"steps"`
	PayRows []PayRow `json:"pay_rows"`
}

type PayRow struct {
	DateTime string `json:"date_time"`
	Salary   []int  `json:"salary"`
}
