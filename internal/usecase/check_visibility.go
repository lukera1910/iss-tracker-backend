package usecase

import (
	"math"
	"iss-tracker-backend/internal/domain"
)

func IsVisible(userLat, userLon float64, iss *domain.ISSPosition) bool {
	if iss == nil {
		return false
	}

	// distância aproximada 
	dist := haversine(userLat, userLon, iss.Latitude, iss.Longitude)

	// a iss pode ser vista até +ou- 2.300 km (quando favorável)
	return dist < 2300
}

// fórmula de haversine para distãncia entre coordenadas
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // raio da terra em km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	lat1 *= math.Pi / 180
	lat2 *= math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}