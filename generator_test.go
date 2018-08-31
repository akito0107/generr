package generr

import (
	"testing"

	"bytes"

	"go/format"
	"go/token"

	"github.com/andreyvit/diff"
)

func TestGenerator_AppendPackage(t *testing.T) {
	g := NewGenerator("main", nil)
	g.AppendPackage()

	exp := "package main\n"
	var buf bytes.Buffer
	format.Node(&buf, token.NewFileSet(), g.f)
	if act := buf.String(); act != exp {
		t.Error(diff.LineDiff(exp, act))
	}
}

func TestGenerator_AppendCheckFunction(t *testing.T) {
	t.Run("return no value", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound()
}
`
		n, s, err := Parse(bytes.NewBufferString(src), "userNotFound")
		if err != nil {
			t.Fatal(err)
		}
		g := NewGenerator(n, s)
		g.AppendPackage()
		if err := g.AppendCheckFunction(); err != nil {
			t.Fatal(err)
		}

		exp := `package main

func IsUserNotFound(err error) bool {
	if e, ok := err.(userNotFound); ok {
		return true
	}
	return false
}
`
		var buf bytes.Buffer
		format.Node(&buf, token.NewFileSet(), g.f)
		if act := buf.String(); act != exp {
			t.Error(diff.LineDiff(exp, act))
		}
	})

	t.Run("return int value", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound() (id int64)
}
`
		n, s, err := Parse(bytes.NewBufferString(src), "userNotFound")
		if err != nil {
			t.Fatal(err)
		}
		g := NewGenerator(n, s)
		g.AppendPackage()
		if err := g.AppendCheckFunction(); err != nil {
			t.Fatal(err)
		}

		exp := `package main

func IsUserNotFound(err error) (bool, int64) {
	var id int64
	if e, ok := err.(userNotFound); ok {
		id = e.UserNotFound()
		return true, id
	}
	return false, id
}
`
		var buf bytes.Buffer
		format.Node(&buf, token.NewFileSet(), g.f)
		if act := buf.String(); act != exp {
			t.Error(diff.LineDiff(exp, act))
		}
	})
}
