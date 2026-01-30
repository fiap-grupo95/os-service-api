package usecase

import (
	"errors"
	"math"
	"time"

	"mecanica_xpto/internal/infrastructure/logs"
)


func CalculateExecutionDurationInHours(startedExecutionDate, finalExecutionDate *time.Time) (float64, error) {
	if startedExecutionDate == nil || finalExecutionDate == nil {
		logs.Logger().Error().Msg("Dates cannot be nil")
		return 0, errors.New("Dates cannot be nil")
	}

	// Check if the final date is earlier than the start date
	if finalExecutionDate.Before(*startedExecutionDate) {
		logs.Logger().Error().Msg("The final date cannot be earlier than the start date")
		return 0, errors.New("The final date cannot be earlier than the start date")
	}

	// Calculate the duration in hours
	duration := finalExecutionDate.Sub(*startedExecutionDate)
	return round(duration.Hours(), 2), nil
}

func round(value float64, precision int) float64 {
	return math.Round(value*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
}
