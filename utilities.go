package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func (g Group) save() {

	path := "./groups/" + g.Identifier + ".json"

	writeFile(path, g)
}

func writeFile(path string, g Group) {

	// Check if file exists
	var _, err = os.Stat(path)

	// Create new file if needed
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		checkError(err)
		defer file.Close()
	}
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
