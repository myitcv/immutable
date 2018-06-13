### `myitcv.io/hybridimporter`

`myitcv.io/hybridimporter` is an implementation of [`go/types.ImporterFrom`](https://godoc.org/go/types#ImporterFrom)
that uses non-stale package dependency targets where they exist, else falls back to a source-file based importer.

This is essentially a work-in-progress and will become obsolete when
[`go/packages`](https://github.com/golang/go/issues/14120#issuecomment-383994980) lands. The importer discussed in that
thread will be able to take advantage of the build cache and also be `vgo`-aware.

Currently relies on Go tip as of
[baf399b02e](https://go.googlesource.com/go/+/baf399b02e7a17add068b185d8969a50ca2fb8a0), which is due to be part of Go
1.11.
