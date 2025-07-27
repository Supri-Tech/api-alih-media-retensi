package models

import "time"

type Log struct {
	ID     int
	Waktu  time.Time
	User   int
	Pesan  string
	Jenis  string
	Status string
}
