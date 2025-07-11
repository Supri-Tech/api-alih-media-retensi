package models

import "time"

type Pasien struct {
	ID           int
	NoRM         string
	NamaPasien   string
	JenisKelamin string
	TanggalLahir time.Time
	NIK          string
	Alamat       string
	Status       string
	CreatedAt    time.Time
}
