# Maskada API

[![Build Status](https://travis-ci.org/gritt/maskada.svg)](https://travis-ci.org/gritt/maskada)
[![Go Report Card](https://goreportcard.com/badge/github.com/gritt/maskada)](https://goreportcard.com/report/github.com/gritt/maskada)
[![Coverage Status](https://codecov.io/gh/gritt/maskada/branch/master/graph/badge.svg)](https://codecov.io/gh/gritt/maskada)

### About

Maskada (from portuguese *mascada*, means money) is a simple money manager api developed in Go. 
It aims to perform very simple operations, to make it easier track your finances.

###### ROADMAP

- ✓︎ Create transactions
- ✓︎ List transactions
- ✓︎︎ Create transactions with category
- ✘ Manage transaction status like: delete/pending/done
- ✘ Create recurring transactions

To make it simple to calculate, all transactions will belong to a type:

[`core.go`](./core/core.go)
```
// Debit is a transaction which is subtracted.
Debit = 1

// Credit is a transaction which is subtracted the next month.
Credit = 2

// Income is a transaction which is summed.
Income = 3
```

All the calculations happen in the client side.

### Architecture

The backend its inspired by the [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) design, 
the core business logic its encapsulated in the `core` package through *Use Cases*, 
the Use Cases communicate with the external world by well-defined interfaces, 
that are implemented by the clients.

### Documentation

To learn more, checkout the Wiki:

- [Setup](./wiki/Setup.md)
- [API Contract](./wiki/API.md)

### Contributing

> Currently, this project is mostly for studying purposes, and it's not hosted anywhere.
> If you're interested in learning golang, ddd, clean architecture and tests, feel free to fork it.
