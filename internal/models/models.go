package models

import "time"

type Smoker struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	StoppedSmoking time.Time `json:"stoppedSmoking"`
}
