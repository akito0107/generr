# generr

generate error utilities for golang.

Inspired by [BANZAI CLOUD's Error Handling Practices in Go](https://banzaicloud.com/blog/error-handling-go/),
this command generate custom error utilities from simple `interface` declaration.

## Getting Started

### Prerequisites
- Go 1.11+
- make

### Installing
```
$ go get -u github.com/akito0107/generr/cmd/generr
```

### How to use
1. You declare `interface` which must have a single function and *named* return values (if needs).
```go
type userNotFound interface {
	UserNotFound() (id int64)
}
```

2. Generate implementation with passing the type name.
```sh
$ generr -type=userNotFound
```

3. Then, you can get implementation file which named `userNotFound_impl.go`.
This file contains the function which identifies whether given error is the one that we are defined before.

```go
func IsUserNotFound(err error) (bool, int64) {
	var id int64
	if e, ok := err.(userNotFound); ok {
		id = e.UserNotFound()
		return true, id
	}
	return false, id
}
```

4. You can also generate struct which implements `error` and the `interface` with `-impl` option.

```sh
$ generr -type=userNotFound -impl
```

```go
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
```

5. You can also use `go generate`

```go
//go:generate generr -type=notFound -impl
type userNotFound interface {
	UserNotFound() (id int64)
}
```

```sh
$ go generate // you can get same results.
```

## Options
- type (required) target interface type name.
- impl (optional, default=false) generate implementation.
- dryrun  (optional, default=false)

## License
This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
