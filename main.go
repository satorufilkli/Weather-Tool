package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		Humidity int     `json:"humidity"`
		WindKph  float64 `json:"wind_kph"`
		WindDir  string  `json:"wind_dir"`
	} `json:"current"`

	Forecast struct {
		ForecastDay []struct {
			Date string `json:"date"`
			Day  struct {
				MaxtempC  float64 `json:"maxtemp_c"`
				MintempC  float64 `json:"mintemp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				MaxwindKph    float64 `json:"maxwind_kph"`
				TotalPrecipMm float64 `json:"totalprecip_mm"`
				AvgHumidity   float64 `json:"avghumidity"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func GetForecast(city string, apiKey string, days uint) (*Weather, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=%d", apiKey, city, days)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %v", err)
	}

	var weather Weather
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	return &weather, nil
}

func (w *Weather) ShowCurrentWeather() {
	fmt.Printf("\n=== Current Weather in %s, %s ===\n", w.Location.Name, w.Location.Country)
	fmt.Printf("Temperature: %.1f°C\n", w.Current.TempC)
	fmt.Printf("Condition: %s\n", w.Current.Condition.Text)
	fmt.Printf("Humidity: %d%%\n", w.Current.Humidity)
	fmt.Printf("Wind: %.1f km/h from %s\n", w.Current.WindKph, w.Current.WindDir)
}

func (w *Weather) ShowForecast() {
	fmt.Printf("\n=== Weather Forecast for %s ===\n", w.Location.Name)
	for _, day := range w.Forecast.ForecastDay {
		fmt.Printf("\nDate: %s\n", day.Date)
		fmt.Printf("  Max Temperature: %.1f°C\n", day.Day.MaxtempC)
		fmt.Printf("  Min Temperature: %.1f°C\n", day.Day.MintempC)
		fmt.Printf("  Condition: %s\n", day.Day.Condition.Text)
		fmt.Printf("  Max Wind: %.1f km/h\n", day.Day.MaxwindKph)
		fmt.Printf("  Precipitation: %.1f mm\n", day.Day.TotalPrecipMm)
		fmt.Printf("  Average Humidity: %.1f%%\n", day.Day.AvgHumidity)
	}
}

func showMenu() {
	fmt.Println("\n=== Weather Information System ===")
	fmt.Println("1. Check current weather")
	fmt.Println("2. View weather forecast")
	fmt.Println("3. Change city")
	fmt.Println("4. Exit")
	fmt.Print("Please enter your choice (1-4): ")
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	const apiKey = "89f887f2bc8241918bd114435250904"
	var city string

	fmt.Print("Enter city name: ")
	city = readInput()

	var weather *Weather
	var err error

	for {
		weather, err = GetForecast(city, apiKey, 3) // 获取3天预报
		if err != nil {
			log.Printf("Error getting forecast: %v", err)
			fmt.Print("Please enter a valid city name: ")
			city = readInput()
			continue
		}
		break
	}

	for {
		showMenu()
		choice := readInput()

		switch choice {
		case "1":
			weather.ShowCurrentWeather()
		case "2":
			weather.ShowForecast()
		case "3":
			fmt.Print("Enter new city name: ")
			city = readInput()
			weather, err = GetForecast(city, apiKey, 3)
			if err != nil {
				log.Printf("Error getting forecast: %v\n", err)
				continue
			}
			fmt.Printf("Changed to %s successfully!\n", city)
		case "4":
			fmt.Println("Thank you for using Weather Information System. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
