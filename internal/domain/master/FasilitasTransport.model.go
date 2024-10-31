package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type FasilitasTransport struct {
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

type RequestFasilitasTransport struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Nama     string    `db:"nama" json:"nama"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappFasilitasTransport = map[string]interface{}{
	"id":        "id",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

func (b *FasilitasTransport) FasilitasTransportFormatRequest(reqFormat RequestFasilitasTransport, userId uuid.UUID, tenantId uuid.UUID) (newFasilitasTransport FasilitasTransport, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newFasilitasTransport = FasilitasTransport{
			ID:        newID,
			Nama:      reqFormat.Nama,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newFasilitasTransport = FasilitasTransport{
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

func (b *FasilitasTransport) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &cTime
	b.UpdatedBy = &userID
}
