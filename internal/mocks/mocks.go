package mocks

import (
	"time"

	"github.com/NarthurN/QuitSmoking/internal/models"
)

var Smokers = map[string]*models.Smoker{
	"arthurCool": {
		ID:             "1",
		Name:           "Arthur",
		Username:       "arthurCool",
		Password:       "123qwe",
		StoppedSmoking: time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
	},
	"victorCool": {
		ID:             "2",
		Name:           "Victor",
		Username:       "victorCool",
		Password:       "qasw",
		StoppedSmoking: time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
	},
}
