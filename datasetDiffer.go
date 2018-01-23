package datasetDiffer

import (
	"fmt"

	"github.com/qri-io/dataset"
)

// Diff contains a string description of a diff
type Diff string

// DiffList contains a slice of diffs in order of descending scope
type DiffList struct {
	diffs []Diff
}

// String returns the first (largest scope) change as a string
func (diffList DiffList) String() string {
	if len(diffList.diffs) > 0 {
		return string(diffList.diffs[0])
	}
	return ""
}

// DiffDatasets calculates diffs between two datasets and returns a
// dataset. Differences are checked in order of descending scope
// - dataset.Dataset.path
// - dataset.Dataset.Structure.path
// - dataset.Dataset.Data.path
// TODO: make diffs non-trivial
func DiffDatasets(a, b *dataset.Dataset) (*DiffList, error) {
	diffList := &DiffList{}
	diffDescription := Diff("")
	if len(a.Structure.Path().String()) <= 1 || len(b.Structure.Path().String()) <= 1 {
		return nil, fmt.Errorf("error: structure path cannot be empty string")
	}
	if len(a.DataPath) <= 1 || len(b.DataPath) <= 1 {
		return nil, fmt.Errorf("error: data path cannot be empty string")
	}
	if a.Structure.Path() != b.Structure.Path() {
		diffDescription = Diff("Structure Changed.")
		diffList.diffs = append(diffList.diffs, diffDescription)
	}
	if a.DataPath != b.DataPath {
		diffDescription = Diff("Data Changed.")
		diffList.diffs = append(diffList.diffs, diffDescription)
	}
	return diffList, nil
}
