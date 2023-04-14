package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Welcome to iCinema")

	e := echo.New()

	port := "80"
	log.Println("Start at port:" + port)

	log.Fatal(e.Start(":" + port))
}
