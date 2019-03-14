package dsdiff

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
		{"testdata/orig.json", "testdata/newStructure.json", "simple", "Structure Changed. (3 changes)", ""},
		{"testdata/orig.json", "testdata/newStructure.json", "delta", `{
  "checksum": [
    "@@ -33,7 +33,7 @@\n y9Jc\n-ud9\n+aaa\n",
    0,
    2
  ],
  "entries": [
    33,
    35
  ],
  "schema": {
    "items": {
      "items": {
        "0": {
          "title": [
            "rank",
            "ranking"
          ]
        },
        "1": {
          "title": [
            "probability_of_automation",
            "prob_of_automation"
          ]
        },
        "_t": "a"
      }
    }
  }
}
`, ""},
		{"testdata/orig.json", "testdata/newTitle.json", "listKeys", "Transform: 2 changes\n\t- modified config\n\t- modified syntax", ""},
		{"testdata/orig.json", "testdata/newDescription.json", "plusMinusColor", ` {
[30;41m-  "description": "I am a dataset",[0m
[30;42m+  "description": "I am a new description",[0m
   "qri": "md:0",
   "title": "abc"
 }
`, ""},
		{"testdata/orig.json", "testdata/newVisConfig.json", "plusMinus", ` {
-  "format": "abc",
+  "format": "new thing",
   "qri": "vz:0"
 }
`, ""},
		{"testdata/orig.json", "testdata/newTransform.json", "simple", "Transform Changed. (2 changes)", ""},
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
		//  s, err := MapDiffsToFormattedString(got, dsLeft)
		//  if err != nil {
		//    t.Errorf("not today: %s", err.Error())
		//  }
		//  fmt.Println("--------------------------")
		//  fmt.Print(s)
		//  fmt.Println("--------------------------")
		// }

		if stringDiffs != c.expected {
			// texp := []byte(c.expected)
			tgot := []byte(stringDiffs)
			_ = ioutil.WriteFile("got0.txt", tgot, 0775)
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
		{"testdata/orig.json", "testdata/newStructure.json", "abc", "3 diffs", ""},
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

func BenchmarkDiffDatasets(b *testing.B) {
	var t1 = []byte(`{"body":[["a","b","c","d"],["1","2","3","4"],["e","f","g","h"]],"bodyPath":"/ipfs/QmP2tdkqc4RhSDGv1KSWoJw1pwzNu6HzMcYZaVFkLN9PMc","commit":{"author":{"id":"QmSyDX5LYTiwQi861F5NAwdHrrnd1iRGsoEvCyzQMUyZ4W"},"path":"/ipfs/QmbwJNx88xNknXYewLCVBVJqbZ5oaiffr4WYDoCJAuCZ93","qri":"cm:0","signature":"TUREFCfoKEf5J189c0jdKfleRYsGZm8Q6sm6g6lJctXGDDM8BGdpSVjMltGTmmrtN6qtQJKRail5ceG325Rb8hLYoMe4926gXZNWBlMfD0yBHSjo81LsE25UqVeloU2W19Z1MNOrLTDPDRBoM0g3vyJLykGQ0UPRqpUvXNod0E5ONZOKGrQpByp113h12yiAjsiCBR6sAfIScNpcyjzkiDhBCCbMy9cGfMVK8q7wNCmcC41zguGhvv1biDoE+MEVDc1QPN1dYeEaDsvaRu5jWSv44zhVdC3lZtlT8R9qArk8OQVW798ctQ6NJ5kCiZ3C6Z19VPrptr85oknoNNaYxA==","timestamp":"2019-02-04T14:26:43.158109Z","title":"created dataset"},"name":"test_1","path":"/ipfs/QmeSYBYd3LVsFPRp1jiXgT8q22Md3R7swUzd9yt7MPVUcj/dataset.json","peername":"b5","qri":"ds:0","structure":{"depth":2,"errCount":0,"format":"json","qri":"st:0","schema":{"type":"array"}}}`)
	var t2 = []byte(`{"body":[["a","b","c","d"],["1","2","3","4"],["e","f","g","h"]],"bodyPath":"/ipfs/QmP2tdkqc4RhSDGv1KSWoJw1pwzNu6HzMcYZaVFkLN9PMc","commit":{"author":{"id":"QmSyDX5LYTiwQi861F5NAwdHrrnd1iRGsoEvCyzQMUyZ4W"},"path":"/ipfs/QmVZrXZ2d6DF11BL7QLJ8AYFYaNiLgAWVEshZ3HB5ogZJS","qri":"cm:0","signature":"CppvSyFkaLNIY3lIOGxq7ybA18ZzJbgrF7XrIgrxi7pwKB3RGjriaCqaqTGNMTkdJCATN/qs/Yq4IIbpHlapIiwfzVHFUO8m0a2+wW0DHI+y1HYsRvhg3+LFIGHtm4M+hqcDZg9EbNk8weZI+Q+FPKk6VjPKpGtO+JHV+nEFovFPjS4XMMoyuJ96KiAEeZISuF4dN2CDSV+WC93sMhdPPAQJJZjZX+3cc/fOaghOkuhedXaA0poTVJQ05aAp94DyljEnysuS7I+jfNrsE/6XhtazZnOSYX7e0r1PJwD7OdoZYRH73HnDk+Q9wg6RrpU7EehF39o4UywyNGAI5yJkxg==","timestamp":"2019-02-11T17:50:20.501283Z","title":"forced update"},"name":"test_1","path":"/ipfs/QmaAuKZezio5knAFXU4krPcZfBWHnHDWWKEX32Ne9v6niQ/dataset.json","peername":"b5","previousPath":"/ipfs/QmeSYBYd3LVsFPRp1jiXgT8q22Md3R7swUzd9yt7MPVUcj","qri":"ds:0","structure":{"depth":2,"errCount":0,"format":"json","qri":"st:0","schema":{"type":"array"}}}`)
	ds1 := &dataset.Dataset{}
	ds2 := &dataset.Dataset{}
	if err := ds1.UnmarshalJSON(t1); err != nil {
		b.Fatal(err)
	}
	if err := ds2.UnmarshalJSON(t2); err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		if _, err := DiffDatasets(ds1, ds2, nil); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDiff5Mb(b *testing.B) {
	ds1File, err := ioutil.ReadFile("testdata/airport_codes.json")
	if err != nil {
		b.Fatal(err)
	}
	ds1Body, err := ioutil.ReadFile("testdata/airport_codes_body.json")
	if err != nil {
		b.Fatal(err)
	}
	ds1 := &dataset.Dataset{}
	if err := ds1.UnmarshalJSON(ds1File); err != nil {
		b.Fatal(err)
	}

	ds2File, err := ioutil.ReadFile("testdata/airport_codes_2.json")
	if err != nil {
		b.Fatal(err)
	}
	ds2Body, err := ioutil.ReadFile("testdata/airport_codes_body.json")
	if err != nil {
		b.Fatal(err)
	}
	ds2 := &dataset.Dataset{}
	if err := ds1.UnmarshalJSON(ds2File); err != nil {
		b.Fatal(err)
	}
	// delta, err := DiffDatasets(ds1, ds2, &StructuredDataTuple{
	// 	a: &ds1Body,
	// 	b: &ds2Body,
	// })
	// if err != nil {
	// 	b.Fatal(err)
	// }
	// b.Log(delta)
	for i := 0; i < b.N; i++ {
		DiffDatasets(ds1, ds2, &StructuredDataTuple{
			a: &ds1Body,
			b: &ds2Body,
		})
	}
}
