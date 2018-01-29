package datasetDiffer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	// "github.com/qri-io/cafs/memfs"
	"github.com/qri-io/dataset"
	// "github.com/qri-io/dataset/dsfs"
	// "github.com/qri-io/qri/core"
	// "github.com/qri-io/qri/repo"
	// testrepo "github.com/qri-io/qri/repo/test"
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
		expected                string
		err                     string
	}{
		{"testdata/orig.json", "testdata/newStructure.json", "Structure Changed. (3 changes)", ""},
		{"testdata/orig.json", "testdata/newTitle.json", "Metadata Changed. (1 changes)", ""},
		{"testdata/orig.json", "testdata/newDescription.json", "Metadata Changed. (1 changes)", ""},
		{"testdata/orig.json", "testdata/newVisConfig.json", "VisConfig Changed. (1 changes)", ""},
		// {"testdata/orig.json", "testdata/newTransform.json", "Transform Changed. (1 changes)", ""},
		// {"testdata/orig.json", "testdata/newData.json", "Data Changed. (1 changes)", ""},
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
		got, err := DiffDatasets(dsLeft, dsRight)
		if err != nil {
			if err.Error() == c.err {
				continue
			} else {
				t.Errorf("case %d error mismatch: expected '%s', got '%s'", i, c.err, err.Error())
				return
			}
		}
		stringDiffs := MapDiffsToString(got)
		if i == 4 {
			// for k, v := range got {
			// 	fmt.Printf("%s: %s\n---\n", k, v)
			// }
			if dsLeft.Transform == nil {
				fmt.Println("left transform nil")
			}
			if dsRight.Transform == nil {
				fmt.Println("right transform nil")
			}
		}

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
		got, err := DiffJSON(aBytes, bBytes)
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
