package models

import "time"

type Dokumen struct {
	ID          int
	IDKunjungan int
	Nama        string
	Path        string
	CreatedAt   time.Time
}
