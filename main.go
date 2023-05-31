package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)


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
		
		keys := splitKey(key)
		nestedMap := createNestedMap(keys, val)

		// Merge the nested map into the payload
		payload = mergeMaps(payload, nestedMap)
	}

	return payload, nil
}
func createNestedMap(keys []string, value interface{}) map[string]interface{} {
	nestedMap := make(map[string]interface{})
	current := nestedMap

	for i := 0; i < len(keys)-1; i++ {
		current[keys[i]] = make(map[string]interface{})
		current = current[keys[i]].(map[string]interface{})
	}

	current[keys[len(keys)-1]] = value

	return nestedMap
}

func mergeMaps(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		if _, found := dest[k]; found {
			// If the key already exists, recursively merge the nested maps
			if nestedMap, ok := dest[k].(map[string]interface{}); ok {
				if srcNestedMap, ok := v.(map[string]interface{}); ok {
					dest[k] = mergeMaps(nestedMap, srcNestedMap)
					continue
				}
			}
		}

		dest[k] = v
	}

	return dest
}

func getNestedValue(data map[string]interface{}, key string) (interface{}, error) {
	keys := splitKey(key)
	current := data

	for i:=0; i<len(keys); i++ {
		k := keys[i]
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
	return nil, nil
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