package hybridimporter // import "myitcv.io/hybridimporter"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"go/importer"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"

	"myitcv.io/hybridimporter/srcimporter"
)

type pkgInfo struct {
	ImportPath string
	Target     string
	Stale      bool
}

func New(ctxt *build.Context, fset *token.FileSet, path, dir string) (*srcimporter.Importer, error) {
	cmd := exec.Command("go", "list", "-deps", "-test", "-json", path)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to start list for %v in %v: %v\n%v", path, dir, err, string(out))
	}

	lookups := make(map[string]io.ReadCloser)

	dec := json.NewDecoder(bytes.NewReader(out))

	for {
		var p pkgInfo
		err := dec.Decode(&p)
		if err != nil {
			if io.EOF == err {
				break
			}
			return nil, fmt.Errorf("failed to parse list for %v in %v: %v", path, dir, err)
		}
		if p.ImportPath == "unsafe" || p.Stale {
			continue
		}
		fi, err := os.Open(p.Target)
		if err != nil {
			return nil, fmt.Errorf("failed to open %v: %v", p.Target, err)
		}
		lookups[p.ImportPath] = fi
	}

	lookup := func(path string) (io.ReadCloser, error) {
		rc, ok := lookups[path]
		if !ok {
			return nil, fmt.Errorf("failed to resolve %v", path)
		}

		return rc, nil
	}

	gc := importer.For("gc", lookup)

	tpkgs := make(map[string]*types.Package)

	for path := range lookups {
		p, err := gc.Import(path)
		if err != nil {
			return nil, fmt.Errorf("failed to gc import %v: %v", path, err)
		}
		tpkgs[path] = p
	}

	return srcimporter.New(ctxt, fset, tpkgs), nil
}
