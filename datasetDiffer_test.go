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
		{"testdata/orig.json", "testdata/orig.json", "", ""},
		{"testdata/orig.json", "testdata/newData.json", "Data Changed.", ""},
		{"testdata/orig.json", "testdata/newStructure.json", "Structure Changed.", ""},
		{"testdata/oldChecksum.json", "testdata/newChecksum.json", "Structure Changed.", ""},
		{"testdata/orig.json", "testdata/newDataAndStructure.json", "Structure Changed.", ""},
		{"testdata/orig.json", "testdata/blankStructure.json", "", "error: structure path cannot be empty string"},
		{"testdata/orig.json", "testdata/blankData.json", "", "error: data path cannot be empty string"},
		{"testdata/orig_with_meta.json", "testdata/newTitle.json", "Title Changed.", ""},
		{"testdata/orig_with_meta.json", "testdata/newDescription.json", "Description Changed.", ""},
		//TODO: test transform cahnge (need path)
		{"testdata/exampleD"}
		//TODO: test visconfig change (need path)
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
