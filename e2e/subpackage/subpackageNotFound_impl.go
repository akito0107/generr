// Code generated by "generr"; DO NOT EDIT.
package e2e

import "fmt"

type subpackageNotFound struct {
	Info string
}

func (e *subpackageNotFound) SubpackageNotFound() string {
	return e.Info
}
func (e *subpackageNotFound) Error() string {
	return fmt.Sprintf("subpackageNotFound Info: %v", e.Info)
}
