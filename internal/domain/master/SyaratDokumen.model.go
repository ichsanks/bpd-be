package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type SyaratDokumen struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	IdTransaksi    string     `db:"id_transaksi" json:"idTransaksi"`
	IdDokumen      string     `db:"id_dokumen" json:"idDokumen"`
	Nama           string     `db:"nama" json:"nama"`
	JenisTransaksi string     `db:"jenis_transaksi" json:"jenisTransaksi"`
	Mandatori      string     `db:"mandatori" json:"mandatori"`
	IsMandatory    bool       `db:"is_mandatory" json:"isMandatory"`
	TenantID       uuid.UUID  `db:"tenant_id" json:"tenantId"`
	IdBranch       *string    `db:"id_branch" json:"idBranch"`
	CreatedAt      *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy      *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted      bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestSyaratDokumen struct {
	ID          uuid.UUID `db:"id" json:"id"`
	IdTransaksi string    `db:"id_transaksi" json:"idTransaksi"`
	IdDokumen   string    `db:"id_dokumen" json:"idDokumen"`
	IsMandatory bool      `db:"is_mandatory" json:"isMandatory"`
	IdBranch    *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappSyaratDokumen = map[string]interface{}{
	"id":        "id",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

func (b *SyaratDokumen) SyaratDokumenFormatRequest(reqFormat RequestSyaratDokumen, userId uuid.UUID, tenantId uuid.UUID) (newSyaratDokumen SyaratDokumen, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newSyaratDokumen = SyaratDokumen{
			ID:          newID,
			IdTransaksi: reqFormat.IdTransaksi,
			IdDokumen:   reqFormat.IdDokumen,
			IsMandatory: reqFormat.IsMandatory,
			TenantID:    tenantId,
			IdBranch:    reqFormat.IdBranch,
			CreatedAt:   &now,
			CreatedBy:   &userId,
		}
	} else {
		newSyaratDokumen = SyaratDokumen{
			ID:          reqFormat.ID,
			IdTransaksi: reqFormat.IdTransaksi,
			IdDokumen:   reqFormat.IdDokumen,
			IsMandatory: reqFormat.IsMandatory,
			TenantID:    tenantId,
			IdBranch:    reqFormat.IdBranch,
			UpdatedAt:   &now,
			UpdatedBy:   &userId,
		}
	}
	return
}

func (b *SyaratDokumen) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &cTime
	b.UpdatedBy = &userID
}
