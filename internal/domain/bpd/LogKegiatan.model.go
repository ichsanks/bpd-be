package bpd

import (
	"time"

	"github.com/gofrs/uuid"
)

type LogKegiatan struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	Tanggal           string     `db:"tanggal" json:"tanggal"`
	IdPerjalananDinas string     `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdBpdPegawai      *string    `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	Foto              string     `db:"foto" json:"foto"`
	Keterangan        *string    `db:"keterangan" json:"keterangan"`
	Lat               *string    `db:"lat" json:"lat"`
	Long              *string    `db:"long" json:"long"`
	Address           *string    `db:"address" json:"address"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy         *string    `db:"created_by" json:"createdBy"`
	UpdatedAt         *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy         *string    `db:"updated_by" json:"updatedBy"`
	IsDeleted         bool       `db:"is_deleted" json:"isDeleted"`
}

type PerjalananDinasDokumen struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	IdBpdPegawai string     `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	File         string     `db:"file" json:"file"`
	Keterangan   *string    `db:"keterangan" json:"keterangan"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy    *string    `db:"created_by" json:"createdBy"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy    *string    `db:"updated_by" json:"updatedBy"`
	IsDeleted    bool       `db:"is_deleted" json:"isDeleted"`
}

type LogKegiatanRequest struct {
	ID                uuid.UUID `db:"id" json:"id"`
	Tanggal           string    `db:"tanggal" json:"tanggal"`
	IdPerjalananDinas string    `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdBpdPegawai      *string   `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	Foto              string    `db:"foto" json:"foto"`
	Keterangan        *string   `db:"keterangan" json:"keterangan"`
	Lat               *string   `db:"lat" json:"lat"`
	Long              *string   `db:"long" json:"long"`
	Address           *string   `db:"address" json:"address"`
	Filename          string    `json:"-"`
}

type PerjalananDinasDokumenRequest struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	IdBpdPegawai string     `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	File         string     `db:"file" json:"file"`
	Keterangan   *string    `db:"keterangan" json:"keterangan"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy    *string    `db:"created_by" json:"createdBy"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy    *string    `db:"updated_by" json:"updatedBy"`
	IsDeleted    bool       `db:"is_deleted" json:"isDeleted"`
}

func (l *LogKegiatan) NewLogKegiatanFormat(reqFormat LogKegiatanRequest, userID string) (pd LogKegiatan, err error) {
	newID, _ := uuid.NewV4()
	now := time.Now()

	var tgl string
	if reqFormat.Tanggal != "" {
		tgl = reqFormat.Tanggal
	} else {
		tgl = now.Format("2006-01-02")
	}

	if reqFormat.ID == uuid.Nil {
		pd = LogKegiatan{
			ID:                newID,
			Tanggal:           tgl,
			IdPerjalananDinas: reqFormat.IdPerjalananDinas,
			IdBpdPegawai:      reqFormat.IdBpdPegawai,
			Foto:              reqFormat.Foto,
			Keterangan:        reqFormat.Keterangan,
			Lat:               reqFormat.Lat,
			Long:              reqFormat.Long,
			Address:           reqFormat.Address,
			CreatedAt:         time.Now(),
			CreatedBy:         &userID,
		}
	} else {
		pd = LogKegiatan{
			ID:                reqFormat.ID,
			Tanggal:           tgl,
			IdPerjalananDinas: reqFormat.IdPerjalananDinas,
			IdBpdPegawai:      reqFormat.IdBpdPegawai,
			Foto:              reqFormat.Foto,
			Keterangan:        reqFormat.Keterangan,
			Lat:               reqFormat.Lat,
			Long:              reqFormat.Long,
			Address:           reqFormat.Address,
			UpdatedAt:         &now,
			UpdatedBy:         &userID,
		}
	}

	return
}

func (l *PerjalananDinasDokumen) NewPerjalananDinasDokumenFormat(reqFormat PerjalananDinasDokumenRequest, userID string) (pd PerjalananDinasDokumen, err error) {
	newID, _ := uuid.NewV4()
	now := time.Now()

	if reqFormat.ID == uuid.Nil {
		pd = PerjalananDinasDokumen{
			ID:           newID,
			IdBpdPegawai: reqFormat.IdBpdPegawai,
			File:         reqFormat.File,
			Keterangan:   reqFormat.Keterangan,
			CreatedAt:    time.Now(),
			CreatedBy:    &userID,
		}
	} else {
		pd = PerjalananDinasDokumen{
			ID:           reqFormat.ID,
			IdBpdPegawai: reqFormat.IdBpdPegawai,
			File:         reqFormat.File,
			Keterangan:   reqFormat.Keterangan,
			UpdatedAt:    &now,
			UpdatedBy:    &userID,
		}
	}

	return
}

func (l *LogKegiatan) SoftDelete(userId string) {
	now := time.Now()
	l.IsDeleted = true
	l.UpdatedBy = &userId
	l.UpdatedAt = &now
}
