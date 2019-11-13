package main

import (
	"database/sql"
	"log"
)

// configDB configures the db object
func configDB() (*sql.DB, error) {

	db, err := sql.Open("mysql", "keid:swamprat@tcp(tinfoilwizard.net:15001)/keid")
	if err != nil {
		log.Fatalf("\n\nmysql: could not get a connection: %v\n\n", err)
	}

	return db, nil
}

