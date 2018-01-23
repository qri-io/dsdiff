package datasetDiffer

import (
	"fmt"
	"io/ioutil"
	// "path/filepath"
	"testing"
	// 	"github.com/qri-io/qri"
	// 	testrepo "github.com/qri-io/qri/repo/test"
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
	//setup

	//test cases
	cases := []struct {
		dsPathA  string
		dsPathB  string
		expected string
		err      string
	}{
		{"testdata/exampleData_orig.json", "testdata/exampleData_orig.json", "", ""},
		{"testdata/exampleData_orig.json", "testdata/exampleData_newData.json", "Data Changed.", ""},
		{"testdata/exampleData_orig.json", "testdata/exampleData_newStructure.json", "Structure Changed.", ""},
		{"testdata/exampleData_orig.json", "testdata/exampleData_newDataAndStructure.json", "Structure Changed.", ""},
		{"testdata/exampleData_orig.json", "testdata/exampleData_blankStructure.json", "", "error: structure path cannot be empty string"},
		{"testdata/exampleData_orig.json", "testdata/exampleData_blankData.json", "", "error: data path cannot be empty string"},
	}
	// load files and execute tests
	for i, c := range cases {
		pathA := c.dsPathA
		pathB := c.dsPathB
		//load data
		dsA, err := loadTestData(pathA)
		if err != nil {
			t.Errorf("case %d error: error loading file '%s'", i, pathA)
			return
		}
		dsB, err := loadTestData(pathB)
		if err != nil {
			t.Errorf("case %d error: error loading file '%s'", i, pathB)
			return
		}
		got, err := DiffDatasets(dsA, dsB)
		if err != nil {
			if err.Error() == c.err {
				continue
			} else {
				t.Errorf("case %d error mismatch: expected '%s', got '%s'", i, c.err, err.Error())
				return
			}

		}
		if got.String() != c.expected {
			t.Errorf("case %d response mismatch.  expected '%s', got '%s'", i, c.expected, got.String())
			continue
		}
	}
}
