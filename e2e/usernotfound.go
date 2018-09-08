package e2e

//go:generate generr -t userNotFound -i
type userNotFound interface {
	UserNotFound() (id int64)
}

//go:generate generr -t notFound -i
type notFound interface {
	NotFound()
}

//go:generate generr -t emailNotFound -i -m "email %s is not found"
type emailNotFound interface {
	EmailNotFound() (email string)
}
