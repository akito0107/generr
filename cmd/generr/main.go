package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		cli.StringFlag{
			Name:  "message, m",
			Usage: "custom error message (optional)",
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
	impl := ctx.Bool("implementation")
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
		ok, err := generate(f, typename, message, dryrun, impl)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	return errors.Errorf("typename %s is not found", typename)
}

func generate(filename, typename, message string, dryrun, impl bool) (bool, error) {
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
	g := generr.NewGenerator(pkgName, ts)
	g.AppendPackage()
	if err := g.AppendCheckFunction(); err != nil {
		return false, err
	}
	if impl {
		if err := g.AppendErrorImplementation(message); err != nil {
			return false, err
		}
	}

	if dryrun {
		g.Out(os.Stdout)
	} else {
		filename := fmt.Sprintf("%s_impl.go", ts.Name.Name)

		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			os.Remove(filename)
		}
		f, err := os.Create(filename)
		if err != nil {
			return false, err
		}
		defer f.Close()
		if err := g.Out(f); err != nil {
			return false, err
		}
	}

	return true, nil
}
