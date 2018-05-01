### `myitcv.io/hybridimporter`

`myitcv.io/hybridimporter` is an implementation of [`go/types.ImporterFrom`](https://godoc.org/go/types#ImporterFrom)
that uses non-stale package dependency targets where they exist, else falls back to a source-file based importer.

This is essentially a work-in-progress and will become obsolete when
https://github.com/golang/go/issues/14120#issuecomment-383994980 lands. The importer discussed in that thread will be
able to take advantage of the build cache and also be `vgo`-aware.

Currently relies on https://go-review.googlesource.com/c/go/+/107916 which, as of 2018/04/29, is not in a major Go
release.
