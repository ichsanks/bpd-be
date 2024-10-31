package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type Fungsionalitas struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	Nama          string     `db:"nama" json:"nama"`
	Level         *int       `db:"level" json:"level"`
	IsHead        *bool      `db:"is_head" json:"isHead"`
	JenisApproval *string    `db:"jenis_approval" json:"jenisApproval"`
	TenantID      *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch      *string    `db:"id_branch" json:"idBranch"`
	CreatedAt     *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy     *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy     *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted     bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestFungsionalitas struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Nama          string    `db:"nama" json:"nama"`
	Level         *int      `db:"level" json:"level"`
	IsHead        *bool     `db:"is_head" json:"isHead"`
	JenisApproval *string   `db:"jenis_approval" json:"jenisApproval"`
	IdBranch      *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappFungsionalitas = map[string]interface{}{
	"id":            "id",
	"nama":          "nama",
	"level":         "level",
	"isHead":        "is_head",
	"jenisApproval": "jenis_approval",
	"createdAt":     "created_at",
	"createdBy":     "created_by",
	"updatedAt":     "updated_at",
	"updatedBy":     "updated_by",
	"isDeleted":     "is_deleted",
}

func (b *Fungsionalitas) FungsionalitasFormatRequest(reqFormat RequestFungsionalitas, userId uuid.UUID, tenantId uuid.UUID) (newFungsionalitas Fungsionalitas, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newFungsionalitas = Fungsionalitas{
			ID:            newID,
			Nama:          reqFormat.Nama,
			Level:         reqFormat.Level,
			IsHead:        reqFormat.IsHead,
			JenisApproval: reqFormat.JenisApproval,
			TenantID:      &tenantId,
			IdBranch:      reqFormat.IdBranch,
			CreatedAt:     &now,
			CreatedBy:     &userId,
		}
	} else {
		newFungsionalitas = Fungsionalitas{
			ID:            reqFormat.ID,
			Nama:          reqFormat.Nama,
			Level:         reqFormat.Level,
			IsHead:        reqFormat.IsHead,
			JenisApproval: reqFormat.JenisApproval,
			TenantID:      &tenantId,
			IdBranch:      reqFormat.IdBranch,
			UpdatedAt:     &now,
			UpdatedBy:     &userId,
		}
	}
	return
}

func (b *Fungsionalitas) SoftDelete(userID uuid.UUID) {
	var now = time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &now
	b.UpdatedBy = &userID
}
