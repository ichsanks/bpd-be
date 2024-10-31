package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type Golongan struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Kode      string     `db:"kode" json:"kode"`
	Nama      string     `db:"nama" json:"nama"`
	TenantID  *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch  *string    `db:"id_branch" json:"idBranch"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestGolongan struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Kode     string    `db:"kode" json:"kode"`
	Nama     string    `db:"nama" json:"nama"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappGolongan = map[string]interface{}{
	"id":        "id",
	"kode":      "kode",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

func (b *Golongan) GolonganFormatRequest(reqFormat RequestGolongan, userId uuid.UUID, tenantId uuid.UUID) (newGolongan Golongan, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newGolongan = Golongan{
			ID:        newID,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newGolongan = Golongan{
			ID:        reqFormat.ID,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			UpdatedAt: &now,
			UpdatedBy: &userId,
		}
	}
	return
}

func (b *Golongan) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &cTime
	b.UpdatedBy = &userID
}
