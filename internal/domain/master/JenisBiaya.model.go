package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type JenisBiaya struct {
	ID              uuid.UUID       `db:"id" json:"id"`
	Nama            string          `db:"nama" json:"nama"`
	CreatedAt       *time.Time      `db:"created_at" json:"createdAt"`
	CreatedBy       *uuid.UUID      `db:"created_by" json:"createdBy"`
	UpdatedAt       *time.Time      `db:"updated_at" json:"updatedAt"`
	UpdatedBy       *uuid.UUID      `db:"updated_by" json:"updatedBy"`
	IsDeleted       bool            `db:"is_deleted" json:"isDeleted"`
	IsMultiple      bool            `db:"is_multiple" json:"isMultiple"`
	Urut            *int            `db:"urut" json:"urut"`
	TenantID        *uuid.UUID      `db:"tenant_id" json:"tenantId"`
	IdBranch        *string         `db:"id_branch" json:"idBranch"`
	IdKategoriBiaya *string         `db:"id_kategori_biaya" json:"idKategoriBiaya"`
	KelompokBiaya   *string         `db:"kelompok_biaya" json:"kelompokBiaya"`
	Detail          []KomponenBiaya `db:"-" json:"detail"`
}

type KomponenBiaya struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	Nama         string     `db:"nama" json:"nama"`
	IdJenisBiaya string     `db:"id_jenis_biaya" json:"idJenisBiaya"`
	CreatedAt    *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy    *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy    *uuid.UUID `db:"updated_by" json:"updatedBy"`
	IsDeleted    bool       `db:"is_deleted" json:"isDeleted"`
	Urut         int        `db:"urut" json:"urut"`
	IsHarian     bool       `db:"is_harian" json:"isHarian"`
	TenantID     *uuid.UUID `db:"tenant_id" json:"tenantId"`
	IdBranch     *string    `db:"id_branch" json:"idBranch"`
}

// field param di swagger
type JenisBiayaFormat struct {
	ID              uuid.UUID             `json:"id"`
	Nama            string                `json:"nama"`
	IsMultiple      bool                  `json:"isMultiple"`
	Urut            *int                  `db:"urut" json:"urut"`
	IdBranch        *string               `db:"id_branch" json:"idBranch"`
	IdKategoriBiaya *string               `db:"id_kategori_biaya" json:"idKategoriBiaya"`
	KelompokBiaya   *string               `db:"kelompok_biaya" json:"kelompokBiaya"`
	Detail          []KomponenBiayaFormat `json:"detail"`
}

type JenisBiayaDto struct {
	ID              uuid.UUID `db:"id" json:"id"`
	Nama            string    `db:"nama" json:"nama"`
	IdKomponenBiaya *string   `db:"id_komponen_biaya" json:"idKomponenBiaya"`
	JumlahBiaya     *int64    `db:"jumlah_biaya" json:"jumlahBiaya"`
	JumlahHari      *int      `db:"jumlah_hari" json:"jumlahHari"`
	IsMax           *bool     `db:"is_max" json:"isMax"`
}

type KomponenBiayaFormat struct {
	ID       uuid.UUID `json:"id"`
	Nama     string    `json:"nama"`
	IsHarian bool      `json:"isHarian"`
	Urut     int       `json:"urut"`
	IdBranch *string   `db:"id_branch" json:"idBranch"`
}

type JenisBiayaHeader struct {
	ID       string `json:"id"`
	Nama     string `json:"nama"`
	Kategori string `json:"kategori"`
	Urut     int    `json:"urut"`
}

// alis dari json ke db untuk sort table fe
var ColumnMappJenisBiaya = map[string]interface{}{
	"id":        "id",
	"nama":      "nama",
	"createdAt": "created_at",
	"createdBy": "created_by",
	"updatedAt": "updated_at",
	"updatedBy": "updated_by",
	"isDeleted": "is_deleted",
	"isHarian":  "is_harian",
}

// field create dan update
func (jenisBiaya *JenisBiaya) JenisBiayaFormat(reqFormat JenisBiayaFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisBiaya JenisBiaya, err error) {
	newID, _ := uuid.NewV4()
	var now = time.Now()

	if reqFormat.ID == uuid.Nil {
		newJenisBiaya = JenisBiaya{
			ID:              newID,
			Nama:            reqFormat.Nama,
			CreatedAt:       &now,
			CreatedBy:       &userId,
			Urut:            reqFormat.Urut,
			IsMultiple:      reqFormat.IsMultiple,
			TenantID:        &tenantId,
			IdBranch:        reqFormat.IdBranch,
			IdKategoriBiaya: reqFormat.IdKategoriBiaya,
			KelompokBiaya:   reqFormat.KelompokBiaya,
		}
	} else {
		newJenisBiaya = JenisBiaya{
			ID:              reqFormat.ID,
			Nama:            reqFormat.Nama,
			UpdatedAt:       &now,
			UpdatedBy:       &userId,
			Urut:            reqFormat.Urut,
			IsMultiple:      reqFormat.IsMultiple,
			TenantID:        &tenantId,
			IdBranch:        reqFormat.IdBranch,
			IdKategoriBiaya: reqFormat.IdKategoriBiaya,
			KelompokBiaya:   reqFormat.KelompokBiaya,
		}
	}

	details := make([]KomponenBiaya, 0)
	for _, d := range reqFormat.Detail {
		var detID uuid.UUID
		if d.ID == uuid.Nil {
			detID, _ = uuid.NewV4()
		} else {
			detID = d.ID
		}

		newDetail := KomponenBiaya{
			ID:           detID,
			Nama:         d.Nama,
			IdJenisBiaya: newJenisBiaya.ID.String(),
			IsHarian:     d.IsHarian,
			Urut:         d.Urut,
			CreatedAt:    &now,
			CreatedBy:    &userId,
			TenantID:     &tenantId,
			IdBranch:     reqFormat.IdBranch,
		}

		details = append(details, newDetail)
	}

	newJenisBiaya.Detail = details

	return
}

// field delete soft
func (jenisBiaya *JenisBiaya) SoftDelete(userID uuid.UUID) {
	cTime := time.Now()
	jenisBiaya.IsDeleted = true
	jenisBiaya.Nama = jenisBiaya.Nama + " [Deleted at " + cTime.Format("2006-01-02 15:04:05") + "]"
}

type JumlahBiaya struct {
	JumlahBiaya int `db:"jumlah_biaya" json:"jumlahBiaya"`
}
