# generr

[![CircleCI](https://circleci.com/gh/akito0107/generr.svg?style=svg)](https://circleci.com/gh/akito0107/generr)
[![Maintainability](https://api.codeclimate.com/v1/badges/5acb46b675867eaa697e/maintainability)](https://codeclimate.com/github/akito0107/generr/maintainability)[![Test Coverage](https://api.codeclimate.com/v1/badges/5acb46b675867eaa697e/test_coverage)](https://codeclimate.com/github/akito0107/generr/test_coverage)

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
$ generr -t userNotFound
```

3. Then, you can get implementation file which named `userNotFound_check.go`.
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

4. You can also generate struct which implements `error` and the `interface` with `-i` option.

```sh
$ generr -t userNotFound -i -it userNotFound -o ../otherpackage
```
You get generated file named `userNotFound_impl.go`.

```go
package otherpackage

type userNotFound struct {
	Id int64
}

func (e *userNotFound) UserNotFound() int64 {
	return e.Id
}
func (e *userNotFound) Error() string {
	return fmt.Sprintf("userNotFound Id: %v", e.Id)
}
```

You can pass struct name with `-it` option (default case is capitalized name given with `-t`).
`-o` option (default location is current directory and package).

5. You can unify `check` and `impl` files with `-u` option.
```sh
$ generr -t userNotFound -i -u
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

6. You can also use `go generate`

```go
//go:generate generr -t notFound -i -u
type userNotFound interface {
	UserNotFound() (id int64)
}
```

```sh
$ go generate // you can get same results.
```

6. You can pass custom error message with `-m` flag.
```go
//go:generate generr -t emailNotFound -m "email %s not found" -i -u
type emailNotFound interface{
    EmailNotFound() (email string)	
}
```

and run `go generate`

```sh
$ go generate 
```

then, you get custom error utilities with passed error message.

```go
// Code generated by "generr"; DO NOT EDIT.
package e2e

import "fmt"

func IsEmailNotFound(err error) (bool, string) {
	var email string
	if e, ok := err.(emailNotFound); ok {
		email = e.EmailNotFound()
		return true, email
	}
	return false, email
}

type EmailNotFound struct {
	Email string
}

func (e *EmailNotFound) EmailNotFound() string {
	return e.Email
}
func (e *EmailNotFound) Error() string {
	return fmt.Sprintf("email %s is not found", e.Email)
}
```


## Options
```sh
$ generr -h
NAME:
   generr - generate custom error from interface

USAGE:
   generr [OPTIONS]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --type value, -t value                        error interface name (required)
   --dryrun                                      dryrun (default=false)
   --implementation, -i                          generate error implementation (default=false)
   --unify, -u                                   (only affects with --implementation option) unify implementation with checking function (default=false)
   --implementation-output-path value, -o value  (only affects with --implementation option) implementation output path (default=current directory)
   --implementation-type value, --it value       (only affects with --implementation option) implementation type name (default=capitalized given type name)
   --message value, -m value                     custom error message (optional)
   --cause, -c                                   append cause check (default=false)
   --help, -h                                    show help
   --version, -v                                 print the version
```

## License
This project is licensed under the Apache License 2.0 License - see the [LICENSE](LICENSE) file for details
