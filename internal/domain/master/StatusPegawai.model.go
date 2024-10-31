package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type StatusPegawai struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Nama      string     `db:"nama" json:"nama"`
	TenantID  *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch  *string    `db:"id_branch" json:"idBranch"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestStatusPegawai struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Nama     string    `db:"nama" json:"nama"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappStatusPegawai = map[string]interface{}{
	"id":        "id",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

func (b *StatusPegawai) StatusPegawaiFormatRequest(reqFormat RequestStatusPegawai, userId uuid.UUID, tenantId uuid.UUID) (newStatusPegawai StatusPegawai, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newStatusPegawai = StatusPegawai{
			ID:        newID,
			Nama:      reqFormat.Nama,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newStatusPegawai = StatusPegawai{
			ID:        reqFormat.ID,
			Nama:      reqFormat.Nama,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			UpdatedAt: &now,
			UpdatedBy: &userId,
		}
	}
	return
}

func (b *StatusPegawai) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &cTime
	b.UpdatedBy = &userID
}
