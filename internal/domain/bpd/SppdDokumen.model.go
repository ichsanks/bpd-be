package bpd

import (
	"time"

	"github.com/gofrs/uuid"
)

type SppdDokumen struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	IdSppd     string     `db:"id_sppd" json:"idSppd"`
	File       string     `db:"file" json:"file"`
	Keterangan *string    `db:"keterangan" json:"keterangan"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy  *string    `db:"created_by" json:"createdBy"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy  *string    `db:"updated_by" json:"updatedBy"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	IdDokumen  string     `db:"id_dokumen" json:"idDokumen"`
}

type SppdDokumenRequest struct {
	ID         uuid.UUID `db:"id" json:"id"`
	IdSppd     string    `db:"id_sppd" json:"idSppd"`
	File       string    `db:"file" json:"file"`
	Keterangan *string   `db:"keterangan" json:"keterangan"`
	IdDokumen  string    `db:"id_dokumen" json:"idDokumen"`
}

func (s *SppdDokumen) NewSppdDokumenFormat(reqFormat SppdDokumenRequest, userID string) (pd SppdDokumen, err error) {
	newID, _ := uuid.NewV4()
	now := time.Now()

	if reqFormat.ID == uuid.Nil {
		pd = SppdDokumen{
			ID:         newID,
			IdSppd:     reqFormat.IdSppd,
			File:       reqFormat.File,
			Keterangan: reqFormat.Keterangan,
			IdDokumen:  reqFormat.IdDokumen,
			CreatedAt:  &now,
			CreatedBy:  &userID,
		}
	} else {

		pd = SppdDokumen{
			ID:         reqFormat.ID,
			IdSppd:     reqFormat.IdSppd,
			File:       reqFormat.File,
			Keterangan: reqFormat.Keterangan,
			IdDokumen:  reqFormat.IdDokumen,
			UpdatedAt:  &now,
			UpdatedBy:  &userID,
		}
	}

	return
}

var ColumnMappSppdDokumen = map[string]interface{}{
	"id":         "id",
	"idSppd":     "id_sppd",
	"file":       "file",
	"keterangan": "keterangan",
	"idDokumen":  "id_dokumen",
	"createdAt":  "created_at",
	"createdBy":  "created_by",
	"updatedAt":  "updated_at",
	"updatedBy":  "updated_by",
	"isDeleted":  "is_deleted",
}

func (b *SppdDokumen) SoftDelete(userID string) {
	var now = time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &now
	b.UpdatedBy = &userID
}
