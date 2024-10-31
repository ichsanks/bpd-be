package bpd

import (
	"time"

	"github.com/gofrs/uuid"
)

type SuratPerjalananDinas struct {
	ID                   uuid.UUID                     `db:"id" json:"id"`
	TglSurat             string                        `db:"tgl_surat" json:"tglSurat"`
	NomorSurat           string                        `db:"nomor_surat" json:"nomorSurat"`
	IdPegawai            string                        `db:"id_pegawai" json:"idPegawai"`
	JenisTujuan          string                        `db:"jenis_tujuan" json:"jenisTujuan"`
	TujuanDinas          string                        `db:"tujuan_dinas" json:"tujuanDinas"`
	KeperluanDinas       string                        `db:"keperluan_dinas" json:"keperluanDinas"`
	TglBerangkat         string                        `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali           string                        `db:"tgl_kembali" json:"tglKembali"`
	IdFasilitasTransport string                        `db:"id_fasilitas_transport" json:"idFasilitasTransport"`
	IsRombongan          bool                          `db:"is_rombongan" json:"isRombongan"`
	IdBranch             string                        `db:"id_branch" json:"idBranch"`
	IsAntar              bool                          `db:"is_antar" json:"isAntar"`
	IsJemput             bool                          `db:"is_jemput" json:"isJemput"`
	OperasionalHariDinas bool                          `db:"operasional_hari_dinas" json:"operasionalHariDinas"`
	CreatedAt            time.Time                     `db:"created_at" json:"createdAt"`
	CreatedBy            *string                       `db:"created_by" json:"createdBy"`
	UpdatedAt            *time.Time                    `db:"updated_at" json:"updatedAt"`
	UpdatedBy            *string                       `db:"updated_by" json:"updatedBy"`
	IsDeleted            bool                          `db:"is_deleted" json:"isDeleted"`
	TenantId             *string                       `db:"tenant_id" json:"tenantId"`
	JenisSppd            *string                       `db:"jenis_sppd" json:"jenisSppd"`
	IdPegawaiPengaju     *string                       `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	IdRuleApproval       *string                       `db:"id_rule_approval" json:"idRuleApproval"`
	TypeApproval         *string                       `db:"type_approval" json:"typeApproval"`
	IdJenisApproval      *string                       `db:"id_jenis_approval" json:"idJenisApproval"`
	IsPengajuan          *bool                         `db:"is_pengajuan" json:"isPengajuan"`
	Status               *string                       `db:"status" json:"status"`
	File                 *string                       `db:"file" json:"file"`
	LinkFile             *string                       `db:"link_file" json:"linkFile"`
	Detail               []SuratPerjalananDinasPegawai `db:"-" json:"detail"`
	DetailDokumen        []SppdDokumen                 `db:"-" json:"detailDokumen"`
	DetailBiaya          []DetailBiaya                 `db:"-" json:"detailBiaya"`
}

type SuratPerjalananDinasDto struct {
	ID                     uuid.UUID                        `db:"id" json:"id"`
	TglSurat               string                           `db:"tgl_surat" json:"tglSurat"`
	NomorSurat             string                           `db:"nomor_surat" json:"nomorSurat"`
	IdPegawai              string                           `db:"id_pegawai" json:"idPegawai"`
	JenisTujuan            string                           `db:"jenis_tujuan" json:"jenisTujuan"`
	TujuanDinas            string                           `db:"tujuan_dinas" json:"tujuanDinas"`
	KeperluanDinas         string                           `db:"keperluan_dinas" json:"keperluanDinas"`
	TglBerangkat           string                           `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali             string                           `db:"tgl_kembali" json:"tglKembali"`
	IdFasilitasTransport   string                           `db:"id_fasilitas_transport" json:"idFasilitasTransport"`
	IsRombongan            bool                             `db:"is_rombongan" json:"isRombongan"`
	IdBranch               *string                          `db:"id_branch" json:"idBranch"`
	Nip                    *string                          `db:"nip" json:"nip"`
	NamaPegawai            *string                          `db:"nama_pegawai" json:"namaPegawai"`
	IdJabatan              *string                          `db:"id_jabatan" json:"idJabatan"`
	IdPersonGrade          *string                          `db:"id_person_grade" json:"idPersonGrade"`
	IdLevelBod             *string                          `db:"id_level_bod" json:"idLevelBod"`
	Jabatan                *string                          `db:"jabatan" json:"jabatan"`
	Bidang                 *string                          `db:"bidang" json:"bidang"`
	Kode                   *string                          `db:"kode" json:"kode"`
	PersonGrade            *string                          `db:"person_grade" json:"personGrade"`
	NamaJenisTujuan        string                           `db:"nama_jenis_tujuan" json:"namaJenisTujuan"`
	Keterangan             string                           `db:"keterangan" json:"keterangan"`
	NamaFasilitasTransport string                           `db:"nama_fasilitas_transport" json:"namaFasilitasTransport"`
	KodeBranch             *string                          `db:"kode_branch" json:"kodeBranch"`
	NamaBranch             *string                          `db:"nama_branch" json:"namaBranch"`
	IsAntar                *bool                            `db:"is_antar" json:"isAntar"`
	IsJemput               *bool                            `db:"is_jemput" json:"isJemput"`
	OperasionalHariDinas   *bool                            `db:"operasional_hari_dinas" json:"operasionalHariDinas"`
	IdRuleApproval         *string                          `db:"id_rule_approval" json:"idRuleApproval"`
	NamaApproval           *string                          `db:"nama_approval" json:"namaApproval"`
	TypeApproval           *string                          `db:"type_approval" json:"typeApproval"`
	IdJenisApproval        *string                          `db:"id_jenis_approval" json:"idJenisApproval"`
	IsPengajuan            *bool                            `db:"is_pengajuan" json:"isPengajuan"`
	CreatedAt              time.Time                        `db:"created_at" json:"createdAt"`
	TenantId               *string                          `db:"tenant_id" json:"tenantId"`
	JenisSppd              *string                          `db:"jenis_sppd" json:"jenisSppd"`
	IdPegawaiPengaju       *string                          `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	Status                 *string                          `db:"status" json:"status"`
	IsMaxPengajuan         *bool                            `db:"is_max_pengajuan" json:"isMaxPengajuan"`
	IsMaxPenyelesaian      *bool                            `db:"is_max_penyelesaian" json:"isMaxPenyelesaian"`
	IdPengajuanBpdHistori  *string                          `db:"id_pengajuan_bpd_histori" json:"idPengajuanBpdHistori"`
	StatusBpd              *string                          `db:"status_bpd" json:"statusBpd"`
	IdPegawaiApproval      *string                          `db:"id_pegawai_approval" json:"idPegawaiApproval"`
	Esign                  *bool                            `db:"esign" json:"esign"`
	IsDeleted              bool                             `db:"is_deleted" json:"isDeleted"`
	File                   *string                          `db:"file" json:"file"`
	LinkFile               *string                          `db:"link_file" json:"linkFile"`
	NamaFungsionalitas     *string                          `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	Detail                 []SuratPerjalananDinasPegawaiDto `db:"-" json:"detail"`
	DetailDokumen          []SppdDokumenDto                 `db:"-" json:"detailDokumen"`
	DetailBiaya            []DetailBiaya                    `db:"-" json:"detailBiaya"`
}

type SuratPerjalananDinasPegawai struct {
	ID                     uuid.UUID  `db:"id" json:"id"`
	IdSuratPerjalananDinas string     `db:"id_surat_perjalanan_dinas" json:"idSuratPerjalananDinas"`
	IdPegawai              string     `db:"id_pegawai" json:"idPegawai"`
	Status                 *string    `db:"status" json:"status"`
	CreatedAt              time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy              *string    `db:"created_by" json:"createdBy"`
	UpdatedAt              *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy              *string    `db:"updated_by" json:"updatedBy"`
	IsDeleted              bool       `db:"is_deleted" json:"isDeleted"`
}

type SuratPerjalananDinasPegawaiDto struct {
	ID                     uuid.UUID `db:"id" json:"id"`
	IdSuratPerjalananDinas string    `db:"id_surat_perjalanan_dinas" json:"idSuratPerjalananDinas"`
	IdPegawai              string    `db:"id_pegawai" json:"idPegawai"`
	Nip                    *string   `db:"nip" json:"nip"`
	NamaPegawai            *string   `db:"nama_pegawai" json:"namaPegawai"`
	IdJabatan              *string   `db:"id_jabatan" json:"idJabatan"`
	IdPersonGrade          *string   `db:"id_person_grade" json:"idPersonGrade"`
	Jabatan                *string   `db:"jabatan" json:"jabatan"`
	Bidang                 *string   `db:"bidang" json:"bidang"`
	Kode                   string    `db:"kode" json:"kode"`
	PersonGrade            string    `db:"person_grade" json:"personGrade"`
	IdLevelBod             *string   `db:"id_level_bod" json:"idLevelBod"`
}

type SppdDokumenDto struct {
	ID          uuid.UUID `db:"id" json:"id"`
	IdSppd      string    `db:"id_sppd" json:"idSppd"`
	File        string    `db:"file" json:"file"`
	Nama        *string   `db:"nama" json:"nama"`
	Keterangan  *string   `db:"keterangan" json:"keterangan"`
	IdDokumen   string    `db:"id_dokumen" json:"idDokumen"`
	IsMandatory *bool     `db:"is_mandatory" json:"isMandatory"`
}

type SuratPerjalananDinasRequest struct {
	ID                   uuid.UUID                     `db:"id" json:"id"`
	TglSurat             string                        `db:"tgl_surat" json:"tglSurat"`
	NomorSurat           string                        `db:"nomor_surat" json:"nomorSurat"`
	IdPegawai            string                        `db:"id_pegawai" json:"idPegawai"`
	JenisTujuan          string                        `db:"jenis_tujuan" json:"jenisTujuan"`
	TujuanDinas          string                        `db:"tujuan_dinas" json:"tujuanDinas"`
	KeperluanDinas       string                        `db:"keperluan_dinas" json:"keperluanDinas"`
	TglBerangkat         string                        `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali           string                        `db:"tgl_kembali" json:"tglKembali"`
	IdFasilitasTransport string                        `db:"id_fasilitas_transport" json:"idFasilitasTransport"`
	IsRombongan          bool                          `db:"is_rombongan" json:"isRombongan"`
	IdBranch             string                        `db:"id_branch" json:"idBranch"`
	IsAntar              bool                          `db:"is_antar" json:"isAntar"`
	IsJemput             bool                          `db:"is_jemput" json:"isJemput"`
	OperasionalHariDinas bool                          `db:"operasional_hari_dinas" json:"operasionalHariDinas"`
	CreatedAt            time.Time                     `db:"created_at" json:"createdAt"`
	CreatedBy            *string                       `db:"created_by" json:"createdBy"`
	UpdatedAt            *time.Time                    `db:"updated_at" json:"updatedAt"`
	UpdatedBy            *string                       `db:"updated_by" json:"updatedBy"`
	IsDeleted            bool                          `db:"is_deleted" json:"isDeleted"`
	TenantId             *string                       `db:"tenant_id" json:"tenantId"`
	JenisSppd            *string                       `db:"jenis_sppd" json:"jenisSppd"`
	IdPegawaiPengaju     *string                       `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	IdRuleApproval       *string                       `db:"id_rule_approval" json:"idRuleApproval"`
	TypeApproval         *string                       `db:"type_approval" json:"typeApproval"`
	IdJenisApproval      *string                       `db:"id_jenis_approval" json:"idJenisApproval"`
	IsPengajuan          *bool                         `db:"is_pengajuan" json:"isPengajuan"`
	Status               *string                       `db:"status" json:"status"`
	Detail               []SuratPerjalananDinasPegawai `db:"-" json:"detail"`
	DetailDokumen        []SppdDokumen                 `db:"-" json:"detailDokumen"`
}

type UpdateDokumenPendukung struct {
	DetailDokumen []SppdDokumen `db:"-" json:"detailDokumen"`
}

type FilterDetailSPPD struct {
	ID                string `json:"id"`
	IdPegawai         string `json:"idPegawai"`
	IdPegawaiApproval string `json:"idPegawaiApproval"`
	TypeApproval      string `json:"typeApproval"`
	IdBidang          string `json:"idBidang"`
}

type FilesSuratPerjalananDinas struct {
	ID   string `db:"id" json:"id"`
	File string `db:"file" json:"file"`
}

type LinkFilesSuratPerjalananDinas struct {
	ID       string `db:"id" json:"id"`
	LinkFile string `db:"link_file" json:"linkFile"`
}
type ResponseDataNomor struct {
	Status     string `json:"status"`
	NomorSurat string `json:"nomorsurat"`
}

type NomorSppd struct {
	Status     string `json:"status"`
	NomorSurat string `json:"nomorsurat"`
}

type SuratPerjalananDinasListDto struct {
	ID                     uuid.UUID                        `db:"id" json:"id"`
	TglSurat               string                           `db:"tgl_surat" json:"tglSurat"`
	NomorSurat             string                           `db:"nomor_surat" json:"nomorSurat"`
	IdPegawai              string                           `db:"id_pegawai" json:"idPegawai"`
	JenisTujuan            string                           `db:"jenis_tujuan" json:"jenisTujuan"`
	TujuanDinas            string                           `db:"tujuan_dinas" json:"tujuanDinas"`
	KeperluanDinas         string                           `db:"keperluan_dinas" json:"keperluanDinas"`
	TglBerangkat           string                           `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali             string                           `db:"tgl_kembali" json:"tglKembali"`
	IdFasilitasTransport   string                           `db:"id_fasilitas_transport" json:"idFasilitasTransport"`
	IsRombongan            bool                             `db:"is_rombongan" json:"isRombongan"`
	IdBranch               *string                          `db:"id_branch" json:"idBranch"`
	Nip                    *string                          `db:"nip" json:"nip"`
	NamaPegawai            *string                          `db:"nama_pegawai" json:"namaPegawai"`
	IdJabatan              *string                          `db:"id_jabatan" json:"idJabatan"`
	IdPersonGrade          *string                          `db:"id_person_grade" json:"idPersonGrade"`
	Jabatan                *string                          `db:"jabatan" json:"jabatan"`
	Bidang                 *string                          `db:"bidang" json:"bidang"`
	Kode                   *string                          `db:"kode" json:"kode"`
	PersonGrade            *string                          `db:"person_grade" json:"personGrade"`
	NamaJenisTujuan        string                           `db:"nama_jenis_tujuan" json:"namaJenisTujuan"`
	Keterangan             string                           `db:"keterangan" json:"keterangan"`
	NamaFasilitasTransport string                           `db:"nama_fasilitas_transport" json:"namaFasilitasTransport"`
	KodeBranch             *string                          `db:"kode_branch" json:"kodeBranch"`
	NamaBranch             *string                          `db:"nama_branch" json:"namaBranch"`
	IsAntar                *bool                            `db:"is_antar" json:"isAntar"`
	IsJemput               *bool                            `db:"is_jemput" json:"isJemput"`
	OperasionalHariDinas   *bool                            `db:"operasional_hari_dinas" json:"operasionalHariDinas"`
	IdRuleApproval         *string                          `db:"id_rule_approval" json:"idRuleApproval"`
	NamaApproval           *string                          `db:"nama_approval" json:"namaApproval"`
	TypeApproval           *string                          `db:"type_approval" json:"typeApproval"`
	IdJenisApproval        *string                          `db:"id_jenis_approval" json:"idJenisApproval"`
	IsPengajuan            *bool                            `db:"is_pengajuan" json:"isPengajuan"`
	CreatedAt              time.Time                        `db:"created_at" json:"createdAt"`
	TenantId               *string                          `db:"tenant_id" json:"tenantId"`
	JenisSppd              *string                          `db:"jenis_sppd" json:"jenisSppd"`
	IdPegawaiPengaju       *string                          `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	Status                 *string                          `db:"status" json:"status"`
	IsMaxPengajuan         *bool                            `db:"is_max_pengajuan" json:"isMaxPengajuan"`
	IsMaxPenyelesaian      *bool                            `db:"is_max_penyelesaian" json:"isMaxPenyelesaian"`
	IdPengajuanBpdHistori  *string                          `db:"id_pengajuan_bpd_histori" json:"idPengajuanBpdHistori"`
	StatusBpd              *string                          `db:"status_bpd" json:"statusBpd"`
	IdPegawaiApproval      *string                          `db:"id_pegawai_approval" json:"idPegawaiApproval"`
	Esign                  *bool                            `db:"esign" json:"esign"`
	IsDeleted              bool                             `db:"is_deleted" json:"isDeleted"`
	IdLevelBod             *string                          `db:"id_level_bod" json:"idLevelBod"`
	File                   *string                          `db:"file" json:"file"`
	LinkFile               *string                          `db:"link_file" json:"linkFile"`
	NamaFungsionalitas     *string                          `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	Detail                 []SuratPerjalananDinasPegawaiDto `db:"-" json:"detail"`
	DetailDokumen          []SppdDokumenDto                 `db:"-" json:"detailDokumen"`
}

func (s *SuratPerjalananDinas) NewSuratPerjalananDinasFormat(reqFormat SuratPerjalananDinasRequest, userID string) (pd SuratPerjalananDinas, err error) {
	newID, _ := uuid.NewV4()
	now := time.Now()

	if reqFormat.ID == uuid.Nil {
		pd = SuratPerjalananDinas{
			ID:                   newID,
			TglSurat:             reqFormat.TglSurat,
			NomorSurat:           reqFormat.NomorSurat,
			IdPegawai:            reqFormat.IdPegawai,
			JenisTujuan:          reqFormat.JenisTujuan,
			TujuanDinas:          reqFormat.TujuanDinas,
			KeperluanDinas:       reqFormat.KeperluanDinas,
			TglBerangkat:         reqFormat.TglBerangkat,
			TglKembali:           reqFormat.TglKembali,
			IdFasilitasTransport: reqFormat.IdFasilitasTransport,
			IsRombongan:          reqFormat.IsRombongan,
			IdBranch:             reqFormat.IdBranch,
			IsAntar:              reqFormat.IsAntar,
			IsJemput:             reqFormat.IsJemput,
			OperasionalHariDinas: reqFormat.OperasionalHariDinas,
			TenantId:             reqFormat.TenantId,
			JenisSppd:            reqFormat.JenisSppd,
			IdPegawaiPengaju:     reqFormat.IdPegawaiPengaju,
			IdRuleApproval:       reqFormat.IdRuleApproval,
			TypeApproval:         reqFormat.TypeApproval,
			IdJenisApproval:      reqFormat.IdJenisApproval,
			IsPengajuan:          reqFormat.IsPengajuan,
			Status:               reqFormat.Status,
			CreatedAt:            time.Now(),
			CreatedBy:            &userID,
		}
	} else {

		pd = SuratPerjalananDinas{
			ID:                   reqFormat.ID,
			TglSurat:             reqFormat.TglSurat,
			NomorSurat:           reqFormat.NomorSurat,
			IdPegawai:            reqFormat.IdPegawai,
			JenisTujuan:          reqFormat.JenisTujuan,
			TujuanDinas:          reqFormat.TujuanDinas,
			KeperluanDinas:       reqFormat.KeperluanDinas,
			TglBerangkat:         reqFormat.TglBerangkat,
			TglKembali:           reqFormat.TglKembali,
			IdFasilitasTransport: reqFormat.IdFasilitasTransport,
			IsRombongan:          reqFormat.IsRombongan,
			IdBranch:             reqFormat.IdBranch,
			IsAntar:              reqFormat.IsAntar,
			IsJemput:             reqFormat.IsJemput,
			OperasionalHariDinas: reqFormat.OperasionalHariDinas,
			TenantId:             reqFormat.TenantId,
			JenisSppd:            reqFormat.JenisSppd,
			IdPegawaiPengaju:     reqFormat.IdPegawaiPengaju,
			IdRuleApproval:       reqFormat.IdRuleApproval,
			TypeApproval:         reqFormat.TypeApproval,
			IdJenisApproval:      reqFormat.IdJenisApproval,
			IsPengajuan:          reqFormat.IsPengajuan,
			Status:               reqFormat.Status,
			UpdatedAt:            &now,
			UpdatedBy:            &userID,
		}
	}

	details := make([]SuratPerjalananDinasPegawai, 0)
	for _, d := range reqFormat.Detail {
		var detID uuid.UUID
		if d.ID == uuid.Nil {
			detID, _ = uuid.NewV4()
		} else {
			detID = d.ID
		}

		newDetail := SuratPerjalananDinasPegawai{
			ID:                     detID,
			IdSuratPerjalananDinas: pd.ID.String(),
			IdPegawai:              d.IdPegawai,
			Status:                 d.Status,
			CreatedAt:              time.Now(),
			CreatedBy:              &userID,
		}

		details = append(details, newDetail)
	}

	pd.Detail = details
	return
}

func (pd *SuratPerjalananDinas) SoftDelete(userId string) {
	now := time.Now()
	pd.IsDeleted = true
	pd.UpdatedBy = &userId
	pd.UpdatedAt = &now
}

var ColumnMappSuratPerjalananDinasDto = map[string]interface{}{
	"id":                     "id",
	"tglSurat":               "tgl_surat",
	"nomorSurat":             "nomor_surat",
	"idPegawai":              "id_pegawai",
	"jenisTujuan":            "jenis_tujuan",
	"tujuanDinas":            "tujuan_dinas",
	"keperluanDinas":         "keperluan_dinas",
	"tglBerangkat":           "tgl_berangkat",
	"tglKembali":             "tgl_kembali",
	"idFasilitasTransport":   "id_fasilitas_transport",
	"isRombongan":            "is_rombongan",
	"idBranch":               "id_branch",
	"nip":                    "nip",
	"namaPegawai":            "nama_pegawai",
	"idJabatan":              "id_jabatan",
	"idPersonGrade":          "id_person_grade",
	"jabatan":                "jabatan",
	"kode":                   "kode",
	"personGrade":            "person_grade",
	"namaJenisTujuan":        "nama_jenis_tujuan",
	"keterangan":             "keterangan",
	"namaFasilitasTransport": "nama_fasilitas_transport",
	"kodeBranch":             "kode_branch",
	"namaBranch":             "nama_branch",
	"isAntar":                "is_antar",
	"isJemput":               "is_jemput",
	"operasionalHariDinas":   "operasional_hari_dinas",
	"createdAt":              "created_at",
	"tenantId":               "tenant_id",
	"jenisSppd":              "jenis_sppd",
	"idPegawaiPengaju":       "id_pegawai_pengaju",
	"idRuleApproval":         "id_rule_approval",
	"typeApproval":           "type_approval",
	"idJenisApproval":        "id_jenis_approval",
	"isPengajuan":            "is_pengajuan",
	"status":                 "status",
}

type DetailBiaya struct {
	Tanggal             string `db:"tanggal" json:"tanggal"`
	BiayaTransport      *int64 `db:"biaya_transport" json:"biayaTransport"`
	BiayaPerdiem        int64  `db:"biaya_perdiem" json:"biayaPerdiem"`
	BiayaTransportLokal *int64 `db:"biaya_transport_lokal" json:"biayaTransportLokal"`
}
