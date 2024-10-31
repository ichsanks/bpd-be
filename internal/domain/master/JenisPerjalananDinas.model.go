package master

import (
	"time"

	"github.com/gofrs/uuid"
)

// field untuk transaksi
type JenisPerjalananDinas struct {
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

// field param di swagger
type JenisPerjalananDinasFormat struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Nama     string    `db:"nama" json:"nama"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

// alis dari json ke db untuk sort table fe
var ColumnMappJenisPerjalananDinas = map[string]interface{}{
	"id":        "id",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

// field create dan update
func (jenisPerjalananDinas *JenisPerjalananDinas) JenisPerjalananDinasFormat(reqFormat JenisPerjalananDinasFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisPerjalananDinas JenisPerjalananDinas, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newJenisPerjalananDinas = JenisPerjalananDinas{
			ID:        newID,
			Nama:      reqFormat.Nama,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newJenisPerjalananDinas = JenisPerjalananDinas{
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

// field delete soft
func (jenisPerjalananDinas *JenisPerjalananDinas) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	jenisPerjalananDinas.IsDeleted = true
	jenisPerjalananDinas.Nama = jenisPerjalananDinas.Nama + " [Deleted at " + cTime.Format("2006-01-02 15:04:05") + "]"
}
