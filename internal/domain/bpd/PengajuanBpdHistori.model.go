package bpd

import (
	"time"

	"github.com/gofrs/uuid"
)

type PengajuanBpdHistori struct {
	ID                   uuid.UUID                   `db:"id" json:"id"`
	Tanggal              time.Time                   `db:"tanggal" json:"tanggal"`
	IdPerjalananDinas    string                      `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdPegawai            *string                     `db:"id_pegawai" json:"idPegawai"`
	IdFungsionalitas     string                      `db:"id_fungsionalitas" json:"idFungsionalitas"`
	IdUnor               string                      `db:"id_unor" json:"idUnor"`
	IdRuleApprovalDetail string                      `db:"id_rule_approval_detail" json:"idRuleApprovalDetail"`
	Catatan              *string                     `db:"catatan" json:"catatan"`
	Keterangan           *string                     `db:"keterangan" json:"keterangan"`
	Status               string                      `db:"status" json:"status"`
	TypeApproval         string                      `db:"type_approval" json:"typeApproval"`
	CreatedAt            time.Time                   `db:"created_at" json:"createdAt"`
	CreatedBy            *string                     `db:"created_by" json:"createdBy"`
	ApprovedAt           *time.Time                  `db:"approved_at" json:"approvedAt"`
	ApprovedBy           *string                     `db:"approved_by" json:"approvedBy"`
	IdBpdHistoriRevisi   *string                     `db:"id_bpd_histori_revisi" json:"idBpdHistoriRevisi"`
	IdApprovalLine       *string                     `db:"id_approval_line" json:"idApprovalLine"`
	GroupApproval        *int                        `db:"group_approval" json:"groupApproval"`
	IdBpdPegawai         *string                     `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	Detail               []PengajuanBpdHistoriDetail `db:"-" json:"detail"`
}

type PengajuanBpdHistoriDTO struct {
	ID                   uuid.UUID  `db:"id" json:"id"`
	Tanggal              time.Time  `db:"tanggal" json:"tanggal"`
	IdPerjalananDinas    string     `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdPegawai            *string    `db:"id_pegawai" json:"idPegawai"`
	IdFungsionalitas     string     `db:"id_fungsionalitas" json:"idFungsionalitas"`
	IdUnor               string     `db:"id_unor" json:"idUnor"`
	IdRuleApprovalDetail string     `db:"id_rule_approval_detail" json:"idRuleApprovalDetail"`
	Catatan              *string    `db:"catatan" json:"catatan"`
	Keterangan           *string    `db:"keterangan" json:"keterangan"`
	Status               string     `db:"status" json:"status"`
	TypeApproval         string     `db:"type_approval" json:"typeApproval"`
	CreatedAt            time.Time  `db:"created_at" json:"createdAt"`
	CreatedBy            *string    `db:"created_by" json:"createdBy"`
	ApprovedAt           *time.Time `db:"approved_at" json:"approvedAt"`
	ApprovedBy           *string    `db:"approved_by" json:"approvedBy"`
	Nip                  *string    `db:"nip" json:"nip"`
	NamaPegawai          *string    `db:"nama_pegawai" json:"namaPegawai"`
	NamaFungsionalitas   *string    `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	NamaUnor             *string    `db:"nama_unor" json:"namaUnor"`
	KodeUnor             *string    `db:"kode_unor" json:"kodeUnor"`
	IdBpdHistoriRevisi   *string    `db:"id_bpd_histori_revisi" json:"idBpdHistoriRevisi"`
}
type PengajuanBpdHistoriInputRequest struct {
	ID                   string  `db:"id" json:"id"`
	IdPerjalananDinas    string  `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdBpdPegawai         *string `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IdPegawai            *string `db:"id_pegawai" json:"idPegawai"`
	IdRuleApprovalDetail string  `db:"id_rule_approval_detail" json:"idRuleApprovalDetail"`
	Catatan              *string `db:"catatan" json:"catatan"`
	Keterangan           *string `db:"keterangan" json:"keterangan"`
	Status               string  `db:"status" json:"status"`
	TypeApproval         string  `db:"type_approval" json:"typeApproval"`
	IdBpdHistoriRevisi   *string `db:"id_bpd_histori_revisi" json:"idBpdHistoriRevisi"`
	KodeUnor             *string `db:"kode_unor" json:"kodeUnor"`
	Jenis                string  `db:"jenis" json:"jenis"`
	FeedbackTolak        string  `json:"-"`
	IdApprovalLine       *string `json:"-"`
	IdManager            *string `json:"-"`
}

type PengajuanBpdHistoriApproveRequest struct {
	ID           string  `db:"id" json:"id"`
	IdPegawai    *string `db:"id_pegawai" json:"idPegawai"`
	Catatan      *string `db:"catatan" json:"catatan"`
	Keterangan   *string `db:"keterangan" json:"keterangan"`
	Status       string  `db:"status" json:"status"`
	TypeApproval string  `db:"type_approval" json:"typeApproval"`
}

type PengajuanBpdHistoriRequest struct {
	ID                   string                             `db:"id" json:"id"`
	IdPerjalananDinas    string                             `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdBpdPegawai         *string                            `db:"id_bpd_pegawai" json:"idBpdPegawai"`
	IdPegawai            *string                            `db:"id_pegawai" json:"idPegawai"`
	IdFungsionalitas     string                             `db:"id_fungsionalitas" json:"idFungsionalitas"`
	IdUnor               string                             `db:"id_unor" json:"idUnor"`
	IdBidang             string                             `db:"id_bidang" json:"idBidang"`
	IdRuleApprovalDetail string                             `db:"id_rule_approval_detail" json:"idRuleApprovalDetail"`
	Catatan              *string                            `db:"catatan" json:"catatan"`
	Keterangan           *string                            `db:"keterangan" json:"keterangan"`
	Status               string                             `db:"status" json:"status"`
	TypeApproval         string                             `db:"type_approval" json:"typeApproval"`
	IdBpdHistoriRevisi   *string                            `db:"id_bpd_histori_revisi" json:"idBpdHistoriRevisi"`
	Jenis                string                             `db:"jenis" json:"jenis"`
	IdApprovalLine       *string                            `db:"id_approval_line" json:"idApprovalLine"`
	GroupApproval        *int                               `db:"group_approval" json:"groupApproval"`
	Detail               []PengajuanBpdHistoriDetailRequest `db:"-" json:"detail"`
}

type UnorPegawai struct {
	IdUnor   string `db:"id_unor" json:"idUnor"`
	KodeUnor string `db:"kode_unor" json:"kodeUnor"`
	NamaUnor string `db:"nama_unor" json:"namaUnor"`
}

func (s *PengajuanBpdHistori) NewPengajuanBpdHistoriFormat(reqFormat PengajuanBpdHistoriRequest, userID string) (pd PengajuanBpdHistori, err error) {
	newID, _ := uuid.NewV4()
	now := time.Now()
	pd = PengajuanBpdHistori{
		ID:                   newID,
		Tanggal:              now,
		IdPerjalananDinas:    reqFormat.IdPerjalananDinas,
		IdPegawai:            reqFormat.IdPegawai,
		IdFungsionalitas:     reqFormat.IdFungsionalitas,
		IdUnor:               reqFormat.IdUnor,
		IdRuleApprovalDetail: reqFormat.IdRuleApprovalDetail,
		Catatan:              reqFormat.Catatan,
		Keterangan:           reqFormat.Keterangan,
		Status:               reqFormat.Status,
		TypeApproval:         reqFormat.TypeApproval,
		CreatedAt:            time.Now(),
		CreatedBy:            &userID,
		IdBpdHistoriRevisi:   reqFormat.IdBpdHistoriRevisi,
		IdApprovalLine:       reqFormat.IdApprovalLine,
		GroupApproval:        reqFormat.GroupApproval,
		IdBpdPegawai:         reqFormat.IdBpdPegawai,
	}

	// pengajuan detail
	details := make([]PengajuanBpdHistoriDetail, 0)
	for _, d := range reqFormat.Detail {
		detID, _ := uuid.NewV4()
		newDetail := PengajuanBpdHistoriDetail{
			ID:          detID,
			IdPengajuan: pd.ID.String(),
			IdPegawai:   d.IdPegawai,
			CreatedAt:   now,
			CreatedBy:   &userID,
		}

		details = append(details, newDetail)
	}

	pd.Detail = details

	return
}

func (s *PengajuanBpdHistori) UpdatePengajuanBpdHistoriFormat(reqFormat PengajuanBpdHistoriInputRequest, userID string) {
	now := time.Now()
	s.IdPegawai = reqFormat.IdPegawai
	s.Catatan = reqFormat.Catatan
	s.Keterangan = reqFormat.Keterangan
	s.Status = reqFormat.Status
	s.ApprovedAt = &now
	s.ApprovedBy = &userID
}

type StatusBPD struct {
	ID          string `db:"id" json:"id"`
	Status      string `db:"status" json:"status"`
	IsPengajuan bool   `db:"is_pengajuan" json:"isPengajuan"`
}
type PengajuanBpdHistoriDetail struct {
	ID          uuid.UUID `db:"id" json:"id"`
	IdPengajuan string    `db:"id_pengajuan" json:"idPengajuan"`
	IdPegawai   string    `db:"id_pegawai" json:"idPegawai"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	CreatedBy   *string   `db:"created_by" json:"createdBy"`
}

type PengajuanBpdHistoriDetailRequest struct {
	ID        uuid.UUID `db:"id" json:"id"`
	IdPegawai string    `db:"id_pegawai" json:"idPegawai"`
}
type Timeline struct {
	ID                   uuid.UUID `db:"id" json:"id"`
	Tanggal              *string   `db:"tanggal" json:"tanggal"`
	IdPerjalananDinas    *string   `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IdPegawai            *string   `db:"id_pegawai" json:"idPegawai"`
	IdFungsionalitas     *string   `db:"id_fungsionalitas" json:"idFungsionalitas"`
	IdUnor               *string   `db:"id_unor" json:"idUnor"`
	IdRuleApprovalDetail *string   `db:"id_rule_approval_detail" json:"idRuleApprovalDetail"`
	Catatan              *string   `db:"catatan" json:"catatan"`
	Keterangan           *string   `db:"keterangan" json:"keterangan"`
	Status               *string   `db:"status" json:"status"`
	TypeApproval         *string   `db:"type_approval" json:"typeApproval"`
	CreatedAt            *string   `db:"created_at" json:"createdAt"`
	ApprovedAt           *string   `db:"approved_at" json:"approvedAt"`
	Nip                  *string   `db:"nip" json:"nip"`
	NamaPegawai          *string   `db:"nama_pegawai" json:"namaPegawai"`
	NamaFungsionalitas   *string   `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	NamaUnor             *string   `db:"nama_unor" json:"namaUnor"`
	KodeUnor             *string   `db:"kode_unor" json:"kodeUnor"`
	NamaBidang           *string   `db:"nama_bidang" json:"namaBidang"`
	NamaJabatan          *string   `db:"nama_jabatan" json:"namaJabatan"`
	NamaPengaju          *string   `db:"nama_pengaju" json:"namaPengaju"`
	IsPengaju            *string   `db:"is_pengaju" json:"isPengaju"`
	KetStatus            *string   `db:"ket_status" json:"ketStatus"`
	KetTtd               *string   `db:"ket_ttd" json:"ketTtd"`
}

type VerifikasiEsignRequest struct {
	ID                string `json:"id"`
	IdPengajuan       string `json:"idPengajuan"`
	IdPegawaiApproval string `json:"idPegawaiApproval"`
	Passphrase        string `json:"passphrase"`
}

type BatalBpdRequest struct {
	ID     string `db:"id" json:"id"`
	Status string `db:"status" json:"status"`
}

type StatusRevisi struct {
	IdPerjalananDinas string `db:"id_perjalanan_dinas" json:"idPerjalananDinas"`
	IsRevisi          bool   `db:"is_revisi" json:"isRevisi"`
}
