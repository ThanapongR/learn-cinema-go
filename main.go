package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	y := c.QueryParam("year")

	if y == "" {
		return c.JSON(http.StatusOK, movies)
	}

	year, err := strconv.Atoi(y)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ms := []Movie{}

	for _, m := range movies {
		if year == m.Year {
			ms = append(ms, m)
		}
	}

	return c.JSON(http.StatusOK, ms)
}

func getMovieByIdHandler(c echo.Context) error {
	id := c.Param("id")

	for _, m := range movies {
		if m.ImdbID == id {
			return c.JSON(http.StatusOK, m)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{"message": "not found"})
}

func main() {
	fmt.Println("Welcome to iCinema")

	e := echo.New()

	e.GET("/movies", getAllMoviesHandler)
	e.GET("/movies/:id", getMovieByIdHandler)

	port := "80"
	log.Println("Start at port:" + port)

	log.Fatal(e.Start(":" + port))
}
