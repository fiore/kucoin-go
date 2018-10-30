# Kucoin REST-API Golang
[![Build Status](https://travis-ci.org/fiore/kucoin-go.svg?branch=master)](https://travis-ci.org/fiore/kucoin-go)
[![GoDoc](https://godoc.org/github.com/fiore/kucoin-go?status.svg)](https://godoc.org/github.com/fiore/kucoin-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/fiore/kucoin-go)](https://goreportcard.com/report/github.com/fiore/kucoin-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/fiore/kucoin-go/blob/master/LICENSE)

Unofficial [Kucoin API](https://kucoinapidocs.docs.apiary.io/) implementation written on Golang.

## Features
- Ready to go solution. Just import the package
- The most needed methods are implemented
- Simple authorization handling
- Pure and stable code
- Built-in Golang performance

## How to use
```bash
go get -u github.com/fiore/kucoin-go
```
```golang
package main

import (
	"github.com/fiore/kucoin-go"
)

func main() {
	// Set your own API key and secret
	k := kucoin.New("API_KEY", "API_SECRET")
	// Call resource
	k.GetCoinBalance("BTC")
}
```
## Checklist
| API Resource                                 | Type | Done |
| -------------------------------------------- | ---- | ---- |
| Tick (symbols)                               | Open | ✔    |
| Get coin info                                | Open | ✔    |
| List coins                                   | Open | ✔    |
| Tick (symbols) for logged user               | Auth | ✔    |
| Get coin deposit address                     | Auth | ✔    |
| Get balance of coin                          | Auth | ✔    |
| Create an order                              | Auth | ✔    |
| Get user info                                | Auth | ✔    |
| List active orders (Both map and array)      | Auth | ✔    |
| List deposit & withdrawal records            | Auth | ✔    |
| List dealt orders (Both Specific and Merged) | Auth | ✔    |
| Order details                                | Auth | ✔    |
| Create withdrawal apply                      | Auth | ✔    |
| Cancel withdrawal                            | Auth | ✔    |
| Cancel orders                                | Auth | ✔    |
| Cancel all orders                            | Auth | ✔    |
| Order books                                  | Auth | ✔    |

## Donate
Your **★Star** will be best donation to my work
