package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type SettingBiaya struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	Nama              *string    `db:"nama" json:"nama"`
	LevelBod          *string    `db:"level_bod" json:"levelBod"`
	JenisTujuan       *string    `db:"jenis_tujuan" json:"jenisTujuan"`
	IdJenisBiaya      string     `db:"id_jenis_biaya" json:"idJenisBiaya"`
	CreatedAt         *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy         *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt         *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy         *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted         bool       `db:"is_deleted" json:"isDeleted"`
	IsHarian          *bool      `db:"is_harian" json:"isHarian"`
	TenantID          *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch          *string    `db:"id_branch" json:"idBranch"`
	IdBodLevel        *string    `db:"id_bod_level" json:"idBodLevel"`
	JumlahBiaya       *float64   `db:"jumlah_biaya" json:"jumlahBiaya"`
	IdJenisTujuan     *string    `db:"id_jenis_tujuan" json:"idJenisTujuan"`
	IdKategoriBiaya   *string    `db:"id_kategori_biaya" json:"idKategoriBiaya"`
	NamaKategoriBiaya *string    `db:"nama_kategori_biaya" json:"namaKategoriBiaya"`
	IsMax             *bool      `db:"is_max" json:"isMax"`
	JumlahHari        *int       `db:"jumlah_hari" json:"jumlahHari"`
}

type SettingBiayaFormat struct {
	ID            uuid.UUID  `json:"id"`
	IdJenisBiaya  string     `db:"id_jenis_biaya" json:"idJenisBiaya"`
	IsHarian      *bool      `json:"isHarian"`
	IdBranch      *string    `db:"id_branch" json:"idBranch"`
	TenantID      *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBodLevel    *string    `db:"id_bod_level" json:"idBodLevel"`
	JumlahBiaya   *float64   `db:"jumlah_biaya" json:"jumlahBiaya"`
	IdJenisTujuan *string    `db:"id_jenis_tujuan" json:"idJenisTujuan"`
	IsMax         *bool      `db:"is_max" json:"isMax"`
	JumlahHari    *int       `db:"jumlah_hari" json:"jumlahHari"`
}

type SettingBiayaUpdateFormat struct {
	ID            uuid.UUID  `json:"id"`
	IdJenisBiaya  string     `db:"id_jenis_biaya" json:"idJenisBiaya"`
	IsHarian      *bool      `json:"isHarian"`
	IdBranch      *string    `db:"id_branch" json:"idBranch"`
	TenantID      *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBodLevel    *string    `db:"id_bod_level" json:"idBodLevel"`
	JumlahBiaya   *float64   `db:"jumlah_biaya" json:"jumlahBiaya"`
	IdJenisTujuan *string    `db:"id_jenis_tujuan" json:"idJenisTujuan"`
	IsMax         *bool      `db:"is_max" json:"isMax"`
	JumlahHari    *int       `db:"jumlah_hari" json:"jumlahHari"`
}

// alis dari json ke db untuk sort table fe
var ColumnMappSettingBiaya = map[string]interface{}{
	"id":          "a.id",
	"nama":        "b.nama",
	"levelBod":    "c.nama",
	"jenisTujuan": "d.nama",
	"createdAt":   "a.created_at",
	"createdBy":   "a.created_by",
	"updatedAt":   "a.updated_at",
	"updatedBy":   "a.updated_by",
	"isDeleted":   "a.is_deleted",
	"isHarian":    "a.is_harian",
}

// field create dan update
func (settingBiaya *SettingBiaya) SettingBiayaNewFormat(reqFormat SettingBiayaFormat, userId uuid.UUID) (newSettingBiaya SettingBiaya, err error) {

	var now = time.Now()
	newID, _ := uuid.NewV4()
	newSettingBiaya = SettingBiaya{
		ID:            newID,
		IdBodLevel:    reqFormat.IdBodLevel,
		IdJenisBiaya:  reqFormat.IdJenisBiaya,
		IsHarian:      reqFormat.IsHarian,
		CreatedAt:     &now,
		CreatedBy:     &userId,
		TenantID:      reqFormat.TenantID,
		IdBranch:      reqFormat.IdBranch,
		JumlahBiaya:   reqFormat.JumlahBiaya,
		IdJenisTujuan: reqFormat.IdJenisTujuan,
		IsMax:         reqFormat.IsMax,
		JumlahHari:    reqFormat.JumlahHari,
	}

	return
}

func (settingBiaya *SettingBiaya) SettingBiayaUpdateFormat(reqFormat SettingBiayaUpdateFormat, userId uuid.UUID) (newSettingBiaya SettingBiaya, err error) {
	var now = time.Now()
	newSettingBiaya = SettingBiaya{
		ID:            reqFormat.ID,
		IdJenisBiaya:  reqFormat.IdJenisBiaya,
		IsHarian:      reqFormat.IsHarian,
		CreatedAt:     &now,
		CreatedBy:     &userId,
		TenantID:      reqFormat.TenantID,
		IdBranch:      reqFormat.IdBranch,
		IdBodLevel:    reqFormat.IdBodLevel,
		JumlahBiaya:   reqFormat.JumlahBiaya,
		IdJenisTujuan: reqFormat.IdJenisTujuan,
		IsMax:         reqFormat.IsMax,
		JumlahHari:    reqFormat.JumlahHari,
	}

	return
}

// field delete soft
func (settingBiaya *SettingBiaya) SoftDelete(userID uuid.UUID) {
	// cTime := time.Now()
	settingBiaya.IsDeleted = true
	// settingBiaya.Nama = settingBiaya.Nama + " [Deleted at " + cTime.Format("2006-01-02 15:04:05") + "]"
}
