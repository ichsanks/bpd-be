package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type Kendaraan struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	Nama               string     `db:"nama" json:"nama"`
	Nopol              string     `db:"nopol" json:"nopol"`
	IdJenisKendaraan   uuid.UUID  `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
	NamaJenisKendaraan *string    `db:"nama_jenis_kendaraan" json:"namaJenisKendaraan"`
	CreatedAt          *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy          *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy          *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted          bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestKendaraan struct {
	ID               uuid.UUID `db:"id" json:"id"`
	Nopol            string    `db:"nopol" json:"nopol"`
	Nama             string    `db:"nama" json:"nama"`
	IdJenisKendaraan uuid.UUID `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
}

var ColumnMappKendaraan = map[string]interface{}{
	"id":                 "mk.id",
	"nopol":              "mk.nopol",
	"nama":               "mk.nama",
	"namaJenisKendaraan": "mjk.nama", // untuk diurutkan berdasarkan nama jenis kendaraan (mjk. untuk menjelaskan jika hasil dari join dengan tabel lain)
	"createdAt":          "mk.created_at",
	"createdBy":          "mk.created_by",
	"updatedAt":          "mk.updated_at",
	"updatedBy":          "mk.updated_by",
	"isDeleted":          "mk.is_deleted",
}

func (b *Kendaraan) KendaraanFormatRequest(reqFormat RequestKendaraan, userId uuid.UUID) (newKendaraan Kendaraan, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newKendaraan = Kendaraan{
			ID:               newID,
			Nopol:            reqFormat.Nopol,
			Nama:             reqFormat.Nama,
			IdJenisKendaraan: reqFormat.IdJenisKendaraan,
			CreatedAt:        &now,
			CreatedBy:        &userId,
		}
	} else {
		newKendaraan = Kendaraan{
			ID:               reqFormat.ID,
			Nopol:            reqFormat.Nopol,
			Nama:             reqFormat.Nama,
			IdJenisKendaraan: reqFormat.IdJenisKendaraan,
			UpdatedAt:        &now,
			UpdatedBy:        &userId,
		}
	}
	return
}

func (b *Kendaraan) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	b.IsDeleted = true
	b.Nama = b.Nama + " [Deleted at " + cTime.Format("2006-01-02 15:04:05") + "]"
}
