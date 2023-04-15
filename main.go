package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Movie struct {
	ImdbID      string  `json:"imdbID"`
	Title       string  `json:"titel"`
	Year        int     `json:"year"`
	Rating      float32 `json:"rating"`
	IsSuperHero bool    `json:"isSuperHero"`
}

var movies = []Movie{
	{
		ImdbID:      "tt4154796",
		Title:       "Avenger: Endgame",
		Year:        2019,
		Rating:      8.4,
		IsSuperHero: true,
	},
}

func getAllMoviesHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, movies)
}

func main() {
	fmt.Println("Welcome to iCinema")

	e := echo.New()

	e.GET("/movies", getAllMoviesHandler)

	port := "80"
	log.Println("Start at port:" + port)

	log.Fatal(e.Start(":" + port))
}
