package rest

import (
	get "gateway/internal/api/rest/handler"
)

// SetupAPIRoutes
// @title           			Gateway Service
// @version         			1.0.0
// @description     			This APIs return gateway entities
// @BasePath  					/
// @Schemes 					https
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
func (s *Server) SetupAPIRoutes(
	handler get.Handler,
) {
	r := s.engine
	{
		v1 := r.Group("api")
		v1.GET("/fuel-price", handler.GetPrice)
		v1.GET("/weather-status", handler.GetWeather)
	}
}

// todo: handle middleware
