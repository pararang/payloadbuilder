package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type MappingConfig struct {
	FunderA map[string]string `json:"funderA"`
	FunderB map[string]string `json:"funderB"`
}

type FunderAPayload struct {
	FullName string `json:"full_name"`
	FCR      int    `json:"fcr"`
}

type FunderBPayload struct {
	Borrower BorrowerData `json:"borrower"`
}

type BorrowerData struct {
	Name string `json:"name"`
	FCR  int    `json:"fcr"`
}

func loadMappingConfig(filePath string) (MappingConfig, error) {
	var config MappingConfig
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func mapDataToFunderPayload(data map[string]interface{}, funder string) (interface{}, error) {
	config, err := loadMappingConfig("mapping_config.json")
	if err != nil {
		return nil, err
	}

	mappingRules := map[string]string{}
	switch funder {
	case "funderA":
		mappingRules = config.FunderA
	case "funderB":
		mappingRules = config.FunderB
	default:
		return nil, fmt.Errorf("unsupported funder: %s", funder)
	}

	payload := make(map[string]interface{})
	for key, value := range mappingRules {
		val, ok := getValueByJSONPath(data, value)
		if !ok {
			return nil, fmt.Errorf("mapping value not found for key: %s", key)
		}
		payload[key] = val
	}

	switch funder {
	case "funderA":
		funderAPayload := FunderAPayload{}
		err = mapstructure.Decode(payload, &funderAPayload)
		return funderAPayload, err
	case "funderB":
		funderBPayload := FunderBPayload{}
		err = mapstructure.Decode(payload, &funderBPayload)
		return funderBPayload, err
	default:
		return nil, fmt.Errorf("unsupported funder: %s", funder)
	}
}

func getValueByJSONPath(data map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	for i, part := range parts {
		value, ok := data[part]
		if !ok {
			return nil, false
		}

		if i == len(parts)-1 {
			return value, true
		}

		childData, ok := value.(map[string]interface{})
		if !ok {
			return nil, false
		}
		data = childData
	}

	return nil, false
}

func main() {
	jsonData := []byte(`{
		"credit_score": {"fcr": 10},
		"profile": {"name": "member name"}
	}`)

	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal(err)
	}

	funder := "funderA" // Change this based on the chosen funder

	payload, err := mapDataToFunderPayload(data, funder)
	if err != nil {
		log.Fatal(err)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonPayload))
}
