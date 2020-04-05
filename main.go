package main

import (
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/ydlamia/crawl-golang/scrapper"
)

var csvFile string = "jobs.csv"

func handleHome(c echo.Context) error {
	// return c.String(http.StatusOK, "Hello, World!")
	return c.File("home.html")
}

//GET
func handleSearch(c echo.Context) error {
	defer os.Remove(csvFile)
	searchString := c.QueryParam("searchValue")
	scrapper.Scrape(searchString)
	// printvalue := searchString + " Crawling Done"
	return c.Attachment(csvFile, csvFile)
	// return c.String(http.StatusOK, printvalue)
}

//POST
func handleScrape(c echo.Context) error {
	defer os.Remove(csvFile)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	// printvalue := term + " Crawling Done"
	return c.Attachment(csvFile, csvFile)
	// return c.String(http.StatusOK, printvalue)
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.GET("/search", handleSearch)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1233"))
}
