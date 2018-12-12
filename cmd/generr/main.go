package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/akito0107/generr"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "generr"
	app.Usage = "generate custom error from interface"
	app.UsageText = "generr [OPTIONS]"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "type, t",
			Usage: "error interface name (required)",
		},
		cli.BoolFlag{
			Name:  "dryrun",
			Usage: "dryrun (default=false)",
		},
		cli.BoolFlag{
			Name:  "implementation, i",
			Usage: "generate error implementation (default=false)",
		},
		cli.BoolFlag{
			Name:  "unify, u",
			Usage: "(only affects with --implementation option) unify implementation with checking function (default=false)",
		},
		cli.StringFlag{
			Name:  "implementation-output-path, o",
			Usage: "(only affects with --implementation option) implementation output path (default=current directory)",
		},
		cli.StringFlag{
			Name:  "implementation-type, it",
			Usage: "(only affects with --implementation option) implementation type name (default=capitalized given type name)",
		},
		cli.StringFlag{
			Name:  "message, m",
			Usage: "custom error message (optional)",
		},
		cli.BoolFlag{
			Name:  "cause, c",
			Usage: "append cause check (default=false)",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	typename := ctx.String("type")
	message := ctx.String("message")
	dryrun := ctx.Bool("dryrun")
	cause := ctx.Bool("cause")
	impl := ctx.Bool("implementation")
	unify := ctx.Bool("unify")
	outpath := ctx.String("implementation-output-path")
	it := ctx.String("implementation-type")

	if typename == "" {
		return errors.New("type is required")
	}

	var filenames []string
	fileinfos, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range fileinfos {
		if strings.HasSuffix(f.Name(), ".go") && !strings.HasSuffix(f.Name(), "_test.go") {
			filenames = append(filenames, f.Name())
		}
	}

	for _, f := range filenames {
		ok, err := generate(f, typename, message, outpath, it, dryrun, impl, cause, unify)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	return errors.Errorf("typename %s is not found", typename)
}

func generate(filename, typename, message, outpath, it string, dryrun, impl, cause, unify bool) (bool, error) {
	r, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer r.Close()
	pkgName, ts, err := generr.Parse(r, typename)
	if err != nil && generr.IsTypeNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	// only check function
	if !impl {
		g := generr.NewGenerator(pkgName, ts)
		g.AppendPackage()
		if err := g.AppendCheckFunction(cause); err != nil {
			return false, err
		}
		if err := write(g, fmt.Sprintf("%s_check.go", ts.Name.Name), dryrun); err != nil {
			return false, err
		}
		return true, nil
	}

	if unify {
		g := generr.NewGenerator(pkgName, ts)
		g.AppendPackage()
		if err := g.AppendCheckFunction(cause); err != nil {
			return false, err
		}
		if err := g.AppendErrorImplementation(it, message); err != nil {
			return false, err
		}
		if err := write(g, fmt.Sprintf("%s_impl.go", ts.Name.Name), dryrun); err != nil {
			return false, err
		}
		return true, nil
	} else {
		g := generr.NewGenerator(pkgName, ts)
		g.AppendPackage()
		if err := g.AppendCheckFunction(cause); err != nil {
			return false, err
		}
		if err := write(g, fmt.Sprintf("%s_check.go", ts.Name.Name), dryrun); err != nil {
			return false, err
		}

		outpackage := pkgName
		if outpath == "" {
			_, outpackage = filepath.Split(outpath)
		}
		g = generr.NewGenerator(outpackage, ts)
		g.AppendPackage()
		if err := g.AppendErrorImplementation(it, message); err != nil {
			return false, err
		}
		p := filepath.Join(outpath, fmt.Sprintf("%s_impl.go", ts.Name.Name))
		if err := write(g, p, dryrun); err != nil {
			return false, err
		}

		return true, nil
	}
}

func write(g *generr.Generator, fpath string, dryrun bool) error {
	if dryrun {
		return g.Out(os.Stdout)
	}
	if _, err := os.Stat(fpath); !os.IsNotExist(err) {
		os.Remove(fpath)
	}
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	return g.Out(f)
}
