package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type PersonGrade struct {
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

type RequestPersonGrade struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Kode     string    `db:"kode" json:"kode"`
	Nama     string    `db:"nama" json:"nama"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappPersonGrade = map[string]interface{}{
	"id":        "id",
	"kode":      "kode",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

func (b *PersonGrade) PersonGradeFormatRequest(reqFormat RequestPersonGrade, userId uuid.UUID, tenantId uuid.UUID) (newPersonGrade PersonGrade, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newPersonGrade = PersonGrade{
			ID:        newID,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newPersonGrade = PersonGrade{
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

func (b *PersonGrade) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &cTime
	b.UpdatedBy = &userID
}
