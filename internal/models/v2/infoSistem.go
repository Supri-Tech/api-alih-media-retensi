package models

import "time"

type InfoSistem struct {
	ID           int
	NamaAplikasi string
	Logo         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
