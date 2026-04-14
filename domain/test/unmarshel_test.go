package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/atharvyadav96k/SPOTNEARR_SHARED/common"
)

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type TestCase struct {
	Name         string                 `json:"name"`
	Input        map[string]interface{} `json:"input"`
	ExpectedJSON string                 `json:"expected_json"`
}

type TestCases struct {
	TestCases []TestCase `json:"test_cases"`
}

func loadTestCases(t *testing.T) []TestCase {
	testTablePath := filepath.Join("..", "test_table", "test_cases.json")
	file, err := os.Open(testTablePath)
	if err != nil {
		t.Fatalf("Failed to open test cases file: %v", err)
	}
	defer file.Close()

	var tc TestCases
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tc); err != nil {
		t.Fatalf("Failed to decode test cases: %v", err)
	}

	return tc.TestCases
}

func mapsEqual(a, b map[string]interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func TestToJSONFromTable(t *testing.T) {
	testCases := loadTestCases(t)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			jsonBytes, err := common.ToJSON(tc.Input)
			if err != nil {
				t.Fatalf("ToJSON failed: %v", err)
			}

			// Unmarshal both expected and actual to compare as maps
			var expectedMap map[string]interface{}
			if err := json.Unmarshal([]byte(tc.ExpectedJSON), &expectedMap); err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}

			var actualMap map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &actualMap); err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			if !mapsEqual(expectedMap, actualMap) {
				t.Errorf("Expected %v, got %v", expectedMap, actualMap)
			}
		})
	}
}

func TestFromJSONFromTable(t *testing.T) {
	testCases := loadTestCases(t)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var result map[string]interface{}
			res, err := common.FromJSON[map[string]interface{}]([]byte(tc.ExpectedJSON))
			if err != nil {
				t.Fatalf("FromJSON failed: %v", err)
			}
			result = res

			// Compare the maps
			if !mapsEqual(result, tc.Input) {
				t.Errorf("Expected %v, got %v", tc.Input, result)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	ts := TestStruct{Name: "test", Value: 42}
	jsonBytes, err := common.ToJSON(ts)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	expected := `{"name":"test","value":42}`
	if string(jsonBytes) != expected {
		t.Errorf("Expected %s, got %s", expected, string(jsonBytes))
	}
}

func TestFromJSON(t *testing.T) {
	jsonStr := `{"name":"test","value":42}`
	// var ts TestStruct
	result, err := common.FromJSON[TestStruct]([]byte(jsonStr))
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if result.Name != "test" || result.Value != 42 {
		t.Errorf("Expected {test 42}, got %+v", result)
	}
}

func TestToJSONAndFromJSON(t *testing.T) {
	original := TestStruct{Name: "roundtrip", Value: 100}
	jsonBytes, err := common.ToJSON(original)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	result, err := common.FromJSON[TestStruct](jsonBytes)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if result != original {
		t.Errorf("Roundtrip failed: expected %+v, got %+v", original, result)
	}
}
