package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMergeMaps(t *testing.T) {
	dest := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"subkey1": "subvalue1",
		},
	}

	src := map[string]interface{}{
		"key2": map[string]interface{}{
			"subkey2": "subvalue2",
		},
		"key3": "value3",
	}

	expected := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"subkey1": "subvalue1",
			"subkey2": "subvalue2",
		},
		"key3": "value3",
	}

	result := mergeMaps(dest, src)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Merge result is incorrect. Expected: %v, Got: %v", expected, result)
	}
}

func TestCreateNestedMap(t *testing.T) {
	keys := []string{"address", "city"}
	value := "Bandung"

	expected := map[string]interface{}{
		"address": map[string]interface{}{
			"city": "Bandung",
		},
	}

	result := createNestedMap(keys, value)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Nested map creation is incorrect. Expected: %v, Got: %v", expected, result)
	}
}

func TestGetNestedValue(t *testing.T) {
	data := map[string]interface{}{
		"address": map[string]interface{}{
			"city":     "Bandung",
			"district": "Batua",
		},
		"fcr":       10,
		"full_name": "member name",
	}

	t.Run("Existing Key", func(t *testing.T) {
		key := "address.city"
		expected := "Bandung"

		result, err := getNestedValue(data, key)

		if err != nil {
			t.Errorf("Failed to get nested value: %v", err)
		}

		if fmt.Sprintf("%v", result) != expected {
			t.Errorf("Nested value is incorrect. Expected: %v, Got: %v", expected, result)
		}
	})

	t.Run("Non-existing Key", func(t *testing.T) {
		key := "address.country"

		result, err := getNestedValue(data, key)

		if err == nil {
			t.Error("Expected error, but got none")
		}

		if result != nil {
			t.Errorf("Expected nil value, but got: %v", result)
		}

		expectedError := fmt.Sprintf("Key '%s' not found", key)
		if err.Error() != expectedError {
			t.Errorf("Error message is incorrect. Expected: %s, Got: %v", expectedError, err)
		}
	})

	t.Run("Last Key", func(t *testing.T) {
		key := "full_name"
		expected := "member name"

		result, err := getNestedValue(data, key)

		if err != nil {
			t.Errorf("Failed to get nested value: %v", err)
		}

		if fmt.Sprintf("%v", result) != expected {
			t.Errorf("Nested value is incorrect. Expected: %v, Got: %v", expected, result)
		}
	})
}