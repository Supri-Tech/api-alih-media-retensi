package models

type Kasus struct {
	ID            int
	JenisKasus    string
	MasaAktifRJ   int
	MasaInaktifRJ int
	MasaAktifRI   int
	MasaInaktifRI int
	InfoLain      string
}
