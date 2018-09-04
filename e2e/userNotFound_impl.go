// Code generated by "generr"; DO NOT EDIT.
package e2e

import "fmt"

func IsUserNotFound(err error) (bool, int64) {
	var id int64
	if e, ok := err.(userNotFound); ok {
		id = e.UserNotFound()
		return true, id
	}
	return false, id
}

type UserNotFound struct {
	Id int64
}

func (e *UserNotFound) UserNotFound() int64 {
	return e.Id
}
func (e *UserNotFound) Error() string {
	return fmt.Sprintf("userNotFound Id: %v", e.Id)
}