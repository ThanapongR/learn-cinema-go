package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func getAllMoviesHandler(c echo.Context) error {
	return c.JSON(200, []string{"movie1", "movie2"})
}

func main() {
	fmt.Println("Welcome to iCinema")

	e := echo.New()

	e.GET("/movies", getAllMoviesHandler)

	port := "80"
	log.Println("Start at port:" + port)

	log.Fatal(e.Start(":" + port))
}
