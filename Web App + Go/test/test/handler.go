package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

// A struct to hold a quote
type Musician struct {
	Music_id      int
	Full_name     string
	Album         string
	Genre         string
	Date_released time.Time
	Artist        string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/index.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) createMusicianInfo(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/music_artist.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

}

func (app *application) createMusician(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/info", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	full_name := r.PostForm.Get("full_name")
	album := r.PostForm.Get("album")
	genre := r.PostForm.Get("genre")
	date_released := r.PostForm.Get("date_released")
	artist := r.PostForm.Get("artist")

	//check the web form fields to validity
	errors := make(map[string]string)
	//check each field
	if strings.TrimSpace(genre) == "" {
		errors["genre"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(genre) > 20 {
		errors["genre"] = "This field cannot is to large(maximum is 20 characters)"
	}

	if strings.TrimSpace(full_name) == "" {
		errors["full_name"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(full_name) > 25 {
		errors["full_name"] = "This field cannot is to large(maximum is 25 characters)"
	}

	if strings.TrimSpace(album) == "" {
		errors["album"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(album) > 50 {
		errors["album"] = "This field cannot is to large(maximum is 50 characters)"
	}

	if strings.TrimSpace(date_released) == "" {
		errors["date_released"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(date_released) > 50 {
		errors["date_released"] = "This field cannot is to large(maximum is 50 characters)"
	}

	if strings.TrimSpace(artist) == "" {
		errors["artist"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(artist) > 50 {
		errors["artist"] = "This field cannot is to large(maximum is 50 characters)"
	}
	//check if there are any errrors in the map
	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	s := `
	INSERT INTO musician(full_name, album, date_released, genre, artist)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err = app.db.Exec(s, full_name, album, date_released, genre, artist)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func (app *application) displayMusician(w http.ResponseWriter, r *http.Request) {
	//sql statement
	musicAnalys := `
	SELECT *
	FROM musician
	LIMIT 6
	`

	rows, err := app.db.Query(musicAnalys)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	defer rows.Close()

	var music []Musician
	for rows.Next() {
		var m Musician
		err = rows.Scan(&m.Music_id, &m.Full_name, &m.Album, &m.Genre, &m.Date_released, &m.Artist)

		if err != nil {
			log.Println(err.Error())
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		music = append(music, m)
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	//display the quotes using a template
	ts, err := template.ParseFiles("./ui/html/input.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, music)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

}
