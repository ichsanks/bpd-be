package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type JenisTujuan struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	Nama       string     `db:"nama" json:"nama"`
	Keterangan *string    `db:"keterangan" json:"keterangan"`
	TenantID   *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch   *string    `db:"id_branch" json:"idBranch"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy  *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy  *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestJenisTujuan struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Nama       string    `db:"nama" json:"nama"`
	Keterangan *string   `db:"keterangan" json:"keterangan"`
	IdBranch   *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappJenisTujuan = map[string]interface{}{
	"id":         "id",
	"nama":       "nama",
	"keterangan": "keterangan",
	"createdAt":  "created_at",
	"createdBy":  "created_by",
	"updatedAt":  "updated_at",
	"updatedBy":  "updated_by",
	"isDeleted":  "is_deleted",
}

func (b *JenisTujuan) JenisTujuanFormatRequest(reqFormat RequestJenisTujuan, userId uuid.UUID, tenanId uuid.UUID) (newJenisTujuan JenisTujuan, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newJenisTujuan = JenisTujuan{
			ID:         newID,
			Nama:       reqFormat.Nama,
			Keterangan: reqFormat.Keterangan,
			TenantID:   &tenanId,
			IdBranch:   reqFormat.IdBranch,
			CreatedAt:  &now,
			CreatedBy:  &userId,
		}
	} else {
		newJenisTujuan = JenisTujuan{
			ID:         reqFormat.ID,
			Nama:       reqFormat.Nama,
			Keterangan: reqFormat.Keterangan,
			TenantID:   &tenanId,
			IdBranch:   reqFormat.IdBranch,
			UpdatedAt:  &now,
			UpdatedBy:  &userId,
		}
	}
	return
}

func (b *JenisTujuan) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &cTime
	b.UpdatedBy = &userID
}
