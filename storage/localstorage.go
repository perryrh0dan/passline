package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Get data of website
func Get(website string) (Item, error) {
	data := getData()
	for i := 0; i < len(data.Items); i++ {
		if data.Items[i].Website == website {
			return data.Items[i], nil
		}
	}

	return Item{}, fmt.Errorf("No entry for website %s", website)
}

func getData() ItemStorage {
	file, _ := ioutil.ReadFile("test.json")

	data := ItemStorage{}

	_ = json.Unmarshal([]byte(file), &data)

	return data
}
