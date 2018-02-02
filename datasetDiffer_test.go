package datasetDiffer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/qri-io/dataset"
)

func loadTestData(path string) (*dataset.Dataset, error) {
	dataBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	d := &dataset.Dataset{}
	err = d.UnmarshalJSON(dataBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal dataset: %s", err.Error())
	}
	return d, nil
}

func TestDiffDataset(t *testing.T) {
	//test cases
	cases := []struct {
		dsLeftPath, dsRightPath string
		displayFormat           string
		expected                string
		err                     string
	}{
		{"testdata/orig.json", "testdata/newStructure.json", "listKeys", "Structure: 3 changes\n\t- modified checksum\n\t- modified entries\n\t- modified schema", ""},
		{"testdata/orig.json", "testdata/newTitle.json", "listKeys", "Meta: 1 change\n\t- modified title", ""},
		{"testdata/orig.json", "testdata/newDescription.json", "listKeys", "Meta: 1 change\n\t- modified description", ""},
		{"testdata/orig.json", "testdata/newVisConfig.json", "listKeys", "VisConfig: 1 change\n\t- modified format", ""},
		// {"testdata/orig.json", "testdata/newTransform.json", "simple", "Transform Changed. (1 change)", ""},
		// {"testdata/orig.json", "testdata/newData.json", "simple", "Data Changed. (1 change)", ""},
	}
	// execute
	for i, c := range cases {
		//Load data
		dsLeft, err := loadTestData(c.dsLeftPath)
		if err != nil {
			t.Errorf("case %d error: error loading file '%s'", i, c.dsLeftPath)
			return
		}
		dsRight, err := loadTestData(c.dsRightPath)
		if err != nil {
			t.Errorf("case %d error: error loading file '%s'", i, c.dsRightPath)
			return
		}
		got, err := DiffDatasets(dsLeft, dsRight, nil)
		if err != nil {
			if err.Error() == c.err {
				continue
			} else {
				t.Errorf("case %d error mismatch: expected '%s', got '%s'", i, c.err, err.Error())
				return
			}
		}
		stringDiffs, err := MapDiffsToString(got, c.displayFormat)
		if err != nil {
			t.Errorf("case %d error: %s", i, err.Error())
			return
		}
		// if i == 0 {
		// 	s, err := MapDiffsToFormattedString(got, dsLeft)
		// 	if err != nil {
		// 		t.Errorf("not today: %s", err.Error())
		// 	}
		// 	fmt.Println("--------------------------")
		// 	fmt.Print(s)
		// 	fmt.Println("--------------------------")
		// }

		if stringDiffs != c.expected {
			t.Errorf("case %d response mistmatch: expected '%s', got '%s'", i, c.expected, stringDiffs)
		}
	}
}

func TestDiffJSON(t *testing.T) {
	//test cases
	cases := []struct {
		dsLeftPath, dsRightPath string
		description             string
		expected                string
		err                     string
	}{
		{"testdata/orig.json", "testdata/newStructure.json", "abc", "1 diffs", ""},
	}
	// execute
	for i, c := range cases {
		//Load data
		a, err := loadTestData(c.dsLeftPath)
		if err != nil {
			t.Errorf("case %d error: error loading file '%s'", i, c.dsLeftPath)
			return
		}
		b, err := loadTestData(c.dsRightPath)
		if err != nil {
			t.Errorf("case %d error: error loading file '%s'", i, c.dsRightPath)
			return
		}
		aBytes, err := json.Marshal(a)
		if err != nil {
			t.Errorf("error marshalling structure a: %s", err.Error())
			return
		}
		bBytes, err := json.Marshal(b)
		if err != nil {
			t.Errorf("error marshalling structure b: %s", err.Error())
			return
		}
		got, err := DiffJSON(aBytes, bBytes, c.description)
		if err != nil {
			if err.Error() == c.err {
				continue
			} else {
				t.Errorf("case %d error mismatch: expected '%s', got '%s'", i, c.err, err.Error())
				return
			}
		}
		stringDiffs := fmt.Sprintf("%d diffs", len(got.Deltas()))
		if stringDiffs != c.expected {
			t.Errorf("case %d response mistmatch: expected '%s', got '%s'", i, c.expected, stringDiffs)
		}
	}
}
