package datasetDiffer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qri-io/dataset"
	jdiff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

// SubDiff holds the diffs of a Dataset Subcomponent diff
type SubDiff struct {
	jdiff.Diff
	kind string
	a, b []byte
}

// SummarizeToString outputs a substring in a one of a few formats
// - simple (single line describing which component and how many
//   changes)
// - listKeys (lists the keys of what changed)
// - plusMinusColor (git-style plus/minus printout)
// - plusMinus (same as plusMinusColor without color)
func (d *SubDiff) SummarizeToString(how string) (string, error) {
	color := false
	if strings.Contains(how, "Color") {
		color = true
	}
	pluralS := ""
	if d != nil && d.Deltas() != nil && len(d.Deltas()) > 1 {
		pluralS = "s"
	}
	switch how {
	case "simple":
		if d.Modified() {
			componentTitle := strings.Title(d.kind)
			return fmt.Sprintf("%s Changed. (%d change%s)", componentTitle, len(d.Deltas()), pluralS), nil
		}
	case "listKeys":
		if d.Modified() {
			componentTitle := strings.Title(d.kind)
			namedDiffs := ""
			for _, del := range d.Deltas() {
				namedDiffs = fmt.Sprintf("%s\n\t- modified %s", namedDiffs, del)
			}
			return fmt.Sprintf("%s: %d change%s%s", componentTitle, len(d.Deltas()), pluralS, namedDiffs), nil
		}
	case "plusMinusColor", "plusMinus":
		if d.Modified() {
			var aJSON map[string]interface{}
			err := json.Unmarshal(d.a, &aJSON)
			if err != nil {
				return "", fmt.Errorf("error summarizing: %s", err.Error())
			}
			config := formatter.AsciiFormatterConfig{
				ShowArrayIndex: true,
				Coloring:       color,
			}
			form := formatter.NewAsciiFormatter(aJSON, config)
			diffString, err := form.Format(d)
			if err != nil {
				return "", fmt.Errorf("error summarizing: %s", err.Error())
			}
			return diffString, nil
		}
	case "delta":
		if d.Modified() {
			form := formatter.NewDeltaFormatter()
			diffString, err := form.Format(d)
			if err != nil {
				return "", fmt.Errorf("error summarizing: %s", err.Error())
			}
			return diffString, nil
		}
	default:
		return "", nil
	}
	return "", nil
}

// DiffStructure diffs the structure of two datasets
func DiffStructure(a, b *dataset.Structure) (*SubDiff, error) {
	var emptyDiff = &SubDiff{kind: "structure"}
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
	aBytes, err := a.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshalling structure a: %s", err.Error())
	}
	bBytes, err := b.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshalling structure b: %s", err.Error())
	}
	return DiffJSON(aBytes, bBytes, emptyDiff.kind)
}

// DiffData diffs the data of two datasets
func DiffData(a, b *dataset.Dataset) (*SubDiff, error) {
	var emptyDiff = &SubDiff{kind: "data"}
	// differ := jdiff.New()
	if len(a.DataPath) > 1 && len(b.DataPath) > 1 {
		if a.DataPath == b.DataPath {
			return emptyDiff, nil
		}
	}
	// TODO: dereference DataPath and pass to jsondiffer
	return emptyDiff, nil
}

// DiffTransform diffs the transform struct of two datasets
func DiffTransform(a, b *dataset.Transform) (*SubDiff, error) {
	var emptyDiff = &SubDiff{kind: "transform"}
	if len(a.Path().String()) > 1 && len(b.Path().String()) > 1 {
		if a.Path() == b.Path() {
			return emptyDiff, nil
		}
	}
	aBytes, err := a.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshalling transform a: %s", err.Error())
	}
	bBytes, err := b.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshalling transform b: %s", err.Error())
	}
	return DiffJSON(aBytes, bBytes, emptyDiff.kind)
}

// DiffMeta diffs the metadata of two datasets
func DiffMeta(a, b *dataset.Meta) (*SubDiff, error) {
	var emptyDiff = &SubDiff{kind: "meta"}
	if len(a.Path().String()) > 1 && len(b.Path().String()) > 1 {
		if a.Path() == b.Path() {
			return emptyDiff, nil
		}
	} else if a.IsEmpty() && b.IsEmpty() {
		return emptyDiff, nil
	}
	aBytes, err := a.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshalling meta a: %s", err.Error())
	}
	bBytes, err := b.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshalling meta b: %s", err.Error())
	}
	return DiffJSON(aBytes, bBytes, emptyDiff.kind)
}

// DiffVisConfig diffs the dataset.VisConfig structs of two datasets
func DiffVisConfig(a, b *dataset.VisConfig) (*SubDiff, error) {
	var emptyDiff = &SubDiff{kind: "visConfig"}
	if len(a.Path().String()) > 1 && len(b.Path().String()) > 1 {
		if a.Path() == b.Path() {
			return emptyDiff, nil
		}
	}
	aBytes, err := a.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshalling visConfig a: %s", err.Error())
	}
	bBytes, err := b.MarshalJSONObject()
	if err != nil {
		return nil, fmt.Errorf("error marshalling visConfig b: %s", err.Error())
	}
	return DiffJSON(aBytes, bBytes, emptyDiff.kind)
}

//DiffJSON diffs two json byte slices and returns a SubDiff pointer
func DiffJSON(a, b []byte, kind string) (*SubDiff, error) {
	differ := jdiff.New()
	d, err := differ.Compare(a, b)
	if err != nil {
		// return emptyDiff, fmt.Errorf("error comparing %s: %s", kind, err.Error())
		return nil, fmt.Errorf("error comparing %s: %s", kind, err.Error())
	}
	subDiff := &SubDiff{d, kind, a, b}
	return subDiff, nil
}

// StructuredDataTuple provides an additional input for DiffDatasets
// to use fully de-referenced dataset.data so that we can consider
// changes in dataset.Data beyond the hash/path being similar or
// different
type StructuredDataTuple struct {
	a, b *[]byte
}

// DiffDatasets returns a map of pointers to diffs of the components
// of a dataset.  It calls each of the Diff{Component} functions and
// adds the option for including de-referenced dataset.Data via
// the StructuredDataTuple
func DiffDatasets(a, b *dataset.Dataset, deRefData *StructuredDataTuple) (map[string]*SubDiff, error) {
	result := make(map[string]*SubDiff)
	//diff structure
	if a.Structure != nil && b.Structure != nil {
		structureDiffs, err := DiffStructure(a.Structure, b.Structure)
		if err != nil {
			return result, err
		}
		if structureDiffs.Diff != nil {
			result[structureDiffs.kind] = structureDiffs
		}
	}
	// diff data
	if deRefData != nil {
		dataDiffs, err := DiffJSON(*deRefData.a, *deRefData.b, "data")
		if err != nil {
			return nil, err
		}
		result[dataDiffs.kind] = dataDiffs
	} else {
		dataDiffs, err := DiffData(a, b)
		if err != nil {
			return nil, err
		}
		if dataDiffs.Diff != nil {
			result[dataDiffs.kind] = dataDiffs
		}
	}
	// diff meta
	if a.Meta != nil && b.Meta != nil {
		metaDiffs, err := DiffMeta(a.Meta, b.Meta)
		if err != nil {
			return nil, err
		}
		if metaDiffs.Diff != nil {
			result[metaDiffs.kind] = metaDiffs
		}
	}
	// diff transform
	if a.Transform != nil && b.Transform != nil {
		transformDiffs, err := DiffTransform(a.Transform, b.Transform)
		if err != nil {
			return nil, err
		}
		if transformDiffs.Diff != nil {
			result[transformDiffs.kind] = transformDiffs
		}
	}
	// diff visConfig
	if a.VisConfig != nil && b.VisConfig != nil {
		visConfigDiffs, err := DiffVisConfig(a.VisConfig, b.VisConfig)
		if err != nil {
			return nil, err
		}
		if visConfigDiffs.Diff != nil {
			result[visConfigDiffs.kind] = visConfigDiffs
		}
	}
	return result, nil
}

// MapDiffsToString generates a string description from a map of diffs
// Currently the String generated reflects the first/highest priority
// change made.  The priority of changes currently are
//   1. dataset.Structure
//   2. dataset.{Data}
//   3. dataset.Transform
//   4. dataset.Meta
//   5. Dataset.VisConfig
func MapDiffsToString(m map[string]*SubDiff, how string) (string, error) {
	keys := []string{
		"structure",
		"data",
		"transform",
		"meta",
		"visConfig",
	}
	// for _, key := range keys {
	// 	val, ok := m[key]
	// 	fmt.Printf("%s: %s, %t\n===\n", key, val, ok)
	// }
	for _, key := range keys {
		diffs, ok := m[key]
		if ok && diffs != nil {
			summary, err := diffs.SummarizeToString(how)
			if err != nil {
				return "", fmt.Errorf("error summarizing %s: %s", diffs.kind, err.Error())
			}
			if summary != "" {
				return summary, nil
			}
		}
	}
	return "", nil
}
