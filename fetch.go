package main

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

type Country struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Capital    []string `json:"capital"`
	Population int      `json:"population"`
	Region     string   `json:"region"`
	Flags      struct {
		PNG string `json:"png"`
	} `json:"flags"`
}

// Get all data from REST countries and put it into table
func FetchAllCountriesFromAPI(ctx *fasthttp.RequestCtx) ([]Country, error) {
	url := "https://restcountries.com/v3.1/all"
	status, body, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Failed to fetch country data"))
		return nil, err
	}

	var countries []Country
	err = json.Unmarshal(body, &countries)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Failed to fetch country data"))
		return nil, err
	}

	return countries, nil
}

// Get data from table for specific country
func FetchCountryFromClickHouse(countryName string) (*Country, error) {
	var country Country
	query := "SELECT name, capital, population, region, flag_url FROM countries WHERE lower(name) = lower(?)"
	err := db.QueryRow(query, countryName).Scan(&country.Name.Common, &country.Capital, &country.Population, &country.Region, &country.Flags.PNG)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &country, nil
}

// Get data for all countries from table
func FetchAllCountriesFromClickHouse() ([]Country, error) {
	query := "SELECT name, capital, population, region, flag_url FROM countries ORDER BY name"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Query error")
		return nil, err
	}
	defer rows.Close()

	var countries []Country
	for rows.Next() {
		var country Country
		err := rows.Scan(&country.Name.Common, &country.Capital, &country.Population, &country.Region, &country.Flags.PNG)
		if err != nil {
			log.Printf("rows.Scan() error: %v", err)
			return nil, err
		}
		countries = append(countries, country)
	}

	return countries, nil
}

// calculate population stats via ClickHouse query
func FetchPopulationStats() (int, float64, string, error) {
	var totalPopulation int
	var averagePopulation float64
	err := db.QueryRow("SELECT SUM(population), AVG(population) FROM countries").Scan(&totalPopulation, &averagePopulation)
	if err != nil {
		return 0, 0, "", err
	}

	var mostPopulatedCountry string
	err = db.QueryRow("SELECT name FROM countries ORDER BY population DESC LIMIT 1").Scan(&mostPopulatedCountry)
	if err != nil {
		return 0, 0, "", err
	}

	return totalPopulation, averagePopulation, mostPopulatedCountry, nil
}
