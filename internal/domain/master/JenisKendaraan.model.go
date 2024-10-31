package master

import (
	"time"

	"github.com/gofrs/uuid"
)

// field untuk transaksi
type JenisKendaraan struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	Nama           string     `db:"nama" json:"nama"`
	CreatedAt      *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy      *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted      bool       `db:"is_deleted" json:"isDeleted"`
	PilihKendaraan *bool      `db:"pilih_kendaraan" json:"pilihKendaraan"`
}

// field param di swagger
type JenisKendaraanFormat struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Nama           string    `db:"nama" json:"nama"`
	PilihKendaraan *bool     `db:"pilih_kendaraan" json:"pilihKendaraan"`
}

// alis dari json ke db untuk sort table fe
var ColumnMappJenisKendaraan = map[string]interface{}{
	"id":        "id",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

// field create dan update
func (jenisKendaraan *JenisKendaraan) JenisKendaraanFormat(reqFormat JenisKendaraanFormat, userId uuid.UUID) (newJenisKendaraan JenisKendaraan, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newJenisKendaraan = JenisKendaraan{
			ID:             newID,
			Nama:           reqFormat.Nama,
			PilihKendaraan: reqFormat.PilihKendaraan,
			CreatedAt:      &now,
			CreatedBy:      &userId,
		}
	} else {
		newJenisKendaraan = JenisKendaraan{
			ID:             reqFormat.ID,
			Nama:           reqFormat.Nama,
			PilihKendaraan: reqFormat.PilihKendaraan,
			UpdatedAt:      &now,
			UpdatedBy:      &userId,
		}
	}
	return
}

// field delete soft
func (jenisKendaraan *JenisKendaraan) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	jenisKendaraan.IsDeleted = true
	jenisKendaraan.UpdatedAt = &cTime
	jenisKendaraan.UpdatedBy = &userID
}
