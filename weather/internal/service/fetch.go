package getservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weather/internal/domain"
	"weather/internal/repository"
)

// todo: move to constants
const baseURL = "https://api.open-meteo.com/v1/forecast"

type Service struct {
	http *http.Client
	repo repository.Repository
}

func New(http *http.Client, repo repository.Repository) *Service {
	return &Service{
		http: http,
		repo: repo,
	}
}

func (s *Service) FetchWeather(ctx context.Context, lat, lon float64) (domain.Weather, error) {
	url := fmt.Sprintf(
		"%s?latitude=%.4f&longitude=%.4f&current_weather=true&daily=temperature_2m_max,temperature_2m_min",
		baseURL,
		lat,
		lon,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("request creation failed: %w", err)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return domain.Weather{}, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return domain.Weather{}, fmt.Errorf("API error: %s", string(body))
	}

	// 4. Parse response
	var result struct {
		Elevation      float64 `json:"elevation"`
		CurrentWeather struct {
			Temperature float64 `json:"temperature"`
			WeatherCode int     `json:"weather_code"`
		} `json:"current_weather"`
		Daily struct {
			Time   []string  `json:"time"`
			Tmax2m []float64 `json:"temperature_2m_max"`
			Tmin2m []float64 `json:"temperature_2m_min"`
		} `json:"daily"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return domain.Weather{}, fmt.Errorf("response parsing failed: %w", err)
	}

	// 5. Map to domain object
	return domain.Weather{
		Forecast:    weatherCodeToString(result.CurrentWeather.WeatherCode),
		Temperature: result.CurrentWeather.Temperature,
		Altitude:    result.Elevation,
	}, nil
}

func weatherCodeToString(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code >= 1 && code <= 3:
		return "Partly cloudy"
	case code >= 45 && code <= 48:
		return "Fog"
	case code >= 51 && code <= 67:
		return "Rain"
	case code >= 71 && code <= 77:
		return "Snow"
	case code >= 80 && code <= 82:
		return "Showers"
	case code >= 95 && code <= 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}
