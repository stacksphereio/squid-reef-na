package feeds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeatherData represents weather information
type WeatherData struct {
	Summary      string  `json:"summary"`
	TemperatureC float64 `json:"temperatureC"`
	FeelsLikeC   float64 `json:"feelsLikeC"`
}

// Coordinates represents latitude and longitude
type Coordinates struct {
	Lat float64
	Lon float64
}

// OpenMeteoResponse represents the API response from Open-Meteo
type OpenMeteoResponse struct {
	Current struct {
		Temperature         float64 `json:"temperature_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		WeatherCode         int     `json:"weather_code"`
	} `json:"current"`
}

// North America country coordinates (major cities)
var naCountryCoordinates = map[string]Coordinates{
	"US": {Lat: 40.7128, Lon: -74.0060}, // New York
	"CA": {Lat: 43.6532, Lon: -79.3832}, // Toronto
	"MX": {Lat: 19.4326, Lon: -99.1332}, // Mexico City
}

// Weather code to description mapping (WMO Weather interpretation codes)
var weatherCodeDescriptions = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Foggy",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

// FetchWeather fetches weather data for a given country using Open-Meteo API
func FetchWeather(country string) (*WeatherData, error) {
	// Get coordinates for the country
	coords, ok := naCountryCoordinates[country]
	if !ok {
		// Default to New York if country not found
		coords = naCountryCoordinates["US"]
	}

	// Build Open-Meteo API URL
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,apparent_temperature,weather_code",
		coords.Lat, coords.Lon,
	)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make API request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("weather API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	// Parse response
	var apiResp OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse weather response: %w", err)
	}

	// Convert weather code to description
	description, ok := weatherCodeDescriptions[apiResp.Current.WeatherCode]
	if !ok {
		description = "Unknown"
	}

	return &WeatherData{
		Summary:      description,
		TemperatureC: apiResp.Current.Temperature,
		FeelsLikeC:   apiResp.Current.ApparentTemperature,
	}, nil
}
