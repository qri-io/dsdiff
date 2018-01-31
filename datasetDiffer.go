package datasetDiffer

import (
	"encoding/json"
	"fmt"

	"github.com/qri-io/dataset"
	diff "github.com/yudai/gojsondiff"
)

// DiffStructure diffs the structure of two datasets
func DiffStructure(a, b *dataset.Structure) (diff.Diff, error) {
	var emptyDiff diff.Diff
	differ := diff.New()
	if len(a.Path().String()) > 1 && len(b.Path().String()) > 1 {
		if a.Path() == b.Path() {
			return emptyDiff, nil
		}
	}
	if len(a.Checksum) > 1 && len(b.Checksum) > 1 {
		if a.Checksum == b.Checksum {
			return emptyDiff, nil
		}
	}
	// If we couldn't determine that there were no changes  using the
	// path or checksum...
	aBytes, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("error marshalling structure a: %s", err.Error())
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("error marshalling structure b: %s", err.Error())
	}
	d, err := differ.Compare(aBytes, bBytes)
	if err != nil {
		return nil, fmt.Errorf("error comparing structure: %s", err.Error())
	}
	return d, nil
}

// DiffData diffs the data of two datasets
func DiffData(a, b *dataset.Dataset) (diff.Diff, error) {
	var emptyDiff diff.Diff
	// differ := diff.New()
	if len(a.DataPath) > 1 && len(b.DataPath) > 1 {
		if a.DataPath == b.DataPath {
			return emptyDiff, nil
		}
	}
	// TODO: dereference DataPath and pass to jsondiffer
	return emptyDiff, nil
}

// DiffTransform diffs the transform struct of two datasets
func DiffTransform(a, b *dataset.Transform) (diff.Diff, error) {
	var emptyDiff diff.Diff
	differ := diff.New()
	if len(a.Path().String()) > 1 && len(b.Path().String()) > 1 {
		if a.Path() == b.Path() {
			return emptyDiff, nil
		}
	}
	aBytes, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("error marshalling transform a: %s", err.Error())
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("error marshalling transform b: %s", err.Error())
	}
	d, err := differ.Compare(aBytes, bBytes)
	if err != nil {
		return nil, fmt.Errorf("error comparing transforms: %s", err.Error())
	}
	return d, nil
}

// DiffMeta diffs the metadata of two datasets
func DiffMeta(a, b *dataset.Meta) (diff.Diff, error) {
	var emptyDiff diff.Diff
	differ := diff.New()
	if len(a.Path().String()) > 1 && len(b.Path().String()) > 1 {
		if a.Path() == b.Path() {
			return emptyDiff, nil
		}
	} else if a.IsEmpty() && b.IsEmpty() {
		return emptyDiff, nil
	}

	aBytes, err := a.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshaling meta a: %s", err.Error())
	}
	bBytes, err := b.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshaling meta b: %s", err.Error())
	}
	d, err := differ.Compare(aBytes, bBytes)
	if err != nil {
		return nil, fmt.Errorf("error comparing Meta: %s", err.Error())
	}
	return d, nil
}

// DiffVisConfig diffs the dataset.VisConfig structs of two datasets
func DiffVisConfig(a, b *dataset.VisConfig) (diff.Diff, error) {
	var emptyDiff diff.Diff
	differ := diff.New()
	if len(a.Path().String()) > 1 && len(b.Path().String()) > 1 {
		if a.Path() == b.Path() {
			return emptyDiff, nil
		}
	}
	aBytes, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("error marshalling visConfig a: %s", err.Error())
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("error marshalling visConfig b: %s", err.Error())
	}
	d, err := differ.Compare(aBytes, bBytes)
	if err != nil {
		return nil, fmt.Errorf("error comparing VisConfigs: %s", err.Error())
	}
	return d, nil
}

// DiffDatasets returns a map of diffs of the components of a dataset
func DiffDatasets(a, b *dataset.Dataset) (map[string]diff.Diff, error) {
	result := map[string]diff.Diff{}
	//diff structure
	if a.Structure != nil && b.Structure != nil {
		structureDiffs, err := DiffStructure(a.Structure, b.Structure)
		if err != nil {
			return nil, err
		}
		result["structure"] = structureDiffs
	}
	// diff data
	dataDiffs, err := DiffData(a, b)
	if err != nil {
		return nil, err
	}
	result["data"] = dataDiffs
	// diff transform
	if a.Transform != nil && b.Transform != nil {
		transformDiffs, err := DiffTransform(a.Transform, b.Transform)
		if err != nil {
			return nil, err
		}
		result["transform"] = transformDiffs
	}
	// diff meta
	if a.Meta != nil && b.Meta != nil {
		metaDiffs, err := DiffMeta(a.Meta, b.Meta)
		if err != nil {
			return nil, err
		}
		result["meta"] = metaDiffs
	}
	// diff visConfig
	if a.VisConfig != nil && b.VisConfig != nil {
		visConfigDiffs, err := DiffVisConfig(a.VisConfig, b.VisConfig)
		if err != nil {
			return nil, err
		}
		result["visConfig"] = visConfigDiffs
	}
	return result, nil
}

// DiffJSON diffs two json files independent of any Dataset structures
func DiffJSON(a, b []byte) (diff.Diff, error) {
	differ := diff.New()
	d, err := differ.Compare(a, b)
	if err != nil {
		return nil, fmt.Errorf("error comparing json: %s", err.Error())
	}
	return d, nil
}

// MapDiffsToString generates a string description from a map of diffs
// Currently the String generated reflects the first/highest priority
// change made.  The priority of changes currently are
//   1. dataset.Structure
//   2. dataset.{Data} // TODO: use dereferenced data
//   3. dataset.Transform
//   4. dataset.Meta
//   5. Dataset.VisConfig
func MapDiffsToString(m map[string]diff.Diff) string {
	if m["structure"] != nil {
		structureDiffs := m["structure"]
		deltas := structureDiffs.Deltas()
		if len(deltas) > 0 {
			// for i, d := range deltas {
			// 	fmt.Printf("%d. %s: (%T)\n", i+1, d)
			// }
			return fmt.Sprintf("Structure Changed. (%d changes)", len(deltas))
		}
	}
	if m["data"] != nil {
		dataDiffs := m["data"]
		deltas := dataDiffs.Deltas()
		if len(deltas) > 0 {
			return fmt.Sprintf("Data Changed. (%d changes)", len(deltas))
		}
	}
	if m["transform"] != nil {
		transformDiffs := m["transform"]
		deltas := transformDiffs.Deltas()
		if len(deltas) > 0 {
			return fmt.Sprintf("Transform Changed. (%d changes)", len(deltas))
		}
	}
	if m["meta"] != nil {
		metaDiffs := m["meta"]
		deltas := metaDiffs.Deltas()
		if len(deltas) > 0 {
			return fmt.Sprintf("Metadata Changed. (%d changes)", len(deltas))
		}
	}
	if m["visConfig"] != nil {
		visConfigDiffs := m["visConfig"]
		deltas := visConfigDiffs.Deltas()
		if len(deltas) > 0 {
			return fmt.Sprintf("VisConfig Changed. (%d changes)", len(deltas))
		}
	}
	return ""
}
