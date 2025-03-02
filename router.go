package main

import (
	"github.com/fasthttp/router"
)

func InitRouter() *router.Router {
	r := router.New()

	// Homepage
	r.GET("/", HomeHandler)

	// Counties's data and population stats
	r.GET("/countries", ListCountriesHandler)

	// Search results page
	r.GET("/country", SearchCountryHandler)

	// Fill database
	r.POST("/filldb", FillDatabaseHandler)

	// Serve static files
	r.ServeFiles("/static/{filepath:*}", "./static")

	return r
}
