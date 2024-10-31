package master

import (
	"time"

	"github.com/gofrs/uuid"
)

// field untuk transaksi
type Jabatan struct {
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
type JabatanFormat struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Nama     string    `db:"nama" json:"nama"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

// alis dari json ke db untuk sort table fe
var ColumnMappJabatan = map[string]interface{}{
	"id":        "id",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

// field create dan update
func (jabatan *Jabatan) JabatanFormat(reqFormat JabatanFormat, userId uuid.UUID, tenantId uuid.UUID) (newJabatan Jabatan, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newJabatan = Jabatan{
			ID:        newID,
			Nama:      reqFormat.Nama,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newJabatan = Jabatan{
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
func (jabatan *Jabatan) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	jabatan.IsDeleted = true
	jabatan.Nama = jabatan.Nama + " [Deleted at " + cTime.Format("2006-01-02 15:04:05") + "]"
}
