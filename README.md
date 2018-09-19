# echo-prometheus
A prometheus exporter for echo

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/0neSe7en/echo-prometheus)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/0neSe7en/echo-prometheus/master/LICENSE)

## Install

`$ go get github.com/0neSe7en/echo-prometheus`

## Usage

```go
package main

import (
	"github.com/0neSe7en/echo-prometheus"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	e := echo.New()

	e.Use(echoprometheus.NewMetric())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
```
