package main

type Group struct {
	Name       string
	Identifier string
	URL        string
	PayScales  []PayScale
}

type PayScale struct {
	Name     string
	Steps    int
	PayScale map[string][]int
}
