package pricing

import "context"

const (
	BaseCost          float64 = 100.0  // Base cost per unit (e.g., $100)
	AltitudeFactor    float64 = 0.0002 // Cost multiplier per meter (e.g., 0.0002)
	TemperatureFactor float64 = 0.01   //Cost multiplier per °C (e.g., 0.01)
)

func (s *Service) Calculate(ctx context.Context) (float64, error) {
	weather, err := s.weatherService.GetWeather(ctx)
	if err != nil {
		return 0, err
	}
	// Adjust for altitude (higher altitude = higher cost)
	altitudeAdjustment := 1 + (AltitudeFactor * weather.Altitude)

	// Adjust for temperature (colder = higher cost)
	temperatureAdjustment := 1 + (TemperatureFactor * weather.Temperature)

	return BaseCost * altitudeAdjustment * temperatureAdjustment, nil
}
