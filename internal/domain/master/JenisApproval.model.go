package master

type JenisApproval struct {
	ID               string  `db:"id" json:"id"`
	Nama             string  `db:"nama" json:"nama"`
	IdFungsionalitas *string `db:"id_fungsionalitas" json:"idFungsionalitas"`
}
