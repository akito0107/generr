package e2e

// go:generate generr -t userNotFound -i -u
type userNotFound interface {
	UserNotFound() (id int64)
}

//go:generate generr -t notFound -i
type notFound interface {
	NotFound()
}

// go:generate generr -t emailNotFound -i -m "email %s is not found" -u
type emailNotFound interface {
	EmailNotFound() (email string)
}

//go:generate generr -t subpackageNotFound -i -it subpackageNotFound -o ./subpackage
type subpackageNotFound interface {
	SubpackageNotFound() (info string)
}
