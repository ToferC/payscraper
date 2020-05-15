package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func (g Group) save() {

	irregString := ""

	if g.IrregularFormat == true {
		irregString = "irr_"
	}

	path := "./rates_of_pay_groups/" + irregString + g.Identifier + ".json"

	writeFile(path, g)
}

func saveGroupData(groups []Group, filename string) {
	path := fmt.Sprintf("./rates_of_pay_groups/%s", filename)

	// Create new file if needed
	var file, err = os.Create(path)
	checkError(err)
	defer file.Close()

	fmt.Println("==> done creating file", path+"\n")

	data, err := json.Marshal(groups)
	checkError(err)
	err = ioutil.WriteFile(path, data, 0644)
	checkError(err)

}

func writeFile(path string, g Group) {

	// Create new file if needed
	var file, err = os.Create(path)
	checkError(err)
	defer file.Close()

	fmt.Println("==> done creating file", path+"\n")

	data, err := json.Marshal(g)
	checkError(err)
	err = ioutil.WriteFile(path, data, 0644)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func afterTimeSpan(inForce, check time.Time) bool {
	return check.After(inForce)
}

func sum(slice []int) int {
	total := 0

	for i := 0; i < len(slice); i++ {
		total += slice[i]
	}

	return total
}
