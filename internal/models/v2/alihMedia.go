package models

import "time"

type AlihMedia struct {
	ID         int
	Status     string
	TglLaporan time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
