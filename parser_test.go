package generr

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	src := `package main

type userNotFound interface {
	UserNotFound() bool
}
`
	t.Run("can parse package name", func(t *testing.T) {
		s, _, err := Parse(bytes.NewBufferString(src), "userNotFound")
		if err != nil {
			t.Fatal(err)
		}
		if s != "main" {
			t.Errorf("must be main but %s", s)
		}
	})

	t.Run("can parse interface with given type name", func(t *testing.T) {
		_, node, err := Parse(bytes.NewBufferString(src), "userNotFound")
		if err != nil {
			t.Fatal(err)
		}
		if node.Name.Name != "userNotFound" {
			t.Errorf("type name must be userNotFound but %s", node.Name.Name)
		}
	})

	t.Run("return err if given typename is not found", func(t *testing.T) {
		_, _, err := Parse(bytes.NewBufferString(src), "xxx")
		if err == nil {
			t.Fatal(err)
		}
	})
}
