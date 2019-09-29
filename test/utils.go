package test

import "github.com/brianvoe/gofakeit"

func RandomNumber() int {
	return gofakeit.Number(0, 100)
}

func RandomName() string {
	return gofakeit.Name()
}

func RandomUsername() string {
	return gofakeit.Username()
}

func RandomPassword() string {
	return gofakeit.Password(true, true, true, true, false, 20)
}

func RandomDomain() string {
	return gofakeit.DomainName()
}
