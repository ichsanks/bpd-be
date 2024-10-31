package auth

import (
	"time"

	"github.com/gofrs/uuid"
)

type Dashboard struct {
	JmlPegawai string `db:"jml_pegawai" json:"jmlPegawai"`
}
type DashboardBpd struct {
	JmlBelumProsesPengajuan *string `db:"belum_proses_pengajuan" json:"jmlBelumProsesPengajuan"`
	JmlPengajuan            *string `db:"pengajuan" json:"jmlPengajuan"`
	JmlSedangDinas          *string `db:"sedang_dinas" json:"jmlSedangDinas"`
	JmlDalamPenyelesaian    *string `db:"dalam_penyelesaian" json:"jmlDalamPenyelesaian"`
}

type DashboardSppd struct {
	JmlBelumProsesPengajuan *string `db:"belum_proses_pengajuan" json:"jmlBelumProsesPengajuan"`
	JmlPengajuan            *string `db:"pengajuan" json:"jmlPengajuan"`
	JmlPengajuanDisetujui   *string `db:"pengajuan_disetujui" json:"jmlPengajuanDisetujui"`
	JmlRevisi               *string `db:"revisi" json:"jmlRevisi"`
}
type JumlahSppd struct {
	JmlPengajuan *string `db:"pengajuan" json:"jmlPengajuan"`
}
type DashboardRequest struct {
	IdPegawai         string `json:"idPegawai"`
	IdPegawaiApproval string `json:"idPegawaiApproval"`
	IdBidang          string `json:"idBidang"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
}

type DataAktifBpd struct {
	ID                   *string `db:"id" json:"id"`
	IdPerjalananDinas    *string `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	Nomor                *string `db:"nomor" json:"nomor"`
	NamaBpd              *string `db:"nama_bpd" json:"namaBpd"`
	Tujuan               *string `db:"tujuan" json:"tujuan"`
	Keperluan            *string `db:"keperluan" json:"keperluan"`
	TglBerangkat         *string `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali           *string `db:"tgl_kembali" json:"tglKembali"`
	Status               *string `db:"status" json:"status"`
	StatusBpd            *string `db:"status_bpd" json:"statusBpd"`
	JenisPerjalananDinas *string `db:"jenis_perjalanan_dinas" json:"jenisPerjalananDinas"`
	JenisKendaraan       *string `db:"jenis_kendaraan" json:"jenisKendaraan"`
}
type DataAktifSppd struct {
	ID                     uuid.UUID `db:"id" json:"id"`
	TglSurat               string    `db:"tgl_surat" json:"tglSurat"`
	NomorSurat             string    `db:"nomor_surat" json:"nomorSurat"`
	IdPegawai              string    `db:"id_pegawai" json:"idPegawai"`
	JenisTujuan            string    `db:"jenis_tujuan" json:"jenisTujuan"`
	TujuanDinas            string    `db:"tujuan_dinas" json:"tujuanDinas"`
	KeperluanDinas         string    `db:"keperluan_dinas" json:"keperluanDinas"`
	TglBerangkat           string    `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali             string    `db:"tgl_kembali" json:"tglKembali"`
	IdFasilitasTransport   string    `db:"id_fasilitas_transport" json:"idFasilitasTransport"`
	IsRombongan            bool      `db:"is_rombongan" json:"isRombongan"`
	IdBranch               *string   `db:"id_branch" json:"idBranch"`
	Nip                    *string   `db:"nip" json:"nip"`
	NamaPegawai            *string   `db:"nama_pegawai" json:"namaPegawai"`
	IdJabatan              *string   `db:"id_jabatan" json:"idJabatan"`
	IdPersonGrade          *string   `db:"id_person_grade" json:"idPersonGrade"`
	Jabatan                *string   `db:"jabatan" json:"jabatan"`
	Kode                   *string   `db:"kode" json:"kode"`
	PersonGrade            *string   `db:"person_grade" json:"personGrade"`
	NamaJenisTujuan        string    `db:"nama_jenis_tujuan" json:"namaJenisTujuan"`
	Keterangan             string    `db:"keterangan" json:"keterangan"`
	NamaFasilitasTransport string    `db:"nama_fasilitas_transport" json:"namaFasilitasTransport"`
	KodeBranch             *string   `db:"kode_branch" json:"kodeBranch"`
	NamaBranch             *string   `db:"nama_branch" json:"namaBranch"`
	IsAntar                *bool     `db:"is_antar" json:"isAntar"`
	IsJemput               *bool     `db:"is_jemput" json:"isJemput"`
	OperasionalHariDinas   *bool     `db:"operasional_hari_dinas" json:"operasionalHariDinas"`
	IdRuleApproval         *string   `db:"id_rule_approval" json:"idRuleApproval"`
	NamaApproval           *string   `db:"nama_approval" json:"namaApproval"`
	TypeApproval           *string   `db:"type_approval" json:"typeApproval"`
	IdJenisApproval        *string   `db:"id_jenis_approval" json:"idJenisApproval"`
	IsPengajuan            *bool     `db:"is_pengajuan" json:"isPengajuan"`
	CreatedAt              time.Time `db:"created_at" json:"createdAt"`
	TenantId               *string   `db:"tenant_id" json:"tenantId"`
	JenisSppd              *string   `db:"jenis_sppd" json:"jenisSppd"`
	IdPegawaiPengaju       *string   `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	Status                 *string   `db:"status" json:"status"`
	IsMaxPengajuan         *bool     `db:"is_max_pengajuan" json:"isMaxPengajuan"`
	IsMaxPenyelesaian      *bool     `db:"is_max_penyelesaian" json:"isMaxPenyelesaian"`
	IdPengajuanBpdHistori  *string   `db:"id_pengajuan_bpd_histori" json:"idPengajuanBpdHistori"`
	StatusBpd              *string   `db:"status_bpd" json:"statusBpd"`
	IdPegawaiApproval      *string   `db:"id_pegawai_approval" json:"idPegawaiApproval"`
	Esign                  *bool     `db:"esign" json:"esign"`
	IsDeleted              bool      `db:"is_deleted" json:"isDeleted"`
	IdLevelBod             *string   `db:"id_level_bod" json:"idLevelBod"`
	File                   *string   `db:"file" json:"file"`
	LinkFile               *string   `db:"link_file" json:"Linkfile"`
	NamaFungsionalitas     *string   `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
}

type DataAktifBpdNew struct {
	ID                       string     `db:"id" json:"id"`
	Nomor                    *string    `db:"nomor" json:"nomorSurat"`
	Nama                     *string    `db:"nama" json:"namaJenisTujuan"`
	Tujuan                   *string    `db:"tujuan" json:"tujuanDinas"`
	Keperluan                *string    `db:"keperluan" json:"keperluanDinas"`
	TglBerangkat             *string    `db:"tgl_berangkat" json:"tglBerangkat"`
	TglKembali               *string    `db:"tgl_kembali" json:"tglKembali"`
	IdJenisPerjalananDinas   *string    `db:"id_jenis_perjalanan_dinas" json:"idJenisPerjalananDinas"`
	IdJenisKendaraan         *string    `db:"id_jenis_kendaraan" json:"idJenisKendaraan"`
	IsRombongan              *bool      `db:"is_rombongan" json:"isRombongan"`
	Status                   *string    `db:"status" json:"status"`
	NamaJenisPerjalananDinas *string    `db:"nama_jenis_perjalanan_dinas" json:"namaJenisPerjalananDinas"`
	NamaJenisKendaraan       *string    `db:"nama_jenis_kendaraan" json:"namaJenisKendaraan"`
	NamaKendaraan            *string    `db:"nama_kendaraan" json:"namaKendaraan"`
	NamaApproval             *string    `db:"nama_approval" json:"namaApproval"`
	IdRuleApproval           *string    `db:"id_rule_approval" json:"idRuleApproval"`
	IdBpdPegawai             *string    `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IdPegawai                *string    `db:"id_pegawai" json:"idPegawai"`
	Nip                      *string    `db:"nip" json:"nip"`
	NamaPegawai              *string    `db:"nama_pegawai" json:"namaPegawai"`
	IdBidang                 *string    `db:"id_bidang" json:"idBidang"`
	IdPengajuanBpdHistori    *string    `db:"id_pengajuan_bpd_histori" json:"idPengajuanBpdHistori"`
	StatusBpd                *string    `db:"status_bpd" json:"statusBpd"`
	TypeApproval             *string    `db:"type_approval" json:"typeApproval"`
	IsMaxPengajuan           *bool      `db:"is_max_pengajuan" json:"isMaxPengajuan"`
	IsMaxPenyelesaian        *bool      `db:"is_max_penyelesaian" json:"isMaxPenyelesaian"`
	CreatedAt                *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy                *string    `db:"created_by" json:"createdBy"`
	UpdatedAt                *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy                *string    `db:"updated_by" json:"updatedBy"`
	IsDeleted                *bool      `db:"is_deleted" json:"isDeleted"`
	IdJenisApproval          *string    `db:"id_jenis_approval" json:"idJenisApproval"`
	IdPegawaiPengaju         *string    `db:"id_pegawai_pengaju" json:"idPegawaiPengaju"`
	IsPengaju                *bool      `db:"is_pengaju" json:"isPengaju"`
	File                     *string    `db:"file" json:"file"`
	IdSppd                   *string    `db:"id_sppd" json:"idSppd"`
	SppId                    *int       `db:"spp_id" json:"sppId"`
	KeteranganTujuan         *string    `db:"keterangan_tujuan" json:"keteranganTujuan"`
	IdPegawaiApproval        *string    `db:"id_pegawai_approval" json:"idPegawaiApproval"`
	Esign                    *bool      `db:"esign" json:"esign"`
	PilihKendaraan           *bool      `db:"pilih_kendaraan" json:"pilihKendaraan"`
	IdJenisTujuan            *string    `db:"id_jenis_tujuan" json:"idJenisTujuan"`
	KetTujuan                *string    `db:"ket_tujuan" json:"ketTujuan"`
	NamaFungsionalitas       *string    `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	FileSppd                 *string    `db:"file_sppd" json:"fileSppd"`
	JenisSppd                *string    `db:"jenis_sppd" json:"jenisSppd"`
}
