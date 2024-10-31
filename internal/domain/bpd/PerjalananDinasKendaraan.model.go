package bpd

import (
	"time"

	"github.com/gofrs/uuid"
)

type PerjalananDinasKendaraan struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	IdPerjalananDinas string     `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdKendaraan       string     `db:"id_kendaraan" json:"idKendaraan"`
	IdPegawai         *string    `db:"id_pegawai" json:"idPegawai"`
	NamaPengemudi     *string    `db:"nama_pengemudi" json:"namaPengemudi"`
	CreatedAt         *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy         *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt         *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy         *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted         bool       `db:"is_deleted" json:"isDeleted"`
}

type PerjalananDinasKendaraanDTO struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	IdPerjalananDinas  string     `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdKendaraan        string     `db:"id_kendaraan" json:"idKendaraan"`
	Nopol              string     `db:"nopol" json:"nopol"`
	NamaKendaraan      string     `db:"nama_kendaraan" json:"namaKendaraan"`
	IdJenisKendaraan   string     `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
	NamaJenisKendaraan string     `db:"nama_jenis_kendaraan" json:"namaJenisKendaraan"`
	IdPegawai          *string    `db:"id_pegawai" json:"idPegawai"`
	NamaPengemudi      *string    `db:"nama_pengemudi" json:"namaPengemudi"`
	CreatedAt          *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy          *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy          *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted          bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestPerjalananDinasKendaraan struct {
	ID                uuid.UUID `db:"id" json:"id"`
	IdPerjalananDinas string    `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdKendaraan       string    `db:"id_kendaraan" json:"idKendaraan"`
	IdPegawai         *string   `db:"id_pegawai" json:"idPegawai"`
	NamaPengemudi     *string   `db:"nama_pengemudi" json:"namaPengemudi"`
}

func (b *PerjalananDinasKendaraan) PerjalananDinasKendaraanFormatRequest(reqFormat RequestPerjalananDinasKendaraan, userId uuid.UUID) (newBpd PerjalananDinasKendaraan, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newBpd = PerjalananDinasKendaraan{
			ID:                newID,
			IdPerjalananDinas: reqFormat.IdPerjalananDinas,
			IdKendaraan:       reqFormat.IdKendaraan,
			IdPegawai:         reqFormat.IdPegawai,
			NamaPengemudi:     reqFormat.NamaPengemudi,
			CreatedAt:         &now,
			CreatedBy:         &userId,
		}
	} else {
		newBpd = PerjalananDinasKendaraan{
			ID:                reqFormat.ID,
			IdPerjalananDinas: reqFormat.IdPerjalananDinas,
			IdKendaraan:       reqFormat.IdKendaraan,
			IdPegawai:         reqFormat.IdPegawai,
			NamaPengemudi:     reqFormat.NamaPengemudi,
			UpdatedAt:         &now,
			UpdatedBy:         &userId,
		}
	}
	return
}

var ColumnMappPDKendaraan = map[string]interface{}{
	"id":                 "pk.id",
	"idPerjalananDinas":  "pk.id_perjalanan_dinas",
	"idKendaraan":        "pk.id_kendaraan",
	"nopol":              "mk.nopol",
	"namaKendaraan":      "mk.nama",
	"idJenisKendaraan":   "mk.id_jenis_kendaraan",
	"namaJenisKendaraan": "mj.nama",
	"idPegawai":          "pk.id_pegawai",
	"namaPengemudi":      "coalesce(pk.nama_pengemudi, mp.nama)",
	"createdAt":          "pk.created_at",
	"createdBy":          "pk.created_by",
	"updatedAt":          "pk.updated_at",
	"updatedBy":          "pk.updated_by",
	"isDeleted":          "pk.is_deleted",
}

func (b *PerjalananDinasKendaraan) SoftDelete(userID uuid.UUID) {
	var now = time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &now
	b.UpdatedBy = &userID
}
