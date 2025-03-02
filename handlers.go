package main

import (
	"fmt"
	"html/template"
	"log"

	"github.com/valyala/fasthttp"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Homepage handler
func HomeHandler(ctx *fasthttp.RequestCtx) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, nil)
}

// Display data for specific country
func SearchCountryHandler(ctx *fasthttp.RequestCtx) {
	countryName := string(ctx.FormValue("country"))

	exists, err := CountryExists(countryName)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Database error"))
		return
	}

	var country Country
	if exists {
		log.Printf("Fetching %s from ClickHouse", countryName)
		countryData, err := FetchCountryFromClickHouse(countryName)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBody([]byte("Error fetching from database"))
			return
		}
		country = *countryData
	} else {
		s := fmt.Sprintf("%s doesn't exist", countryName)
		log.Print(s)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(s))
		return
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("text/html")
	log.Printf("Rendering template with country: %+v", country)
	tmpl.Execute(ctx, map[string]interface{}{"Country": country})
}

// Fills database with data from REST countries
func FillDatabaseHandler(ctx *fasthttp.RequestCtx) {
	countries, err := FetchAllCountriesFromAPI(ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Failed to fetch country data from API"))
		return
	}

	err = SaveAllToClickHouse(countries)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Failed to save countries to ClickHouse"))
		return
	}

	ctx.SetContentType("text/plain")
	ctx.SetBody([]byte("All countries have been added to ClickHouse!"))
}

// Displays all countries' data and population stats
func ListCountriesHandler(ctx *fasthttp.RequestCtx) {
	countries, err := FetchAllCountriesFromClickHouse()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Database error"))
		return
	}

	tmpl, err := template.ParseFiles("templates/list.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	totalPopulation, averagePopulation, mostPopulatedCountry, err := FetchPopulationStats()
	if err != nil {
		log.Printf("Error calculating stats: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	frenchPrinter := message.NewPrinter(language.French)
	data := map[string]interface{}{
		"Countries":            countries,
		"TotalPopulation":      frenchPrinter.Sprintf("%d", totalPopulation),
		"AveragePopulation":    frenchPrinter.Sprintf("%d", int64(averagePopulation)),
		"MostPopulatedCountry": mostPopulatedCountry,
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, data)
}
