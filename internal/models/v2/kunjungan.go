package models

import "time"

type Kunjungan struct {
	ID             int
	IDPasien       int
	IDKasus        int
	TanggalMasuk   time.Time
	JenisKunjungan string
}

type KunjunganJoin struct {
	ID             int    // kunjungan
	NamaPasien     string // pasien
	NoRM           string
	NIK            string
	JenisKelamin   string
	TglLahir       time.Time
	Alamat         string
	Status         string
	TglMasuk       time.Time // kunjungan
	IDKasus        int
	JenisKasus     string // kasus
	MasaAktifRi    int
	MasaInaktifRi  int
	MasaAktifRj    int
	MasaInaktifRj  int
	InfoLain       string
	JenisKunjungan string // kunjungan
	Dokumen        string // dokumen
}
