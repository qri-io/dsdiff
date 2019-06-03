<a name="v0.1.1"></a>
# [v0.1.1](https://github.com/qri-io/dsdiff/compare/v0.1.0...v0.1.1) (2019-06-03)

We had a cross-dependency issue with github.com/qri-io/jsonschema. This release should fix that.


# 0.1.0 (2019-05-30)

`dsdiff` is a utility for diffing datasets, currently a basic placeholder

This is the first proper release of `dsdiff`. In preparation for go 1.13, in which `go.mod` files and go modules are the primary way to handle go dependencies, we are going to do an official release of all our modules. This will be version v0.1.0 of `dsdiff`.

### Bug Fixes

* **DatasetDiff:** fixes false "no changes detected" returns ([8ed0e4c](https://github.com/qri-io/dsdiff/commit/8ed0e4c))
* **Differ:** removed early bail that was preventing structure changes from being recognized ([#9](https://github.com/qri-io/dsdiff/issues/9)) ([fd9745b](https://github.com/qri-io/dsdiff/commit/fd9745b))
* **DiffMeta:** force metat to marshal to objects for comparison ([f89a626](https://github.com/qri-io/dsdiff/commit/f89a626))
* added dataset as dependency to circleci config ([9278a9b](https://github.com/qri-io/dsdiff/commit/9278a9b))
* added missing comment for List() method ([c14bc72](https://github.com/qri-io/dsdiff/commit/c14bc72))
* updated dependencies to include gojsondiff ([5cf9260](https://github.com/qri-io/dsdiff/commit/5cf9260))
* updated diffTransform diffVisConfig to handle nil values ([e1d5c67](https://github.com/qri-io/dsdiff/commit/e1d5c67))


### Features

* **TransformDiff:** Transform path differences are ignored, if non-blank ([1060916](https://github.com/qri-io/dsdiff/commit/1060916))
* allow nil when diffing md,st,tf,vc ([ac2534e](https://github.com/qri-io/dsdiff/commit/ac2534e))
* **MarshalJSON:** marshals slice of SubDiff.diffs to JSON ([effe0ac](https://github.com/qri-io/dsdiff/commit/effe0ac))
* **MarshalJSON:** marshals slice of SubDiff.diffs to JSON ([54dc77e](https://github.com/qri-io/dsdiff/commit/54dc77e))
* added exported List() method to access full list of diffs ([7f38af6](https://github.com/qri-io/dsdiff/commit/7f38af6))
* added generic differ for two slices of bytes ([9019d2c](https://github.com/qri-io/dsdiff/commit/9019d2c))
* initial commit ([329a79f](https://github.com/qri-io/dsdiff/commit/329a79f))
* made `MapDiffsToSTring` more descriptive of changes ([73ed2cf](https://github.com/qri-io/dsdiff/commit/73ed2cf))



