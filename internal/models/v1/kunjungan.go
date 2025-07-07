package models

import "time"

type Kunjungan struct {
	ID             int
	IDPasien       int
	IDKasus        int
	TanggalMasuk   time.Time
	JenisKunjungan time.Time
}
