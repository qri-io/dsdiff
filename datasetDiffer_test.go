package datasetDiffer

import (
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
	//setup
	// mr, err := testrepo.NewTestRepo()
	// if err != nil {
	// 	fmt.Printf("error allocating test repo: %s\n", err.Error())
	// 	return
	// }
	// // make new request
	// req := core.NewDatasetRequests(mr, nil)
	// File 1
	// dsRef1 := &repo.DatasetRef{}
	// initParams := &core.InitParams{
	// 	DataFilename: jobsByAutomationFile.FileName(),
	// 	Data:         jobsByAutomationFile,
	// 	// MetadataFilename: jobsMeta.FileName(),
	// 	// Metadata:         jobsMeta,
	// }
	// err = req.Init(initParams, dsRef1)
	// if err != nil {
	// 	fmt.Println("couldn't load file 1")
	// }
	// dsBase, err := dsfs.LoadDataset(mr.Store(), dsRef1.Path)
	// if err != nil {
	// 	fmt.Printf("error loading dataset 1: %s", err.Error())
	// }
	// // File 2
	// dsRef2 := &repo.DatasetRef{}
	// initParams = &core.InitParams{
	// 	DataFilename: jobsByAutomationFile2.FileName(),
	// 	Data:         jobsByAutomationFile2,
	// }
	// err = req.Init(initParams, dsRef2)
	// if err != nil {
	// 	fmt.Println("couldn't load second file")
	// }
	// dsNewStructure, err := dsfs.LoadDataset(mr.Store(), dsRef2.Path)
	// if err != nil {
	// 	fmt.Println("error loading dataset: %s", err.Error())
	// }
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

// func TestDiffDataset(t *testing.T) {
// 	//setup

// 	//test cases
// 	cases := []struct {
// 		dsPathA  string
// 		dsPathB  string
// 		expected string
// 		err      string
// 	}{
// 		{"testdata/orig.json", "testdata/orig.json", "", ""},
// 		{"testdata/orig.json", "testdata/newData.json", "Data Changed.", ""},
// 		{"testdata/orig.json", "testdata/newStructure.json", "Structure Changed.", ""},
// 		{"testdata/oldChecksum.json", "testdata/newChecksum.json", "Structure Changed.", ""},
// 		{"testdata/orig.json", "testdata/newDataAndStructure.json", "Structure Changed.", ""},
// 		{"testdata/orig.json", "testdata/blankStructure.json", "", "error: structure path cannot be empty string"},
// 		{"testdata/orig.json", "testdata/blankData.json", "", "error: data path cannot be empty string"},
// 		{"testdata/orig_with_meta.json", "testdata/newTitle.json", "Title Changed.", ""},
// 		{"testdata/orig_with_meta.json", "testdata/newDescription.json", "Description Changed.", ""},
// 		//TODO: test transform change (need path)
// 		// {"orig_with_vis_config_and_transform.json", "orig_with_vis_config_and_transform.json", "Transform Changed.", ""}
// 		//TODO: test visconfig change (need path)
// 		// {"orig_with_vis_config_and_transform.json", "orig_with_vis_config_and_transform.json", "VisConfig Changed.", ""}
// 	}
// 	// load files and execute tests
// 	for i, c := range cases {
// 		pathA := c.dsPathA
// 		pathB := c.dsPathB
// 		//load data
// 		dsA, err := loadTestData(pathA)
// 		if err != nil {
// 			t.Errorf("case %d error: error loading file '%s'", i, pathA)
// 			return
// 		}
// 		dsB, err := loadTestData(pathB)
// 		if err != nil {
// 			t.Errorf("case %d error: error loading file '%s'", i, pathB)
// 			return
// 		}
// 		got, err := DiffDatasets(dsA, dsB)
// 		if err != nil {
// 			if err.Error() == c.err {
// 				continue
// 			} else {
// 				t.Errorf("case %d error mismatch: expected '%s', got '%s'", i, c.err, err.Error())
// 				return
// 			}

// 		}
// 		if got.String() != c.expected {
// 			t.Errorf("case %d response mismatch.  expected '%s', got '%s'", i, c.expected, got.String())
// 			continue
// 		}
// 	}
// }

// func TestDiffStructure2(t *testing.T) {
// 	//test cases
// 	cases := []struct {
// 		dsPathA  string
// 		dsPathB  string
// 		expected string
// 		err      string
// 	}{
// 		{"testdata/structureJsonSchemaOrig.json", "testdata/structureJsonSchemaNew.json", "", ""},
// 	}
// 	for i, c := range cases {
// 		pathA := c.dsPathA
// 		pathB := c.dsPathB
// 		//load data
// 		dsA, err := loadTestData(pathA)
// 		if err != nil {
// 			t.Errorf("case %d error: error loading file '%s'", i, pathA)
// 			return
// 		}
// 		dsB, err := loadTestData(pathB)
// 		if err != nil {
// 			t.Errorf("case %d error: error loading file '%s'", i, pathB)
// 			return
// 		}
// 		got, err := diffStructure2(dsA.Structure, dsB.Structure, true)
// 		if err != nil {
// 			if err.Error() == c.err {
// 				continue
// 			} else {
// 				t.Errorf("case %d error mismatch: expected '%s', got '%s'", i, c.err, err.Error())
// 				return
// 			}
// 		}
// 		fmt.Println("------------------")
// 		fmt.Print(got.JSONDiffAsciiStr)
// 		fmt.Println("------------------")
// 		if got.JSONDiffAsciiStr != c.expected {
// 			t.Errorf("case %d response mismatch.  expected '%s', got '%s'", i, c.expected, got.JSONDiffAsciiStr)
// 			continue
// 		}
// 	}
// }

// formatter := formatter.NewAsciiFormatter(sd1Json, config)
// diffString, err := formatter.Format(d)
// if err != nil {
// 	fmt.Printf("something went wrong: %s", err.Error())
// 	return
// }
