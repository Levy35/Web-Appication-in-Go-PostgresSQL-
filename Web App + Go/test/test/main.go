package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Third party package
)

func setUpDB() (*sql.DB, error) {
	// Provide the credentials for our database
	const (
		host     = "localhost"
		port     = 5432
		user     = "music"
		password = "#swordfish#"
		dbname   = "music"
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Establish a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Test our connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Dependencies (things/variables)
// Dependency Injection (passing)
type application struct {
	db *sql.DB
}

func main() {
	var db, err = setUpDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Always do this before exiting
	app := &application{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/info", app.createMusicianInfo)
	mux.HandleFunc("/info-add", app.createMusician)
	mux.HandleFunc("/display", app.displayMusician)
	log.Println("Starting server on port :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
