package report

import "encoding/json"

type ReportRekapBpd struct {
	ID                *string  `db:"id" json:"id"`
	IdPerjalananDinas *string  `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	Nomor             *string  `db:"nomor" json:"nomor"`
	Nip               *string  `db:"nip" json:"nip"`
	Nama              *string  `db:"nama" json:"nama"`
	TglBerangkat      *string  `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali        *string  `db:"tgl_kembali" json:"tglKembali"`
	NamaUnor          *string  `db:"nama_unor" json:"namaUnor"`
	NamaBidang        *string  `db:"nama_bidang" json:"namaBidang"`
	NamaJabatan       *string  `db:"nama_jabatan" json:"namaJabatan"`
	KodeGolongan      *string  `db:"kode_golongan" json:"kodeGolongan"`
	NamaGolongan      *string  `db:"nama_golongan" json:"namaGolongan"`
	NamaDinas         *string  `db:"nama_dinas" json:"namaDinas"`
	Tujuan            *string  `db:"tujuan" json:"tujuan"`
	Keperluan         *string  `db:"keperluan" json:"keperluan"`
	BiayaDinas        *float64 `db:"biaya_dinas" json:"biayaDinas"`
	Reimbursement     *float64 `db:"reimbursement" json:"reimbursement"`
	Akomodasi         *float64 `db:"akomodasi" json:"akomodasi"`
}

type ReportRekapBpdBagian struct {
	Id             *string  `db:"id" json:"id"`
	Nomor          *string  `db:"nomor" json:"nomor"`
	NamaBpd        *string  `db:"nama_bpd" json:"namaBpd"`
	Tujuan         *string  `db:"tujuan" json:"tujuan"`
	Keperluan      *string  `db:"keperluan" json:"keperluan"`
	TglBerangkat   *string  `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali     *string  `db:"tgl_kembali" json:"tglKembali"`
	IsRombongan    *bool    `db:"is_rombongan" json:"isRombongan"`
	JenisBpd       *string  `db:"jenis_bpd" json:"jenisBpd"`
	JenisKendaraan *string  `db:"jenis_kendaraan" json:"jenisKendaraan"`
	NamaBidang     *string  `db:"nama_bidang" json:"namaBidang"`
	Total          *float64 `db:"total" json:"total"`
}

type ReportRekapTotalBpd struct {
	IdBidang           string   `db:"id_bidang" json:"idBidang"`
	Bidang             string   `db:"bidang" json:"bidang"`
	IdUnor             string   `db:"id_unor" json:"idUnor"`
	Unor               string   `db:"unor" json:"unor"`
	JmlPengajuan       *int     `db:"jml_pengajuan" json:"jmlPengajuan"`
	JmlBlmPenyelesaian *int     `db:"jml_blm_penyelesaian" json:"jmlBlmPenyelesaian"`
	JmlSdhPenyelesaian *int     `db:"jml_sdh_penyelesaian" json:"jmlSdhPenyelesaian"`
	Total              *float64 `db:"total" json:"total"`
}

type GroupRekapTotalBpd struct {
	IdBidang string                `db:"id_bidang" json:"idBidang"`
	Bidang   string                `db:"bidang" json:"Bidang"`
	Total    float64               `db:"total" json:"total"`
	Detail   []ReportRekapTotalBpd `db:"detail" json:"detail"`
}
type FilterReport struct {
	TglAwal  string `json:"tglAwal"`
	TglAkhir string `json:"tglAkhir"`
	IdBidang string `json:"idBidang"`
	IdUnor   string `json:"idUnor"`
	Type     string `json:"type"`
}

type ReportRekapAkReim struct {
	IdPerjalananDinas *string  `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdPegawai         *string  `db:"id_pegawai" json:"idPegawai"`
	Nip               *string  `db:"nip" json:"nip"`
	Nama              *string  `db:"nama" json:"nama"`
	Nomor             *string  `db:"nomor" json:"nomor"`
	NamaBpd           *string  `db:"nama_bpd" json:"namaBpd"`
	Tujuan            *string  `db:"tujuan" json:"tujuan"`
	Keperluan         *string  `db:"keperluan" json:"keperluan"`
	TglBerangkat      *string  `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali        *string  `db:"tgl_kembali" json:"tglKembali"`
	IsRombongan       *bool    `db:"is_rombongan" json:"isRombongan"`
	IdBidang          *string  `db:"id_bidang" json:"idBidang"`
	NamaBidang        *string  `db:"nama_bidang" json:"namaBidang"`
	NamaJabatan       *string  `db:"nama_jabatan" json:"namaJabatan"`
	Status            *string  `db:"status" json:"status"`
	JmlHari           *int     `db:"jml_hari" json:"jmlHari"`
	Total             *float64 `db:"total" json:"total"`
}

type RekapBiayaAkomodasi struct {
	ID                *string          `db:"id" json:"id"`
	IdPerjalananDinas *string          `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	Nomor             *string          `db:"nomor" json:"nomor"`
	Nip               *string          `db:"nip" json:"nip"`
	Nama              *string          `db:"nama" json:"nama"`
	TglBerangkat      *string          `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali        *string          `db:"tgl_kembali" json:"tglKembali"`
	NamaBidang        *string          `db:"nama_bidang" json:"namaBidang"`
	NamaJabatan       *string          `db:"nama_jabatan" json:"namaJabatan"`
	KodeGolongan      *string          `db:"kode_golongan" json:"kodeGolongan"`
	NamaGolongan      *string          `db:"nama_golongan" json:"namaGolongan"`
	NamaDinas         *string          `db:"nama_dinas" json:"namaDinas"`
	Tujuan            *string          `db:"tujuan" json:"tujuan"`
	Keperluan         *string          `db:"keperluan" json:"keperluan"`
	TotalAkomodasi    *float64         `db:"total_akomodasi" json:"totalAkomodasi"`
	Details           *json.RawMessage `db:"details" json:"details"`
}

type RekapBiayaAkomodasiDet struct {
	ID           *string  `db:"id" json:"id"`
	IDJenisBiaya *string  `db:"idJenisBiaya" json:"idJenisBiaya"`
	NamaBiaya    *string  `db:"namaBiaya" json:"namaBiaya"`
	Nominal      *float64 `db:"nominal" json:"nominal"`
	Keterangan   *string  `db:"keterangan" json:"keterangan"`
}
