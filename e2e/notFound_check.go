// Code generated by "generr"; DO NOT EDIT.
package e2e

func IsNotFound(err error) bool {
	if _, ok := err.(notFound); ok {
		return true
	}
	return false
}