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

	helper := func(t *testing.T, src, typename, exp string) {
		t.Helper()
		n, s, err := Parse(bytes.NewBufferString(src), typename)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGenerator(n, s)
		g.AppendPackage()
		if err := g.AppendCheckFunction(); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		format.Node(&buf, token.NewFileSet(), g.f)
		if act := buf.String(); act != exp {
			t.Error(diff.LineDiff(exp, act))
		}
	}

	t.Run("return no value", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound()
}
`
		exp := `package main

func IsUserNotFound(err error) bool {
	if e, ok := err.(userNotFound); ok {
		return true
	}
	return false
}
`
		helper(t, src, "userNotFound", exp)
	})

	t.Run("return int value", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound() (id int64)
}
`

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
		helper(t, src, "userNotFound", exp)
	})

	t.Run("return multi value value", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound() (id int64, name string)
}
`

		exp := `package main

func IsUserNotFound(err error) (bool, int64, string) {
	var id int64
	var name string
	if e, ok := err.(userNotFound); ok {
		id, name = e.UserNotFound()
		return true, id, name
	}
	return false, id, name
}
`
		helper(t, src, "userNotFound", exp)
	})
}

func TestGenerator_AppendErrorImplementation(t *testing.T) {

	helper := func(t *testing.T, src, typename, exp string) {
		t.Helper()
		n, s, err := Parse(bytes.NewBufferString(src), typename)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGenerator(n, s)
		g.AppendPackage()
		if err := g.AppendErrorImplementation(); err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		format.Node(&buf, token.NewFileSet(), g.f)
		if act := buf.String(); act != exp {
			t.Error(diff.LineDiff(exp, act))
		}
	}

	t.Run("return no value", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound()
}
`
		exp := `package main

type UserNotFound struct {
}

func (e *UserNotFound) UserNotFound() {
}
func (e *UserNotFound) Error() string {
	return fmt.Sprint("userNotFound")
}
`
		helper(t, src, "userNotFound", exp)
	})
}
