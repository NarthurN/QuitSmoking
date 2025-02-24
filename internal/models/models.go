package models

import "time"

type Smoker struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Experience     int       `json:"experience"`
	StoppedSmoking time.Time `json:"StoppedSmoking"`
}
