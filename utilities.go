package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func (g Group) save() {

	path := "./rates_of_pay_groups/" + g.Identifier + ".json"

	writeFile(path, g)
}

func writeFile(path string, g Group) {

	// Create new file if needed
	var file, err = os.Create(path)
	checkError(err)
	defer file.Close()

	fmt.Println("==> done creating file", path+"\n")

	data, _ := json.Marshal(g)
	err = ioutil.WriteFile(path, data, 0644)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}
