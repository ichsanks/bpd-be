package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type STtd struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	IdPegawai string     `db:"id_pegawai" json:"idPegawai"`
	IdJabatan string     `db:"id_jabatan" json:"idJabatan"`
	Jenis     string     `db:"jenis" json:"jenis"`
	TenantID  *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch  *string    `db:"id_branch" json:"idBranch"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
}

type STtdDto struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	IdPegawai   string     `db:"id_pegawai" json:"idPegawai"`
	Nip         string     `db:"nip" json:"nip"`
	NamaPegawai string     `db:"nama_pegawai" json:"namaPegawai"`
	IdJabatan   string     `db:"id_jabatan" json:"idJabatan"`
	NamaJabatan string     `db:"nama_jabatan" json:"namaJabatan"`
	NamaJenis   string     `db:"nama_jenis" json:"namaJenis"`
	Jenis       string     `db:"jenis" json:"jenis"`
	TenantID    *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch    *string    `db:"id_branch" json:"idBranch"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy   *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy   *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
}

type RequestSTtd struct {
	ID        uuid.UUID `db:"id" json:"id"`
	IdPegawai string    `db:"id_pegawai" json:"idPegawai"`
	IdJabatan string    `db:"id_fungsionalitas" json:"idJabatan"`
	Jenis     string    `db:"jenis" json:"jenis"`
	IdBranch  *string   `db:"id_branch" json:"idBranch"`
}

var ColumnMappSTtd = map[string]interface{}{
	"id":          "a.id",
	"namaPegawai": "b.nama",
	"namaJabatan": "c.nama",
	"createdAt":   "a.created_at",
	"createdBy":   "a.created_by",
	"updatedAt":   "a.updated_at",
	"updatedBy":   "a.updated_by",
	"isDeleted":   "a.is_deleted",
}

func (b *STtd) STtdFormatRequest(reqFormat RequestSTtd, userId uuid.UUID, tenantId uuid.UUID) (newSTtd STtd, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newSTtd = STtd{
			ID:        newID,
			IdPegawai: reqFormat.IdPegawai,
			IdJabatan: reqFormat.IdJabatan,
			Jenis:     reqFormat.Jenis,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			CreatedAt: &now,
			CreatedBy: &userId,
		}
	} else {
		newSTtd = STtd{
			ID:        reqFormat.ID,
			IdPegawai: reqFormat.IdPegawai,
			IdJabatan: reqFormat.IdJabatan,
			Jenis:     reqFormat.Jenis,
			TenantID:  &tenantId,
			IdBranch:  reqFormat.IdBranch,
			UpdatedAt: &now,
			UpdatedBy: &userId,
		}
	}
	return
}

func (b *STtd) SoftDelete(userID uuid.UUID) {
	var now = time.Now()
	b.IsDeleted = true
	b.UpdatedAt = &now
	b.UpdatedBy = &userID
}
