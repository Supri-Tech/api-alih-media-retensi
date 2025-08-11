package models

import "time"

type Pemusnahan struct {
	ID         int
	Status     string
	TglLaporan time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PemusnahanJoin struct {
	ID             int // Pemusnahan
	TglLaporan     time.Time
	Status         string
	TglMasuk       time.Time // Kunjungan
	JenisKunjungan string
	NoRM           string // Pasien
	NamaPasien     string
	JenisKelamin   string
	TglLahir       time.Time
	Alamat         string
	StatusPasien   string
	JenisKasus     string // Kasus
	MasaAktifRi    int
	MasaInaktifRi  int
	MasaAktifRj    int
	MasaInaktifRj  int
	InfoLain       string
}
