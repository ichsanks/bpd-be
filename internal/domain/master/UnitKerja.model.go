package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type UnitOrganisasiKerja struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Kode      string     `db:"kode" json:"kode"`
	Nama      string     `db:"nama" json:"nama"`
	IdBidang  string     `db:"id_bidang" json:"idBidang"`
	TenantID  *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch  *string    `db:"id_branch" json:"idBranch"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type UnitOrganisasiKerjaDTO struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	Kode       string     `db:"kode" json:"kode"`
	Nama       string     `db:"nama" json:"nama"`
	IdBidang   *string    `db:"id_bidang" json:"idBidang"`
	KodeBidang *string    `db:"kode_bidang" json:"kodeBidang"`
	NamaBidang *string    `db:"nama_bidang" json:"namaBidang"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy  *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy  *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestUnor struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Kode     string    `db:"kode" json:"kode"`
	Nama     string    `db:"nama" json:"nama"`
	IdBidang string    `db:"id_bidang" json:"idBidang"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappUnor = map[string]interface{}{
	"id":        "id",
	"kode":      "kode",
	"nama":      "nama",
	"idBidang":  "id_bidang",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
}

func (b *UnitOrganisasiKerja) UnorFormatRequest(reqFormat RequestUnor, userId uuid.UUID, tenantId uuid.UUID) (newUnor UnitOrganisasiKerja, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newUnor = UnitOrganisasiKerja{
			ID:        newID,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			IdBidang:  reqFormat.IdBidang,
			IdBranch:  reqFormat.IdBranch,
			TenantID:  &tenantId,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newUnor = UnitOrganisasiKerja{
			ID:        reqFormat.ID,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			IdBidang:  reqFormat.IdBidang,
			IdBranch:  reqFormat.IdBranch,
			TenantID:  &tenantId,
			UpdatedAt: &now,
			UpdatedBy: &userId,
		}
	}
	return
}

func (b *UnitOrganisasiKerja) SoftDelete(userID uuid.UUID) {
	var now = time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &now
	b.UpdatedBy = &userID
}
