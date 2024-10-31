package bpd

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type PerjalananDinasBiaya struct {
	ID                    uuid.UUID       `db:"id" json:"id"`
	IDBpdPegawai          string          `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IDJenisBiaya          string          `db:"id_jenis_biaya" json:"idJenisBiaya"`
	IDKomponenBiaya       string          `db:"id_komponen_biaya" json:"idKomponenBiaya"`
	Keterangan            string          `db:"keterangan" json:"keterangan"`
	Qty                   *int            `db:"qty" json:"qty"`
	Nominal               decimal.Decimal `db:"nominal" json:"nominal"`
	CreatedAt             *time.Time      `db:"created_at" json:"createdAt"`
	CreatedBy             *uuid.UUID      `db:"created_by" json:"createdBy"`
	UpdatedAt             *time.Time      `db:"updated_at" json:"updatedAt"`
	UpdatedBy             *uuid.UUID      `db:"updated_by" json:"updatedBy"`
	IsDeleted             bool            `db:"is_deleted" json:"isDeleted"`
	Tanggal               *string         `db:"tanggal" json:"tanggal"`
	IdUangCuci            string          `db:"id_uang_cuci" json:"idUangCuci"`
	NominalUangCuci       decimal.Decimal `db:"nominal_uang_cuci" json:"nominalUangCuci"`
	IdBiayaPenginapan     string          `db:"id_biaya_penginapan" json:"idBiayaPenginapan"`
	NominalPenginapan     decimal.Decimal `db:"nominal_penginapan" json:"nominalPenginapan"`
	IdTiketBerangkat      string          `db:"id_tiket_berangkat" json:"idTiketBerangkat"`
	NominalTiketBerangkat decimal.Decimal `db:"nominal_tiket_berangkat" json:"nominalTiketBerangkat"`
	IdTiketPulang         string          `db:"id_tiket_pulang" json:"idTiketPulang"`
	NominalTiketPulang    decimal.Decimal `db:"nominal_tiket_pulang" json:"nominalTiketPulang"`
	IdPegawai             string          `db:"id_pegawai" json:"idPegawai"`
	File                  string          `db:"file" json:"file"`
	DeletedAt             *time.Time      `db:"deleted_at" json:"deletedAt"`
	DeletedBy             *string         `db:"deleted_by" json:"deletedBy"`
	IsReimbursement       bool            `db:"is_reimbursement" json:"isReimbursement"`
	Aksi                  string          `db:"aksi" json:"aksi"`
}

type RequestPerjalananDinasBiaya struct {
	Data           []PerjalananDinasBiayaDetail `json:"data"`
	IsRevisi       bool                         `json:"isRevisi"`
	IdBpdPegawai   string                       `json:"idBpdPegawai"`
	IsUm           bool                         `json:"isUm"`
	PersentaseUm   float64                      `json:"persentaseUm"`
	PersentaseSisa float64                      `json:"persentaseSisa"`
	ShowUm         bool                         `json:"showUm"`
	ShowSisa       bool                         `json:"showSisa"`
	TotalUm        float64                      `json:"totalUm"`
	SisaUm         float64                      `json:"sisaUm"`
	InsertBiaya    bool                         `json:"insertBiaya"`
}

type PerjalananDinasBiayaDetail struct {
	ID                    string          `json:"id"`
	IDBpdPegawai          string          `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IDJenisBiaya          string          `db:"id_jenis_biaya" json:"idJenisBiaya"`
	IDKomponenBiaya       string          `db:"id_komponen_biaya" json:"idKomponenBiaya"`
	Keterangan            string          `json:"keterangan"`
	Qty                   int             `json:"qty"`
	Nominal               decimal.Decimal `db:"nominal" json:"nominal"`
	Tanggal               *string         `db:"tanggal" json:"tanggal"`
	IdUangCuci            string          `db:"id_uang_cuci" json:"idUangCuci"`
	NominalUangCuci       decimal.Decimal `db:"nominal_uang_cuci" json:"nominalUangCuci"`
	IdBiayaPenginapan     string          `db:"id_biaya_penginapan" json:"idBiayaPenginapan"`
	NominalPenginapan     decimal.Decimal `db:"nominal_penginapan" json:"nominalPenginapan"`
	IdTiketBerangkat      string          `db:"id_tiket_berangkat" json:"idTiketBerangkat"`
	NominalTiketBerangkat decimal.Decimal `db:"nominal_tiket_berangkat" json:"nominalTiketBerangkat"`
	IdTiketPulang         string          `db:"id_tiket_pulang" json:"idTiketPulang"`
	NominalTiketPulang    decimal.Decimal `db:"nominal_tiket_pulang" json:"nominalTiketPulang"`
	BiayaTransport        decimal.Decimal `db:"biaya_transport" json:"biayaTransport"`
	BiayaPerdiem          decimal.Decimal `db:"biaya_perdiem" json:"biayaPerdiem"`
	BiayaTransportLokal   decimal.Decimal `db:"biaya_transport_lokal" json:"biayaTransportLokal"`
	IdPegawai             string          `db:"id_pegawai" json:"idPegawai"`
	File                  string          `db:"file" json:"file"`
	DeletedAt             time.Time       `db:"deleted_at" json:"deletedAt"`
	DeletedBy             string          `db:"deleted_by" json:"deletedBy"`
	IsReimbursement       bool            `db:"is_reimbursement" json:"isReimbursement"`
}

type ListKomponenBiaya struct {
	ID   string `db:"id" json:"id"`
	Nama string `db:"nama" json:"nama"`
}
type PerjalananDinasBiayaDTO struct {
	IDJenisBiaya      string                   `db:"id_jenis_biaya" json:"idJenisBiaya"`
	NamaJenisBiaya    string                   `db:"nama_jenis_biaya" json:"namaJenisBiaya"`
	IsMultiple        bool                     `db:"is_multiple" json:"isMultiple"`
	ListKomponenBiaya []KomponenBiayaDetailDTO `db:"-" json:"listKomponenBiaya"`
	DetailBiaya       []KomponenBiayaDetailDTO `db:"-" json:"detailBiaya"`
}

type JmlReimbursement struct {
	Nominal decimal.Decimal `db:"nominal" json:"nominal"`
}

type BiayaPerjalananDinasDto struct {
	Id              string          `db:"id" json:"id"`
	IdBpdPegawai    string          `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IdJenisBiaya    string          `db:"id_jenis_biaya" json:"idJenisBiaya"`
	IdKomponenBiaya string          `db:"id_komponen_biaya" json:"idKomponenBiaya"`
	Nominal         decimal.Decimal `db:"nominal" json:"nominal"`
	File            string          `db:"file" json:"file"`
	Nama            *string         `db:"nama" json:"nama"`
	Keterangan      *string         `db:"keterangan" json:"keterangan"`
	Tanggal         *string         `db:"tanggal" json:"tanggal"`
	IdPegawai       *string         `db:"id_pegawai" json:"idPegawai"`
	UpdatedAt       *string         `db:"updated_at" json:"updatedAt"`
}

type HistoriPerjalananDinas struct {
	Id              string          `db:"id" json:"id"`
	IdBpdPegawai    string          `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IdJenisBiaya    string          `db:"id_jenis_biaya" json:"idJenisBiaya"`
	IdKomponenBiaya string          `db:"id_komponen_biaya" json:"idKomponenBiaya"`
	Nominal         decimal.Decimal `db:"nominal" json:"nominal"`
	File            string          `db:"file" json:"file"`
	Nama            *string         `db:"nama" json:"nama"`
	Keterangan      *string         `db:"keterangan" json:"keterangan"`
	Tanggal         *string         `db:"tanggal" json:"tanggal"`
	IdPegawai       *string         `db:"id_pegawai" json:"idPegawai"`
	DeletedAt       *string         `db:"deleted_at" json:"deletedAt"`
	DeletedBy       *string         `db:"deleted_by" json:"deletedBy"`
	UpdatedAt       *string         `db:"updated_at" json:"updatedAt"`
	UpdatedBy       *string         `db:"updated_by" json:"updatedBy"`
	CreatedAt       *string         `db:"created_at" json:"createdAt"`
	CreatedBy       *string         `db:"created_by" json:"createdBy"`
	Aksi            *string         `db:"aksi" json:"aksi"`
	UserCreate      *string         `db:"user_create" json:"userCreate"`
	UserUpdate      *string         `db:"user_update" json:"userUpdate"`
	UserDelete      *string         `db:"user_delete" json:"userDelete"`
	IsReimbursement *string         `db:"is_reimbursement" json:"isReimbursement"`
}

type KomponenBiayaDetailDTO struct {
	ID              *string          `db:"id" json:"id"`
	IDBpdPegawai    *string          `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IDKomponenBiaya string           `db:"id_komponen_biaya" json:"idKomponenBiaya"`
	IDJenisBiaya    string           `db:"id_jenis_biaya" json:"idJenisBiaya"`
	Nama            string           `db:"nama" json:"nama"`
	IsHarian        bool             `db:"is_harian" json:"isHarian"`
	Keterangan      *string          `db:"keterangan" json:"keterangan"`
	Qty             *int             `db:"qty" json:"qty"`
	Nominal         *decimal.Decimal `db:"nominal" json:"nominal"`
}

func (p *PerjalananDinasBiaya) PerjalananDinasBiayaFormatRequest(reqFormat RequestPerjalananDinasBiaya, userId uuid.UUID) (newBpd []PerjalananDinasBiaya, err error) {
	var now = time.Now()
	details := make([]PerjalananDinasBiaya, 0)

	for _, d := range reqFormat.Data {
		var detID uuid.UUID
		if d.ID == "" {
			detID, _ = uuid.NewV4()
		} else {
			detID, _ = uuid.FromString(d.ID)
		}

		newDetail := PerjalananDinasBiaya{
			ID:                detID,
			IDBpdPegawai:      d.IDBpdPegawai,
			IdPegawai:         d.IdPegawai,
			IDJenisBiaya:      d.IDJenisBiaya,
			IDKomponenBiaya:   d.IDKomponenBiaya,
			Keterangan:        d.Keterangan,
			Qty:               &d.Qty,
			Nominal:           d.Nominal,
			Tanggal:           d.Tanggal,
			IdUangCuci:        d.IdUangCuci,
			NominalUangCuci:   d.NominalUangCuci,
			IdBiayaPenginapan: d.IdBiayaPenginapan,
			NominalPenginapan: d.NominalPenginapan,
			IdTiketBerangkat:  d.IdTiketBerangkat,
			IdTiketPulang:     d.IdTiketPulang,
			CreatedAt:         &now,
			CreatedBy:         &userId,
		}

		details = append(details, newDetail)
	}

	newBpd = details
	return
}

func (p *PerjalananDinasBiaya) SoftDelete(userID uuid.UUID) {
	var now = time.Now()
	p.IsDeleted = true
	p.UpdatedAt = &now
	p.UpdatedBy = &userID
}

type DocPenyelesaianBpdPegawai struct {
	ID       string `db:"id" json:"id"`
	File     string `db:"file" json:"file"`
	IsRevisi *bool  `db:"is_revisi" json:"isRevisi"`
}

type UangMukaBpd struct {
	ID             string  `db:"id" json:"id"` //  idBpdPegawai
	IsUm           bool    `db:"is_um" json:"isUm"`
	PersentaseUm   float64 `db:"persentase_um" json:"persentaseUm"`
	PersentaseSisa float64 `db:"persentase_sisa" json:"persentaseSisa"`
	ShowUm         bool    `db:"show_um" json:"showUm"`
	ShowSisa       bool    `db:"show_sisa" json:"showSisa"`
	TotalUm        float64 `db:"total_um" json:"totalUm"`
	SisaUm         float64 `db:"sisa_um" json:"sisaUm"`
}

func (p *PerjalananDinasBiaya) BiayaRequest(reqFormat PerjalananDinasBiayaDetail, userId uuid.UUID) (newBpd PerjalananDinasBiaya, err error) {
	var now = time.Now()
	var detID uuid.UUID
	fmt.Println("id", reqFormat.ID)
	if reqFormat.ID == "null" {
		detID, _ = uuid.NewV4()
		newBpd = PerjalananDinasBiaya{
			ID:              detID,
			IDBpdPegawai:    reqFormat.IDBpdPegawai,
			IDJenisBiaya:    reqFormat.IDJenisBiaya,
			IDKomponenBiaya: reqFormat.IDKomponenBiaya,
			Nominal:         reqFormat.Nominal,
			CreatedAt:       &now,
			CreatedBy:       &userId,
			IdPegawai:       reqFormat.IdPegawai,
			File:            reqFormat.File,
			IsReimbursement: reqFormat.IsReimbursement,
			Keterangan:      reqFormat.Keterangan,
			Tanggal:         reqFormat.Tanggal,
			Aksi:            "CREATE",
		}
	} else {
		detID, _ = uuid.FromString(reqFormat.ID)
		newBpd = PerjalananDinasBiaya{
			ID:              detID,
			IDBpdPegawai:    reqFormat.IDBpdPegawai,
			IDJenisBiaya:    reqFormat.IDJenisBiaya,
			IDKomponenBiaya: reqFormat.IDKomponenBiaya,
			Nominal:         reqFormat.Nominal,
			UpdatedAt:       &now,
			UpdatedBy:       &userId,
			IdPegawai:       reqFormat.IdPegawai,
			File:            reqFormat.File,
			IsReimbursement: reqFormat.IsReimbursement,
			Keterangan:      reqFormat.Keterangan,
			Tanggal:         reqFormat.Tanggal,
			Aksi:            "UPDATE",
		}

	}

	return
}

func (p *PerjalananDinasBiaya) SoftDeleteBiaya(userId string) {
	now := time.Now()
	p.IsDeleted = true
	p.DeletedBy = &userId
	p.DeletedAt = &now
	p.Aksi = "DELETE"
}
