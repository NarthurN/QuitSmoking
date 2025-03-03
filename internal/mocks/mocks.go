package mocks

import (
	"time"

	"github.com/NarthurN/QuitSmoking/internal/models"
)

var Smokers = map[string]*models.Smoker{
	"1": {
		ID:             "1",
		Name:           "Arthur",
		StoppedSmoking: time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
	},
	"2": {
		ID:             "2",
		Name:           "Victor",
		StoppedSmoking: time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
	},
}
