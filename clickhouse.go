package main

import (
	"database/sql"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"
)

var db *sql.DB

// Connect to db and create table if it does not already exist
func InitClickHouse() {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "myuser",
			Password: "mypassword",
		},
	})

	if err := conn.Ping(); err != nil {
		log.Fatal("ClickHouse connection failed:", err)
	}
	db = conn

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS countries (
    		name String,
    		capital Array(String),
    		population UInt64,
    		region String,
    		flag_url String
		) ENGINE = MergeTree ORDER BY name;
	`)

	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

// Check if the country already exists in table
func CountryExists(countryName string) (bool, error) {
	query := "SELECT COUNT(*) FROM countries WHERE lower(name) = lower(?)"
	var count int

	err := db.QueryRow(query, countryName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Save all countries' data if they does not already exist
func SaveAllToClickHouse(countries []Country) error {
	query := "INSERT INTO countries (name, capital, population, region, flag_url) VALUES (?, ?, ?, ?, ?)"
	for _, country := range countries {
		exists, err := CountryExists(country.Name.Common)
		if err != nil {
			return err
		}

		if !exists {
			_, err = db.Exec(query, country.Name.Common, country.Capital, country.Population, country.Region, country.Flags.PNG)
			if err != nil {
				return err
			}
			log.Printf("Inserted %s into ClickHouse", country.Name.Common)
		} else {
			log.Printf("Country %s already exists in ClickHouse", country.Name.Common)
		}
	}
	log.Println("All countries inserted into ClickHouse successfully!")
	return nil
}
