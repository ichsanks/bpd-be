package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type RuleApproval struct {
	ID               uuid.UUID            `db:"id" json:"id"`
	Nama             string               `db:"nama" json:"nama"`
	IdFungsionalitas string               `db:"id_fungsionalitas" json:"idFungsionalitas"`
	Jenis            string               `db:"jenis" json:"jenis"`
	CreatedAt        time.Time            `db:"created_at" json:"createdAt"`
	CreatedBy        *string              `db:"created_by" json:"createdBy"`
	UpdatedAt        *time.Time           `db:"updated_at" json:"updatedAt"`
	UpdatedBy        *string              `db:"updated_by" json:"updatedBy"`
	IsDeleted        bool                 `db:"is_deleted" json:"isDeleted"`
	GroupRule        int                  `db:"group_rule" json:"groupRule"`
	IdPegawai        *string              `db:"id_pegawai" json:"idPegawai"`
	TenantID         *uuid.UUID           `db:"tenant_id" json:"tenantId"`
	IdBranch         *string              `db:"id_branch" json:"idBranch"`
	Detail           []RuleApprovalDetail `db:"-" json:"detail"`
}

type RuleApprovalDetail struct {
	ID               uuid.UUID `db:"id" json:"id"`
	IdRuleApproval   string    `db:"id_rule_approval" json:"idRuleApproval"`
	IdFungsionalitas string    `db:"id_fungsionalitas" json:"idFungsionalitas"`
	IdUnor           *string   `db:"id_unor" json:"idUnor"`
	IdBidang         *string   `db:"id_bidang" json:"idBidang"`
	TypeApproval     string    `db:"type_approval" json:"typeApproval"`
	Urut             int       `db:"urut" json:"urut"`
	GroupApproval    int       `db:"group_approval" json:"groupApproval"`
	FeedbackTolak    *string   `db:"feedback_tolak" json:"feedbackTolak"`
	ApprovalLine     *int      `db:"approval_line" json:"approvalLine"`
	IdPegawai        *string   `db:"id_pegawai" json:"idPegawai"`
	Esign            *bool     `db:"esign" json:"esign"`
	KetTtd           *string   `db:"ket_ttd" json:"ketTttd"`
}

type RuleApprovalDTO struct {
	ID                 uuid.UUID               `db:"id" json:"id"`
	Nama               string                  `db:"nama" json:"nama"`
	IdFungsionalitas   *string                 `db:"id_fungsionalitas" json:"idFungsionalitas"`
	NamaFungsionalitas *string                 `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	Jenis              *string                 `db:"jenis" json:"jenis"`
	CreatedAt          time.Time               `db:"created_at" json:"createdAt"`
	CreatedBy          *string                 `db:"created_by" json:"createdBy"`
	UpdatedAt          *time.Time              `db:"updated_at" json:"updatedAt"`
	UpdatedBy          *string                 `db:"updated_by" json:"updatedBy"`
	IsDeleted          bool                    `db:"is_deleted" json:"isDeleted"`
	GroupRule          int                     `db:"group_rule" json:"groupRule"`
	IdPegawai          *string                 `db:"id_pegawai" json:"idPegawai"`
	NamaGroupRule      *string                 `db:"nama_group_rule" json:"namaGroupRule"`
	GroupValue         *string                 `db:"group_value" json:"groupValue"`
	NamaJenis          *string                 `db:"nama_jenis" json:"namaJenis"`
	IdBranch           *string                 `db:"id_branch" json:"idBranch"`
	Detail             []RuleApprovalDetailDTO `db:"-" json:"detail"`
}

type RuleApprovalDetailDTO struct {
	ID                 uuid.UUID `db:"id" json:"id"`
	IdRuleApproval     string    `db:"id_rule_approval" json:"idRuleApproval"`
	IdFungsionalitas   string    `db:"id_fungsionalitas" json:"idFungsionalitas"`
	NamaFungsionalitas *string   `db:"nama_fungsionalitas" json:"namaFungsionalitas"`
	IdUnor             *string   `db:"id_unor" json:"idUnor"`
	KodeUnor           *string   `db:"kode_unor" json:"kodeUnor"`
	NamaUnor           *string   `db:"nama_unor" json:"namaUnor"`
	IdBidang           *string   `db:"id_bidang" json:"idBidang"`
	KodeBidang         *string   `db:"kode_bidang" json:"kodeBidang"`
	NamaBidang         *string   `db:"nama_bidang" json:"namaBidang"`
	TypeApproval       string    `db:"type_approval" json:"typeApproval"`
	Urut               int       `db:"urut" json:"urut"`
	GroupApproval      int       `db:"group_approval" json:"groupApproval"`
	FeedbackTolak      *string   `db:"feedback_tolak" json:"feedbackTolak"`
	ApprovalLine       *int      `db:"approval_line" json:"approvalLine"`
	IdPegawai          *string   `db:"id_pegawai" json:"idPegawai"`
	Pegawai            *string   `db:"pegawai" json:"pegawai"`
	IsHead             *bool     `db:"is_head" json:"isHead"`
	Esign              *bool     `db:"esign" json:"esign"`
	KetTtd             *string   `db:"ket_ttd" json:"ketTtd"`
}

type RuleApprovalRequest struct {
	ID               uuid.UUID                   `db:"id" json:"id"`
	Nama             string                      `db:"nama" json:"nama"`
	IdFungsionalitas string                      `db:"id_fungsionalitas" json:"idFungsionalitas"`
	Jenis            string                      `db:"jenis" json:"jenis"`
	GroupRule        int                         `db:"group_rule" json:"groupRule"`
	IdPegawai        *string                     `db:"id_pegawai" json:"idPegawai"`
	IdBranch         *string                     `db:"id_branch" json:"idBranch"`
	Detail           []RuleApprovalDetailRequest `db:"-" json:"detail"`
}

type RuleApprovalDetailRequest struct {
	ID               string  `db:"id" json:"id"`
	IdRuleApproval   string  `db:"id_rule_approval" json:"idRuleApproval"`
	IdFungsionalitas string  `db:"id_fungsionalitas" json:"idFungsionalitas"`
	IdUnor           *string `db:"id_unor" json:"idUnor"`
	IdBidang         *string `db:"id_bidang" json:"idBidang"`
	TypeApproval     string  `db:"type_approval" json:"typeApproval"`
	Urut             int     `db:"urut" json:"urut"`
	GroupApproval    int     `db:"group_approval" json:"groupApproval"`
	FeedbackTolak    *string `db:"feedback_tolak" json:"feedbackTolak"`
	ApprovalLine     *int    `db:"approval_line" json:"approvalLine"`
	IdPegawai        *string `db:"id_pegawai" json:"idPegawai"`
	Esign            *bool   `db:"esign" json:"esign"`
	KetTtd           *string `db:"ket_ttd" json:"ketTtd"`
}

type RuleApprovalTtd struct {
	ID          string  `db:"id" json:"id"`
	IdPegawai   *string `db:"id_pegawai" json:"idPegawai"`
	Nip         *string `db:"nip" json:"nip"`
	NamaPegawai *string `db:"nama_pegawai" json:"namaPegawai"`
	NamaBidang  *string `db:"nama_bidang" json:"namaBidang"`
	NamaJabatan *string `db:"nama_jabatan" json:"namaJabatan"`
}

func (s *RuleApproval) NewRuleApprovalFormat(reqFormat RuleApprovalRequest, userID string, tenantId uuid.UUID) (rule RuleApproval, err error) {
	newID, _ := uuid.NewV4()
	now := time.Now()

	if reqFormat.ID == uuid.Nil {
		rule = RuleApproval{
			ID:               newID,
			Nama:             reqFormat.Nama,
			IdFungsionalitas: reqFormat.IdFungsionalitas,
			Jenis:            reqFormat.Jenis,
			GroupRule:        reqFormat.GroupRule,
			IdPegawai:        reqFormat.IdPegawai,
			TenantID:         &tenantId,
			IdBranch:         reqFormat.IdBranch,
			CreatedAt:        time.Now(),
			CreatedBy:        &userID,
		}
	} else {
		rule = RuleApproval{
			ID:               reqFormat.ID,
			Nama:             reqFormat.Nama,
			IdFungsionalitas: reqFormat.IdFungsionalitas,
			Jenis:            reqFormat.Jenis,
			GroupRule:        reqFormat.GroupRule,
			IdPegawai:        reqFormat.IdPegawai,
			TenantID:         &tenantId,
			IdBranch:         reqFormat.IdBranch,
			UpdatedAt:        &now,
			UpdatedBy:        &userID,
		}
	}

	details := make([]RuleApprovalDetail, 0)
	for _, d := range reqFormat.Detail {
		var detID uuid.UUID
		if d.ID == "" {
			detID, _ = uuid.NewV4()
		} else {
			detID, _ = uuid.FromString(d.ID)
		}

		newDetail := RuleApprovalDetail{
			ID:               detID,
			IdRuleApproval:   rule.ID.String(),
			IdFungsionalitas: d.IdFungsionalitas,
			IdUnor:           d.IdUnor,
			IdBidang:         d.IdBidang,
			TypeApproval:     d.TypeApproval,
			Urut:             d.Urut,
			GroupApproval:    d.GroupApproval,
			FeedbackTolak:    d.FeedbackTolak,
			ApprovalLine:     d.ApprovalLine,
			IdPegawai:        d.IdPegawai,
			Esign:            d.Esign,
			KetTtd:           d.KetTtd,
		}

		details = append(details, newDetail)
	}

	rule.Detail = details

	return
}

var ColumnMappRuleApproval = map[string]interface{}{
	"id":            "r.id",
	"nama":          "r.nama",
	"namaGroupRule": "r.nama_group_rule",
	"groupValue":    "r.group_value",
	"jenis":         "r.jenis",
	"namaJenis":     "r.nama_jenis",
	"createdBy":     "r.created_by",
	"createdAt":     "r.created_at",
	"updatedBy":     "r.updated_by",
	"updatedAt":     "r.updated_at",
	"isDeleted":     "r.is_deleted",
}

func (rule *RuleApproval) SoftDelete(userId string) {
	now := time.Now()
	rule.IsDeleted = true
	rule.UpdatedBy = &userId
	rule.UpdatedAt = &now
}

type RuleParams struct {
	Jenis            string `json:"jenis"`
	TypeApproval     string `json:"typeApproval"`
	IdPegawai        string `json:"idPegawai"`
	IdUnor           string `json:"idUnor"`
	IdBidang         string `json:"idBidang"`
	IdFungsionalitas string `json:"idFungsionalitas"`
	GroupRule        int    `json:"groupRule"`
}

type RuleDetailParams struct {
	IdPegawai        string `json:"idPegawai"`
	IdApprovalLine   string `json:"idApprovalLine"`
	IdManager        string `json:"idManager"`
	IdFungsionalitas string `json:"idFungsionalitas"`
	IdUnor           string `json:"idUnor"`
	IdBidang         string `json:"idBidang"`
	GroupApproval    int    `json:"groupRule"`
	ApprovalLine     int    `json:"approvalLine"`
}
