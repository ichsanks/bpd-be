package bpd

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type PerjalananDinas struct {
	ID                     uuid.UUID                      `db:"id" json:"id"`
	Nomor                  string                         `db:"nomor" json:"nomor"`
	Nama                   string                         `db:"nama" json:"nama"`
	Tujuan                 string                         `db:"tujuan" json:"tujuan"`
	Keperluan              string                         `db:"keperluan" json:"keperluan"`
	TglBerangkat           string                         `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali             string                         `db:"tgl_kembali" json:"tglKembali"`
	IdJenisPerjalananDinas string                         `db:"id_jenis_perjalanan_dinas" json:"idJenisPerjalananDinas"`
	IdJenisKendaraan       string                         `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
	IsRombongan            bool                           `db:"is_rombongan" json:"isRombongan"`
	Status                 string                         `db:"status" json:"status"`
	IdRuleApproval         string                         `db:"id_rule_approval" json:"idRuleApproval"`
	IdJenisApproval        string                         `db:"id_jenis_approval" json:"idJenisApproval"`
	IdPegawai              string                         `db:"id_pegawai" json:"idPegawai"`
	IdPegawaiPengaju       string                         `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	File                   *string                        `db:"file" json:"file"`
	TypeApproval           *string                        `db:"type_approval" json:"typeApproval"`
	IsPengajuan            *bool                          `db:"is_pengajuan" json:"isPengajuan"`
	IdSppd                 *string                        `db:"id_sppd" json:"idSppd"`
	TenantId               *string                        `db:"tenant_id" json:"tenantId"`
	IdBranch               *string                        `db:"id_branch" json:"idBranch"`
	CreatedAt              time.Time                      `db:"created_at" json:"createdAt"`
	CreatedBy              *string                        `db:"created_by" json:"createdBy"`
	UpdatedAt              *time.Time                     `db:"updated_at" json:"updatedAt"`
	UpdatedBy              *string                        `db:"updated_by" json:"updatedBy"`
	IsDeleted              bool                           `db:"is_deleted" json:"isDeleted"`
	Detail                 []PerjalananDinasPegawaiDetail `db:"-" json:"detail"`
}

type PerjalananDinasPegawaiDetail struct {
	ID                uuid.UUID              `db:"id" json:"id"`
	Nomor             string                 `db:"nomor" json:"nomor"`
	IdPerjalananDinas string                 `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdPegawai         string                 `db:"id_pegawai" json:"idPegawai"`
	IsPic             bool                   `db:"is_pic" json:"isPic"`
	Status            string                 `db:"status" json:"status"`
	File              *string                `db:"file" json:"file"`
	CreatedAt         time.Time              `db:"created_at" json:"createdAt"`
	CreatedBy         *string                `db:"created_by" json:"createdBy"`
	UpdatedAt         *time.Time             `db:"updated_at" json:"updatedAt"`
	UpdatedBy         *string                `db:"updated_by" json:"updatedBy"`
	IsDeleted         bool                   `db:"is_deleted" json:"isDeleted"`
	TglBerangkat      *string                `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali        *string                `db:"tgl_kembali" json:"tglKembali"`
	Tujuan            *string                `db:"tujuan" json:"tujuan"`
	Keperluan         *string                `db:"keperluan" json:"keperluan"`
	BpdPenyesuaian    *string                `db:"bpd_penyesuaian" json:"bpdPenyesuaian"`
	IsRevisi          *bool                  `db:"is_revisi" json:"isRevisi"`
	Dokumen           []PerjalananDinasBiaya `db:"-" json:"dokumen"`
}

type PerjalananDinasDokumenDetail struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	IdBpdPegawai    string     `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	File            *string    `db:"file" json:"file"`
	Keterangan      *string    `db:"keterangan" json:"keterangan"`
	IdSyaratDokumen *string    `db:"id_syarat_dokumen" json:"idSyaratDokumen"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy       *string    `db:"created_by" json:"createdBy"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy       *string    `db:"updated_by" json:"updatedBy"`
	IsDeleted       bool       `db:"is_deleted" json:"isDeleted"`
}

type PerjalananDinasDokumenDto struct {
	ID              uuid.UUID `db:"id" json:"id"`
	IdBpdPegawai    string    `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	File            *string   `db:"file" json:"file"`
	Keterangan      *string   `db:"keterangan" json:"keterangan"`
	IdSyaratDokumen *string   `db:"id_syarat_dokumen" json:"idSyaratDokumen"`
	IdDokumen       *string   `db:"id_dokumen" json:"idDokumen"`
	IsMandatory     *bool     `db:"is_mandatory" json:"isMandatory"`
	Nama            *string   `db:"nama" json:"nama"`
}

type PerjalananDinasDTO struct {
	ID                       string                            `db:"id" json:"id"`
	Nomor                    *string                           `db:"nomor" json:"nomor"`
	Nama                     *string                           `db:"nama" json:"nama"`
	Tujuan                   *string                           `db:"tujuan" json:"tujuan"`
	Keperluan                *string                           `db:"keperluan" json:"keperluan"`
	TglBerangkat             *string                           `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali               *string                           `db:"tgl_kembali" json:"tglKembali"`
	IdJenisPerjalananDinas   *string                           `db:"id_jenis_perjalanan_dinas" json:"idJenisPerjalananDinas"`
	IdJenisKendaraan         *string                           `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
	IsRombongan              *bool                             `db:"is_rombongan" json:"isRombongan"`
	Status                   *string                           `db:"status" json:"status"`
	NamaJenisPerjalananDinas *string                           `db:"nama_jenis_perjalanan_dinas" json:"namaJenisPerjalananDinas"`
	NamaJenisKendaraan       *string                           `db:"nama_jenis_kendaraan" json:"namaJenisKendaraan"`
	NamaKendaraan            *string                           `db:"nama_kendaraan" json:"namaKendaraan"`
	NamaApproval             *string                           `db:"nama_approval" json:"namaApproval"`
	IdRuleApproval           *string                           `db:"id_rule_approval" json:"idRuleApproval"`
	IdBpdPegawai             *string                           `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IdPegawai                *string                           `db:"id_pegawai" json:"idPegawai"`
	Nip                      *string                           `db:"nip" json:"nip"`
	NamaPegawai              *string                           `db:"nama_pegawai" json:"namaPegawai"`
	IdBidang                 *string                           `db:"id_bidang" json:"idBidang"`
	IdPengajuanBpdHistori    *string                           `db:"id_pengajuan_bpd_histori" json:"idPengajuanBpdHistori"`
	StatusBpd                *string                           `db:"status_bpd" json:"statusBpd"`
	TypeApproval             *string                           `db:"type_approval" json:"typeApproval"`
	IsMaxPengajuan           *bool                             `db:"is_max_pengajuan" json:"isMaxPengajuan"`
	IsMaxPenyelesaian        *bool                             `db:"is_max_penyelesaian" json:"isMaxPenyelesaian"`
	CreatedAt                *time.Time                        `db:"created_at" json:"createdAt"`
	CreatedBy                *string                           `db:"created_by" json:"createdBy"`
	UpdatedAt                *time.Time                        `db:"updated_at" json:"updatedAt"`
	UpdatedBy                *string                           `db:"updated_by" json:"updatedBy"`
	IsDeleted                *bool                             `db:"is_deleted" json:"isDeleted"`
	Detail                   []PerjalananDinasPegawaiDetailDTO `db:"-" json:"detail"`
	DetailKendaraan          []PerjalananDinasKendaraanDTO     `db:"-" json:"detailKendaraan"`
	DetailLogKegiatan        []LogKegiatan                     `db:"-" json:"detailLogKegiatan"`
	DetailDokumen            []SppdDokumenDto                  `db:"-" json:"detailDokumen"`
	IdJenisApproval          *string                           `db:"id_jenis_approval" json:"idJenisApproval"`
	IdPegawaiPengaju         *string                           `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	IsPengaju                *bool                             `db:"is_pengaju" json:"isPengaju"`
	File                     *string                           `db:"file" json:"file"`
	IdSppd                   *string                           `db:"id_sppd" json:"idSppd"`
	SppId                    *int                              `db:"spp_id" json:"sppId"`
	KeteranganTujuan         *string                           `db:"keterangan_tujuan" json:"keteranganTujuan"`
	IdPegawaiApproval        *string                           `db:"id_pegawai_approval" json:"idPegawaiApproval"`
	Esign                    *bool                             `db:"esign" json:"esign"`
	PilihKendaraan           *bool                             `db:"pilih_kendaraan" json:"pilihKendaraan"`
	IdJenisTujuan            *string                           `db:"id_jenis_tujuan" json:"idJenisTujuan"`
	KetTujuan                *string                           `db:"ket_tujuan" json:"ketTujuan"`
	NamaFungsionalitas       *string                           `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	FileSppd                 *string                           `db:"file_sppd" json:"fileSppd"`
	JenisSppd                *string                           `db:"jenis_sppd" json:"jenisSppd"`
}

type ListBpdNew struct {
	ID                     string                            `db:"id" json:"id"`
	Nomor                  *string                           `db:"nomor" json:"nomor"`
	Nama                   *string                           `db:"nama" json:"nama"`
	Tujuan                 *string                           `db:"tujuan" json:"tujuan"`
	Keperluan              *string                           `db:"keperluan" json:"keperluan"`
	TglBerangkat           *string                           `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali             *string                           `db:"tgl_kembali" json:"tglKembali"`
	IsRombongan            *bool                             `db:"is_rombongan" json:"isRombongan"`
	Status                 *string                           `db:"status" json:"status"`
	IdBpdPegawai           *string                           `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IdPegawaiDetail        *string                           `db:"id_pegawai_detail" json:"idPegawaiDetail"`
	IdPegawai              *string                           `db:"id_pegawai" json:"idPegawai"`
	Nip                    *string                           `db:"nip" json:"nip"`
	NamaPegawai            *string                           `db:"nama_pegawai" json:"namaPegawai"`
	IdRuleApproval         *string                           `db:"id_rule_approval" json:"idRuleApproval"`
	IdJenisPerjalananDinas *string                           `db:"id_jenis_perjalanan_dinas" json:"idJenisPerjalananDinas"`
	IdJenisKendaraan       *string                           `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
	CreatedAt              *time.Time                        `db:"created_at" json:"createdAt"`
	CreatedBy              *string                           `db:"created_by" json:"createdBy"`
	UpdatedAt              *time.Time                        `db:"updated_at" json:"updatedAt"`
	UpdatedBy              *string                           `db:"updated_by" json:"updatedBy"`
	IsDeleted              *bool                             `db:"is_deleted" json:"isDeleted"`
	Detail                 []PerjalananDinasPegawaiDetailDTO `db:"-" json:"detail"`
	DetailKendaraan        []PerjalananDinasKendaraanDTO     `db:"-" json:"detailKendaraan"`
	DetailLogKegiatan      []LogKegiatan                     `db:"-" json:"detailLogKegiatan"`
	DetailDokumen          []PerjalananDinasDokumen          `db:"-" json:"detailDokumen"`
	IdJenisApproval        *string                           `db:"id_jenis_approval" json:"idJenisApproval"`
	IdPegawaiPengaju       *string                           `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	File                   *string                           `db:"file" json:"file"`
	IdSppd                 string                            `db:"id_sppd" json:"idSppd"`
	KeteranganTujuan       *string                           `db:"keterangan_tujuan" json:"keteranganTujuan"`
	JenisSppd              *string                           `db:"jenis_sppd" json:"jenisSppd"`
}

type BpdPegawaiDTO struct {
	ID                      string                            `db:"id" json:"id"`
	IdPerjalananDinas       string                            `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	Nomor                   *string                           `db:"nomor" json:"nomor"`
	NamaBpd                 *string                           `db:"nama_bpd" json:"namaBpd"`
	Tujuan                  *string                           `db:"tujuan" json:"tujuan"`
	Keperluan               *string                           `db:"keperluan" json:"keperluan"`
	TglBerangkat            *string                           `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali              *string                           `db:"tgl_kembali" json:"tglKembali"`
	IdJenisPerjalananDinas  *string                           `db:"id_jenis_perjalanan_dinas" json:"idJenisPerjalananDinas"`
	IdJenisKendaraan        *string                           `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
	IsRombongan             *bool                             `db:"is_rombongan" json:"isRombongan"`
	Status                  *string                           `db:"status" json:"status"`
	JenisPerjalananDinas    *string                           `db:"jenis_perjalanan_dinas" json:"jenisPerjalananDinas"`
	JenisKendaraan          *string                           `db:"jenis_kendaraan" json:"jenisKendaraan"`
	IdPegawai               *string                           `db:"id_pegawai" json:"idPegawai"`
	Nip                     *string                           `db:"nip" json:"nip"`
	NamaPegawai             *string                           `db:"nama_pegawai" json:"namaPegawai"`
	Unor                    *string                           `db:"unor" json:"unor"`
	Bidang                  *string                           `db:"bidang" json:"bidang"`
	Jabatan                 *string                           `db:"jabatan" json:"jabatan"`
	KodeGolongan            *string                           `db:"kode_golongan" json:"kodeGolongan"`
	Golongan                *string                           `db:"golongan" json:"golongan"`
	IdPengajuanBpdHistori   *string                           `db:"id_pengajuan_bpd_histori" json:"idPengajuanBpdHistori"`
	CreatedAt               *time.Time                        `db:"created_at" json:"createdAt"`
	CreatedBy               *string                           `db:"created_by" json:"createdBy"`
	IsDeleted               *bool                             `db:"is_deleted" json:"isDeleted"`
	IdJenisApproval         *string                           `db:"id_jenis_approval" json:"idJenisApproval"`
	IdPegawaiPengaju        *string                           `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	FileBpd                 *string                           `db:"file_bpd" json:"fileBpd"`
	FilePenyelesaianBpd     *string                           `db:"file_penyelesaian_bpd" json:"filePenyelesaianBpd"`
	BpdPenyesuaian          *string                           `db:"bpd_penyesuaian" json:"bpdPenyesuaian"`
	StatusBpd               *string                           `db:"status_bpd" json:"statusBpd"`
	TypeApproval            *string                           `db:"type_approval" json:"typeApproval"`
	TujuanPenyesuaian       *string                           `db:"tujuan_penyesuaian" json:"tujuanPenyesuaian"`
	KeperluanPenyesuaian    *string                           `db:"keperluan_penyesuaian" json:"keperluanPenyesuaian"`
	TglBerangkatPenyesuaian *string                           `db:"tgl_berangkat_penyesuaian" json:"tglBerangkatPenyesuaian"`
	TglKembaliPenyesuaian   *string                           `db:"tgl_kembali_penyesuaian" json:"tglKembaliPenyesuaian"`
	Total                   *float64                          `db:"total" json:"total"`
	IdPegawaiApproval       *string                           `db:"id_pegawai_approval" json:"idPegawaiApproval"`
	Esign                   *bool                             `db:"esign" json:"esign"`
	IsSppb                  *bool                             `db:"is_sppb" json:"isSppb"`
	IsRevisi                *bool                             `db:"is_revisi" json:"isRevisi"`
	IsUm                    *bool                             `db:"is_um" json:"isUm"`
	PersentaseUm            *float64                          `db:"persentase_um" json:"persentaseUm"`
	PersentaseSisa          *float64                          `db:"persentase_sisa" json:"persentaseSisa"`
	ShowUm                  *bool                             `db:"show_um" json:"showUm"`
	ShowSisa                *bool                             `db:"show_sisa" json:"showSisa"`
	TotalUm                 *float64                          `db:"total_um" json:"totalUm"`
	SisaUm                  *float64                          `db:"sisa_um" json:"sisaUm"`
	Detail                  []PerjalananDinasPegawaiDetailDTO `db:"-" json:"detail"`
	DetailKendaraan         []PerjalananDinasKendaraanDTO     `db:"-" json:"detailKendaraan"`
	DetailLogKegiatan       []LogKegiatan                     `db:"-" json:"detailLogKegiatan"`
	DetailDokumen           []PerjalananDinasDokumen          `db:"-" json:"detailDokumen"`
}
type PerjalananDinasPegawaiDetailDTO struct {
	ID                  uuid.UUID                   `db:"id" json:"id"`
	IdPerjalananDinas   string                      `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdSppd              *string                     `db:"id_sppd" json:"idSppd"`
	IdPegawai           string                      `db:"id_pegawai" json:"idPegawai"`
	NamaPegawai         *string                     `db:"nama_pegawai" json:"namaPegawai"`
	Nip                 *string                     `db:"nip" json:"nip"`
	IdBranch            *string                     `db:"id_branch" json:"idBranch"`
	IdFungsionalitas    *string                     `db:"id_fungsionalitas" json:"idFungsionalitas"`
	IdUnor              *string                     `db:"id_unor" json:"idUnor"`
	IdBidang            *string                     `db:"id_bidang" json:"idBidang"`
	NamaUnor            *string                     `db:"nama_unor" json:"namaUnor"`
	NamaBidang          *string                     `db:"nama_bidang" json:"namaBidang"`
	NamaJabatan         *string                     `db:"nama_jabatan" json:"namaJabatan"`
	KodeGolongan        *string                     `db:"kode_golongan" json:"kodeGolongan"`
	NamaGolongan        *string                     `db:"nama_golongan" json:"namaGolongan"`
	TglBerangkat        *string                     `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali          *string                     `db:"tgl_kembali" json:"tglKembali"`
	IsPic               bool                        `db:"is_pic" json:"isPic"`
	Status              string                      `db:"status" json:"status"`
	CreatedAt           time.Time                   `db:"created_at" json:"createdAt"`
	CreatedBy           *string                     `db:"created_by" json:"createdBy"`
	File                []PerjalananDinasDokumenDto `db:"dokumen" json:"dokumen"`
	IdLevelBod          *string                     `db:"id_level_bod" json:"idLevelBod"`
	Tujuan              *string                     `db:"tujuan" json:"tujuan"`
	Keperluan           *string                     `db:"keperluan" json:"keperluan"`
	Nominal             decimal.Decimal             `db:"nominal" json:"nominal"`
	Detail              []BiayaPegawai              `db:"detail" json:"detail"`
	DetailAkomodasi     []BiayaPerjalananDinasDto   `db:"detail_akomodasi" json:"detailAkomodasi"`
	DetailReimbursement []BiayaPerjalananDinasDto   `db:"detail_reimbursement" json:"detailReimbursement"`
	HistoriBiaya        []HistoriPerjalananDinas    `db:"histori_biaya" json:"historiBiaya"`
}
type PerjalananDinasRequest struct {
	ID                     string                                `db:"id" json:"id"`
	Nomor                  string                                `db:"nomor" json:"nomor"`
	Nama                   string                                `db:"nama" json:"nama"`
	Tujuan                 string                                `db:"tujuan" json:"tujuan"`
	Keperluan              string                                `db:"keperluan" json:"keperluan"`
	TglBerangkat           string                                `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali             string                                `db:"tgl_kembali" json:"tglKembali"`
	IdJenisPerjalananDinas string                                `db:"id_jenis_perjalanan_dinas" json:"idJenisPerjalananDinas"`
	IdJenisKendaraan       string                                `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
	IsRombongan            bool                                  `db:"is_rombongan" json:"isRombongan"`
	Status                 string                                `db:"status" json:"status"`
	IdRuleApproval         string                                `db:"id_rule_approval" json:"idRuleApproval"`
	IdJenisApproval        string                                `db:"id_jenis_approval" json:"idJenisApproval"`
	IdPegawai              string                                `db:"id_pegawai" json:"idPegawai"`
	IdPegawaiPengaju       string                                `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	TypeApproval           *string                               `db:"type_approval" json:"typeApproval"`
	IsPengajuan            *bool                                 `db:"is_pengajuan" json:"isPengajuan"`
	IdSppd                 *string                               `db:"id_sppd" json:"idSppd"`
	TenantId               *string                               `db:"tenant_id" json:"tenantId"`
	IdBranch               *string                               `db:"id_branch" json:"idBranch"`
	Detail                 []PerjalananDinasPegawaiDetailRequest `db:"-" json:"detail"`
}

type PerjalananDinasPegawaiDetailRequest struct {
	ID                string                 `db:"id" json:"id"`
	Nomor             string                 `db:"nomor" json:"nomor"`
	IdPerjalananDinas string                 `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdPegawai         string                 `db:"id_pegawai" json:"idPegawai"`
	IsPic             bool                   `db:"is_pic" json:"isPic"`
	Status            string                 `db:"status" json:"status"`
	File              *string                `db:"file" json:"file"`
	TglBerangkat      *string                `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali        *string                `db:"tgl_kembali" json:"tglKembali"`
	Dokumen           []PerjalananDinasBiaya `db:"-" json:"dokumen"`
}

type PerjalananDinasDokumenDetailRequest struct {
	ID           uuid.UUID `db:"id" json:"id"`
	IdBpdPegawai string    `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	File         *string   `db:"file" json:"file"`
	Keterangan   *string   `db:"keterangan" json:"keterangan"`
}

type BiayaPegawai struct {
	IdBpd           string          `db:"id_bpd" json:"idBpd"`
	IdSppd          string          `db:"id_sppd" json:"idSppd"`
	IdPegawai       string          `db:"id_pegawai" json:"idPegawai"`
	NamaPegawai     string          `db:"nama_pegawai" json:"namaPegawai"`
	Tgl             string          `db:"tgl" json:"tgl"`
	TglBerangkat    string          `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali      string          `db:"tgl_kembali" json:"tglKembali"`
	IDKomponenBiaya string          `db:"id_komponen_biaya" json:"idKomponenBiaya"`
	IDJenisBiaya    string          `db:"id_jenis_biaya" json:"idJenisBiaya"`
	Nominal         decimal.Decimal `db:"biaya" json:"nominal"`
	Urut            string          `db:"urut" json:"urut"`
	JenisBiaya      string          `db:"jenis_biaya" json:"jenisBiaya"`
}

type EsignRequest struct {
	Pdf        string  `json:"pdf"`
	Ttd        *string `json:"ttd"`
	FilePath   string  `json:"filePath"`
	TtdPath    string  `json:"ttdPath"`
	Addr       string  `json:"addr"`
	Nik        string  `json:"nik"`
	Passphrase string  `json:"passphrase"`
	// Page       *string `json:"page"`
	// Width      string  `json:"width"`
	// Height     string  `json:"height"`
	// X          *string `json:"x"`
	// Y          *string `json:"y"`
	// Tempat     string  `json:"tempat"`
	// Tagar      string  `json:"tagar"`
}

type ResponseEsign struct {
	IdDokumen        string  `json:"idDokumen"`
	Base64SignedFile *string `json:"base64SignedFile"`
}

func (s *PerjalananDinas) NewPerjalananDinasFormat(reqFormat PerjalananDinasRequest, userID string) (pd PerjalananDinas, err error) {
	newID, _ := uuid.NewV4()
	now := time.Now()

	if reqFormat.ID == "" {
		pd = PerjalananDinas{
			ID:                     newID,
			Nomor:                  reqFormat.Nomor,
			Nama:                   reqFormat.Nama,
			Tujuan:                 reqFormat.Tujuan,
			Keperluan:              reqFormat.Keperluan,
			TglBerangkat:           reqFormat.TglBerangkat,
			TglKembali:             reqFormat.TglKembali,
			IdJenisPerjalananDinas: reqFormat.IdJenisPerjalananDinas,
			IdJenisKendaraan:       reqFormat.IdJenisKendaraan,
			IsRombongan:            reqFormat.IsRombongan,
			Status:                 reqFormat.Status,
			IdRuleApproval:         reqFormat.IdRuleApproval,
			IdJenisApproval:        reqFormat.IdJenisApproval,
			IdPegawai:              reqFormat.IdPegawai,
			IdPegawaiPengaju:       reqFormat.IdPegawaiPengaju,
			TypeApproval:           reqFormat.TypeApproval,
			IsPengajuan:            reqFormat.IsPengajuan,
			IdSppd:                 reqFormat.IdSppd,
			TenantId:               reqFormat.TenantId,
			IdBranch:               reqFormat.IdBranch,
			CreatedAt:              now,
			CreatedBy:              &userID,
		}
	} else {
		id, _ := uuid.FromString(reqFormat.ID)
		pd = PerjalananDinas{
			ID:                     id,
			Nomor:                  reqFormat.Nomor,
			Nama:                   reqFormat.Nama,
			Tujuan:                 reqFormat.Tujuan,
			Keperluan:              reqFormat.Keperluan,
			TglBerangkat:           reqFormat.TglBerangkat,
			TglKembali:             reqFormat.TglKembali,
			IdJenisPerjalananDinas: reqFormat.IdJenisPerjalananDinas,
			IdJenisKendaraan:       reqFormat.IdJenisKendaraan,
			IsRombongan:            reqFormat.IsRombongan,
			Status:                 reqFormat.Status,
			IdRuleApproval:         reqFormat.IdRuleApproval,
			IdJenisApproval:        reqFormat.IdJenisApproval,
			IdPegawai:              reqFormat.IdPegawai,
			IdPegawaiPengaju:       reqFormat.IdPegawaiPengaju,
			TypeApproval:           reqFormat.TypeApproval,
			IsPengajuan:            reqFormat.IsPengajuan,
			IdSppd:                 reqFormat.IdSppd,
			TenantId:               reqFormat.TenantId,
			IdBranch:               reqFormat.IdBranch,
			UpdatedAt:              &now,
			UpdatedBy:              &userID,
		}
	}

	details := make([]PerjalananDinasPegawaiDetail, 0)
	for _, d := range reqFormat.Detail {
		var detID uuid.UUID
		if d.ID == "" {
			detID, _ = uuid.NewV4()
		} else {
			detID, _ = uuid.FromString(d.ID)
		}

		detailsDokumen := make([]PerjalananDinasBiaya, 0)
		for _, e := range d.Dokumen {
			dokId, _ := uuid.NewV4()
			var now = time.Now()
			detID, _ = uuid.FromString(userID)
			newDokumen := PerjalananDinasBiaya{
				ID:              dokId,
				IDBpdPegawai:    detID.String(),
				IDJenisBiaya:    e.IDJenisBiaya,
				IDKomponenBiaya: e.IDKomponenBiaya,
				File:            e.File,
				Keterangan:      e.Keterangan,
				Nominal:         e.Nominal,
				IsReimbursement: false,
				IdPegawai:       e.IdPegawai,
				CreatedAt:       &now,
				CreatedBy:       &detID,
			}

			detailsDokumen = append(detailsDokumen, newDokumen)

		}

		newDetail := PerjalananDinasPegawaiDetail{
			ID:                detID,
			Nomor:             d.Nomor,
			IdPerjalananDinas: pd.ID.String(),
			IdPegawai:         d.IdPegawai,
			IsPic:             d.IsPic,
			Status:            d.Status,
			TglBerangkat:      d.TglBerangkat,
			TglKembali:        d.TglKembali,
			File:              d.File,
			Dokumen:           detailsDokumen,
			CreatedAt:         time.Now(),
			CreatedBy:         &userID,
		}

		details = append(details, newDetail)
	}

	pd.Detail = details

	return
}

var ColumnMappPerjalananDinas = map[string]interface{}{
	"id":                     "pd.id",
	"nomor":                  "pd.nomor",
	"nama":                   "pd.nama",
	"tujuan":                 "pd.tujuan",
	"keperluan":              "pd.keperluan",
	"tglBerangkat":           "pd.tgl_berangkat",
	"tglKembali":             "pd.tgl_kembali",
	"idJenisPerjalananDinas": "pd.id_jenis_perjalanan_dinas",
	"idJenisKendaraan":       "pd.id_jenis_kendaraan",
	"isRombongan":            "pd.is_rombongan",
	"status":                 "pd.status",
	"idRuleApproval":         "pd.id_rule_approval",
	"namaPegawai":            "pd.nama_pegawai",
	"createdAt":              "pd.created_at",
	"createdBy":              "pd.created_by",
	"updatedAt":              "pd.updated_at",
	"updatedBy":              "pd.updated_by",
	"isDeleted":              "pd.is_deleted",
}

var ColumnMappPenyelesaianBpd = map[string]interface{}{
	"id":                     "pp.id",
	"nomor":                  "pp.nomor",
	"namaBpd":                "pp.nama_bpd",
	"tujuan":                 "pp.tujuan",
	"keperluan":              "pp.keperluan",
	"tglBerangkat":           "pp.tgl_berangkat",
	"tglKembali":             "pp.tgl_kembali",
	"idJenisPerjalananDinas": "pp.id_jenis_perjalanan_dinas",
	"idJenisKendaraan":       "pp.id_jenis_kendaraan",
	"isRombongan":            "pp.is_rombongan",
	"status":                 "pp.status",
	"namaPegawai":            "pp.nama_pegawai",
	"unor":                   "pp.unor",
	"bidang":                 "pp.bidang",
	"createdAt":              "pp.created_at",
	"createdBy":              "pp.created_by",
	"isDeleted":              "pp.is_deleted",
}

func (pd *PerjalananDinas) SoftDelete(userId string) {
	now := time.Now()
	pd.IsDeleted = true
	pd.UpdatedBy = &userId
	pd.UpdatedAt = &now
}

type FilesPerjalananDinas struct {
	ID                     string `db:"id" json:"id"`
	File                   string `db:"file" json:"file"`
	IdJenisPerjalananDinas string `db:"id_jenis_perjalanan_dinas" json:"idJenisPerjalananDinas"`
}

type ResponseSpp struct {
	ID    string `db:"id" json:"id"`
	SppId int    `db:"spp_id" json:"sppId"`
}
type FilterDetailBPD struct {
	ID                string `json:"id"`
	IdPegawai         string `json:"idPegawai"`
	IdPegawaiApproval string `json:"idPegawaiApproval"`
	TypeApproval      string `json:"typeApproval"`
}

type BpdPegawaiRequest struct {
	ID           string `db:"id" json:"id"`
	TglBerangkat string `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali   string `db:"tgl_kembali" json:"tglKembali"`
	Tujuan       string `db:"tujuan" json:"tujuan"`
	Keperluan    string `db:"keperluan" json:"keperluan"`
	File         string `db:"file" json:"file"`
}

func (l *PerjalananDinasPegawaiDetail) FormatUpdateBpdPegawai(req BpdPegawaiRequest, userId string) {
	now := time.Now()
	id, _ := uuid.FromString(req.ID)
	l.ID = id
	l.TglBerangkat = &req.TglBerangkat
	l.TglKembali = &req.TglKembali
	l.Tujuan = &req.Tujuan
	l.Keperluan = &req.Keperluan
	l.File = &req.File
	l.UpdatedBy = &userId
	l.UpdatedAt = &now
}
