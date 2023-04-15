package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/proullon/ramsql/driver"
)

type Movie struct {
	ID          int64   `json:"id"`
	ImdbID      string  `json:"imdbID"`
	Title       string  `json:"titel"`
	Year        int     `json:"year"`
	Rating      float32 `json:"rating"`
	IsSuperHero bool    `json:"isSuperHero"`
}

func getAllMoviesHandler(c echo.Context) error {
	y := c.QueryParam("year")

	if y == "" {
		ms := []Movie{}
		rows, err := db.Query(`SELECT id, imdbID, title, year, rating, isSuperHero FROM iDB`)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			var m Movie

			if err := rows.Scan(&m.ID, &m.ImdbID, &m.Title, &m.Year, &m.Rating, &m.IsSuperHero); err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

			ms = append(ms, m)
		}

		if err := rows.Err(); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, ms)
	} else {
		year, err := strconv.Atoi(y)

		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		ms := []Movie{}
		rows, err := db.Query(`SELECT id, imdbID, title, year, rating, isSuperHero FROM iDB WHERE year=?`, year)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()

		for rows.Next() {
			var m Movie

			if err := rows.Scan(&m.ID, &m.ImdbID, &m.Title, &m.Year, &m.Rating, &m.IsSuperHero); err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

			ms = append(ms, m)
		}

		if err := rows.Err(); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, ms)
	}
}

func getMovieByIdHandler(c echo.Context) error {
	imdbID := c.Param("imdbID")

	row := db.QueryRow(`SELECT id, imdbID, title, year, rating, isSuperHero FROM iDB WHERE imdbID=?`, imdbID)

	m := Movie{}
	err := row.Scan(&m.ID, &m.ImdbID, &m.Title, &m.Year, &m.Rating, &m.IsSuperHero)

	switch {
	case err == nil:
		return c.JSON(http.StatusOK, m)
	case err == sql.ErrNoRows:
		return c.JSON(http.StatusConflict, "movie not found")
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func createMovieHandler(c echo.Context) error {
	m := Movie{}

	err := c.Bind(&m)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error)
	}

	stmt, err := db.Prepare(`
	INSERT INTO iDB (imdbID, title, year, rating, isSuperHero)
	VALUES (?, ?, ?, ?, ?);
	`)
	defer func() {
		_ = stmt.Close()
	}()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	res, err := stmt.Exec(m.ImdbID, m.Title, m.Year, m.Rating, strconv.FormatBool(m.IsSuperHero))

	switch {
	case err == nil:
		var id int64
		id, err = res.LastInsertId()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		m.ID = id
		return c.JSON(http.StatusCreated, m)
	case err.Error() == "UNIQUE constraint violation":
		return c.JSON(http.StatusConflict, "movie already exists")
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

var db *sql.DB

func conn() {
	var err error
	db, err = sql.Open("ramsql", "iDB")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Welcome to iCinema")

	conn()

	createTb := `
	CREATE TABLE IF NOT EXISTS iDB (
	id INT AUTO_INCREMENT,
	imdbID TEXT NOT NULL UNIQUE,
	title TEXT NOT NULL,
	year INT NOT NULL,
	rating FLOAT NOT NULL,
	isSuperHero BOOLEAN NOT NULL,
	PRIMARY KEY (id)
	);
	`

	if _, err := db.Exec(createTb); err != nil {
		log.Fatal("Create table error", err)
	}

	e := echo.New()

	e.GET("/movies", getAllMoviesHandler)
	e.GET("/movies/:imdbID", getMovieByIdHandler)

	e.POST("/movies", createMovieHandler)

	port := "80"
	log.Println("Start at port:" + port)

	log.Fatal(e.Start(":" + port))
}
