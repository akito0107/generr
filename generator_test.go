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
		if err := g.AppendCheckFunction(false); err != nil {
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
	if _, ok := err.(userNotFound); ok {
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

func TestGenerator_AppendCheckFunctionWithCause(t *testing.T) {

	helper := func(t *testing.T, src, typename, exp string) {
		t.Helper()
		n, s, err := Parse(bytes.NewBufferString(src), typename)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGenerator(n, s)
		g.AppendPackage()
		g.AppendPkgErrorImportSpec()
		if err := g.AppendCheckFunction(true); err != nil {
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
	err = errors.Cause(err)
	if _, ok := err.(userNotFound); ok {
		return true
	}
	return false
}
`
		helper(t, src, "userNotFound", exp)
	})

}

func TestGenerator_AppendErrorImplementation(t *testing.T) {

	helper := func(t *testing.T, src, typename, exp string, msg string) {
		t.Helper()
		n, s, err := Parse(bytes.NewBufferString(src), typename)
		if err != nil {
			t.Fatal(err)
		}
		g := NewGenerator(n, s)
		g.AppendPackage()
		if err := g.AppendErrorImplementation("", msg); err != nil {
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
	return
}
func (e *UserNotFound) Error() string {
	return fmt.Sprint("userNotFound")
}
`
		helper(t, src, "userNotFound", exp, "")
	})

	t.Run("return int value", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound() (id int64)
}
`

		exp := `package main

type UserNotFound struct {
	Id int64
}

func (e *UserNotFound) UserNotFound() int64 {
	return e.Id
}
func (e *UserNotFound) Error() string {
	return fmt.Sprintf("userNotFound Id: %v", e.Id)
}
`
		helper(t, src, "userNotFound", exp, "")
	})

	t.Run("return multiple value", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound() (id int64, name string)
}
`

		exp := `package main

type UserNotFound struct {
	Id   int64
	Name string
}

func (e *UserNotFound) UserNotFound() (int64, string) {
	return e.Id, e.Name
}
func (e *UserNotFound) Error() string {
	return fmt.Sprintf("userNotFound Id: %v Name: %v", e.Id, e.Name)
}
`
		helper(t, src, "userNotFound", exp, "")
	})

	t.Run("return multiple value with custom message", func(t *testing.T) {
		src := `package main

type userNotFound interface {
	UserNotFound() (id int64, name string)
}
`

		exp := `package main

type UserNotFound struct {
	Id   int64
	Name string
}

func (e *UserNotFound) UserNotFound() (int64, string) {
	return e.Id, e.Name
}
func (e *UserNotFound) Error() string {
	return fmt.Sprintf("custom message with %d and %s", e.Id, e.Name)
}
`
		msg := "custom message with %d and %s"
		helper(t, src, "userNotFound", exp, msg)
	})
}
