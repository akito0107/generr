package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fmt"

	"github.com/akito0107/generr"
)

var typename = flag.String("type", "", "pass error type name (required)")
var dryrun = flag.Bool("dryrun", false, "dryrun (default=false)")
var genImpl = flag.Bool("impl", false, "generate error implementation")

func main() {
	flag.Parse()
	if *typename == "" {
		log.Fatal("must be passed type name")
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
		ok, err := generate(f)
		if err != nil {
			log.Fatal(err)
		}
		if ok {
			return
		}
	}

	log.Fatalf("typename %s notfound", *typename)
}

func generate(filename string) (bool, error) {
	r, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer r.Close()
	pkgName, ts, err := generr.Parse(r, *typename)
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
	if *genImpl {
		if err := g.AppendErrorImplementation(); err != nil {
			return false, err
		}
	}

	if *dryrun {
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
