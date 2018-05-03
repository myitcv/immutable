### `gg`

`gg` is a dependency-aware wrapper for go generate. Currently requires [`tip`](https://github.com/golang/go), i.e. this
will not work with Go `< 1.11`.

_Very much work in progress._

```bash
go get -u myitcv.io/gg
```

`gg` was born out of the following scenario:

* it's a good idea to clean all generated files as part of a CI build and regenerate; therefore we need a simple,
  reliable means to re-run `go generate` (or similar) on an entire repo of packages
* some `go generate` programs will generate code that itself contains `go generate` directives; this requires `go
  generate` to be called multiple times before a "fixed point" is reached

Whilst it's possible to achieve this on a per-project basis by writing a relatively simple program to wrap things up,
there is some merit in writing a tool to wrap `go generate`:

* the tool can be reused by others
* existing `go generate` programs can be re-used with zero effort

### Usage

_To follow. For now see [`gg_test.go`](https://github.com/myitcv/gg/blob/add_test_framework/gg_test.go)._

### TODO

* Tidy up the code, docs, more tests
* Add support for/switch to using `go run pkg` when it is sped up/its results are cached. [Issue
  raised](https://github.com/golang/go/issues/25416).
* Make config file discovery relative to the directive-containing-file
* Ensure we properly handle when the a `go generate` iteration "grows" a new dependency that is not yet satisfied
* Make `gg` cache-based, a la the `go` command itself
* Use `.DepOnly` field from `go list` [when it lands](https://go-review.googlesource.com/c/go/+/112755)

### Credit

* Russ Cox's [`gt`](https://github.com/rsc/gt)
* `go generate` source code in [the main Go repo](https://github.com/golang/go/tree/master/src/cmd/go)
