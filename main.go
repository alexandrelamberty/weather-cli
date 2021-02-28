package main

import (
	//"time"
	//"log"
	//"strings"
	"fmt"
	"net/http"
	"os"

	"github.com/gocolly/colly"
	"github.com/olekukonko/tablewriter"
)

var websiteURL string = "https://www.meteobelgique.be/"
var cityURL string = "previsions-meteo-belgique/1450/chastre/"

// Day represent the  weather for a given days.
// TODO rename to WeekDay
type Day struct {
	Day     string
	TempMin string
	TempMax string
}

// HourWeather represent the weather for a given hour period.
type HourWeather struct {
  Hour      string
  Condition string
  Temp      string
  Rain      string
  Wind      string
}

// Today represent the weather for today
type TodayWeather struct {
  Temp    string
  Wind    string
  Hours   []HourWeather
}

// IsOnline verify if the website is online, or internet.
func IsOnline() bool {
	_, err := http.Get(websiteURL)
	if err == nil {
		return true
	}
	return false
}

// HasInternet verify if the website is online, or internet.
func HasInternet() bool {
	_, err := http.Get("https://www.google.com")
	if err == nil {
		return true
	}
	return false
}

// Convert a struc array to a string array
func convert(datas []Day) [][]string {
	var data [][]string
	for _, day := range datas {
		row := []string{day.Day, day.TempMin, day.TempMax}
		data = append(data, row)
	}
	return data
}

// Print to output with tablewriter
func printTable(days []Day) {
	var data [][]string = convert(days)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Day", "Min", "Max"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// scrape get data from internet
// TODO Move to it's own package as a specific scrapper 
// for the given website.
// Make another scrapper for meteo.be.
func scrape() {
	days := make([]Day, 0, 12)
	c := colly.NewCollector()
	c.OnHTML("div[class=menu_forecast_day]", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			day := Day{}
			day.Day = el.ChildText(".day_label")
			day.TempMin = el.ChildText(".temp_min")
			day.TempMax = el.ChildText(".temp_max")
			// Check the possibity to use some Nerd fonts.
			//day.Symbol = el.ChildText(".symbol")
			days = append(days, day)
		})
		printTable(days)
	})
	c.Visit(websiteURL + cityURL)
}

func main() {
	if IsOnline() {
		scrape()
	} else {
		fmt.Println("The website is down!")
		if HasInternet() {
			fmt.Println("You are connected to internet.")
		} else {
			fmt.Println("The internet is down!")
		}
	}
}
