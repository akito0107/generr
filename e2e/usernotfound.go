package e2e

//go:generate generr -type=userNotFound -impl
type userNotFound interface {
	UserNotFound() (id int64)
}

//go:generate generr -type=notFound -impl
type notFound interface {
	NotFound()
}
