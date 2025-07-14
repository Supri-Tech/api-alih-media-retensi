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
	ID         int
	NamaPasien string
	NoRM       string
	TglLahir   time.Time
	Alamat     string
	JenisKasus string
}
