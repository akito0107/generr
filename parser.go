package generr

import (
	"go/parser"
	"go/token"

	"go/ast"

	"io"

	"io/ioutil"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
)

func Parse(r io.Reader, tp string) (string, *ast.TypeSpec, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil, errors.Wrap(err, "parse file readAll failed")
	}
	f, err := parser.ParseFile(token.NewFileSet(), "", string(src), parser.ParseComments)
	if err != nil {
		return "", nil, errors.Wrap(err, "parse file with filename is failed")
	}

	var pkgName string
	var ts *ast.TypeSpec
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			pkgName = x.Name.Name
		case *ast.TypeSpec:
			pp.Println(x)
			if _, ok := x.Type.(*ast.InterfaceType); !ok {
				return true
			}
			if x.Name.Name == tp {
				ts = x
				return false
			}
		default:
			return true
		}
		return true
	})

	if ts == nil {
		return "", nil, errors.Errorf("typename: %s not found", tp)
	}

	return pkgName, ts, nil
}
