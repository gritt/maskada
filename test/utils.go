package test

import "github.com/brianvoe/gofakeit"

// RandomNumber returns a fake number.
func RandomNumber() int {
	return gofakeit.Number(0, 100)
}

// RandomName returns fake name.
func RandomName() string {
	return gofakeit.Name()
}

// RandomUsername returns a fake username.
func RandomUsername() string {
	return gofakeit.Username()
}

// RandomPassword returns fake password.
func RandomPassword() string {
	return gofakeit.Password(true, true, true, true, false, 20)
}

// RandomDomain returns a fake domain name.
func RandomDomain() string {
	return gofakeit.DomainName()
}
