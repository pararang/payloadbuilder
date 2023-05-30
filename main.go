package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	// _ "github.com/mitchellh/mapstructure"
)

type FunderAPayload struct {
	FullName string `json:"full_name" mapstructure:"full_name"`
	FCR      int    `json:"fcr"`
}

type FunderBPayload struct {
	Borrower BorrowerData `json:"borrower"`
}

type BorrowerData struct {
	Name string `json:"name"`
	FCR  int    `json:"fcr"`
}

type configJSON map[string]string

type MappingConfig struct {
	FunderA configJSON `json:"funderA"`
	FunderB configJSON `json:"funderB"`
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

func iterateFields(data map[string]interface{}, prefix string) {
	for key, value := range data {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			// Recursively iterate over nested maps
			iterateFields(nestedMap, prefix+key+".")
		} else {
			// Handle leaf nodes (non-object values)
			fmt.Printf("Field: %s%s, Value: %v\n", prefix, key, value)
		}
	}
}

func mapDataToFunderPayload(data map[string]interface{}, funder string) (map[string]interface{}, error) {
	config, err := loadMappingConfig("mapping_config.json")
	if err != nil {
		return nil, err
	}

	mappingRules := configJSON{}
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
		val, err := getNestedValue(data, value)
		if err != nil {
			log.Println("failed get value for ", key)
			continue
		}

		
		payload[key] = val
	}

	return payload, nil

	// switch funder {
	// case "funderA":
	// 	funderAPayload := FunderAPayload{}
	// 	err = mapstructure.Decode(payload, &funderAPayload)
	// 	return funderAPayload, err
	// case "funderB":
	// 	funderBPayload := FunderBPayload{}
	// 	err = mapstructure.Decode(payload, &funderBPayload)
	// 	return funderBPayload, err
	// default:
	// 	return nil, fmt.Errorf("unsupported funder: %s", funder)
	// }
}

func getNestedValue(data map[string]interface{}, key string) (interface{}, error) {
	keys := splitKey(key)
	current := data

	for _, k := range keys {
		value, found := current[k]
		if !found {
			return nil, fmt.Errorf("Key '%s' not found", key)
		}

		// If the value is a nested map, continue traversing
		if nestedMap, ok := value.(map[string]interface{}); ok {
			current = nestedMap
		} else {
			return value, nil
		}
	}
	return current, nil
}

func splitKey(key string) []string {
	return strings.Split(key, ".")
}

func main() {
	jsonData := []byte(`{
		"credit_score": {"fcr": 10},
		"profile": {"name": "member name"},
		"district": "Batua",
		"city": "Bandung"
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