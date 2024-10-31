package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type LevelBod struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Kode      string     `db:"kode" json:"kode"`
	Nama      string     `db:"nama" json:"nama"`
	Level     string     `db:"level" json:"level"`
	TenantID  *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch  *string    `db:"id_branch" json:"idBranch"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestLevelBod struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Kode     string    `db:"kode" json:"kode"`
	Nama     string    `db:"nama" json:"nama"`
	Level    string    `db:"level" json:"level"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappLevelBod = map[string]interface{}{
	"id":        "id",
	"kode":      "kode",
	"nama":      "nama",
	"level":     "level",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

func (b *LevelBod) LevelBodFormatRequest(reqFormat RequestLevelBod, userId uuid.UUID, tenantId uuid.UUID) (newLevelBod LevelBod, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newLevelBod = LevelBod{
			ID:        newID,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			Level:     reqFormat.Level,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newLevelBod = LevelBod{
			ID:        reqFormat.ID,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			Level:     reqFormat.Level,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			UpdatedAt: &now,
			UpdatedBy: &userId,
		}
	}
	return
}

func (b *LevelBod) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &cTime
	b.UpdatedBy = &userID
}
