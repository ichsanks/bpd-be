package master

type JenisKelamin struct {
	Kode string `db:"kode" json:"kode"`
	Nama string `db:"nama" json:"nama"`
}
