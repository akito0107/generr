package e2e

import (
	"fmt"
)

//go:generate -Type=userNotFound
type userNotFound interface {
	UserNotFound() (id int64)
}

func IsUserNotFound(err error) (bool, int64) {
	var id int64
	if e, ok := err.(userNotFound); ok {
		return true, e.UserNotFound()
	}
	return false, id
}

type UserNotFound struct {
	Id int64
}

func (u *UserNotFound) UserNotFound() (id int64) {
	return id
}

func (u *UserNotFound) Error() string {
	return fmt.Sprintf("userNotFound with Id: %v", u.Id)
}
