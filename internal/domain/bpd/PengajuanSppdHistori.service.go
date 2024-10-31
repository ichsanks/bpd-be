package bpd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/internal/domain/master"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
)

type PengajuanSppdHistoriService interface {
	Create(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error)
	CreatePenyelesaian(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error)
	Approve(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error)
	GetTimeline(IdSuratPerjalananDinas string, idBpdPegawai string) (data []TimelineSppd, err error)
	Batal(reqFormat PengajuanSppdHistoriInputRequest, userID string) (err error)
	RevisiPenyelesaianBiaya(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error)
	RevisiPengajuan(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error)
}

type PengajuanSppdHistoriServiceImpl struct {
	PengajuanSppdHistoriRepository PengajuanSppdHistoriRepository
	SuratPerjalananDinasRepository SuratPerjalananDinasRepository
	RuleApprovalRepository         master.RuleApprovalRepository
	PegawaiRepository              master.PegawaiRepository
	PerjalananDinasRepository      PerjalananDinasRepository
	Config                         *configs.Config
}

func ProvidePengajuanSppdHistoriServiceImpl(repository PengajuanSppdHistoriRepository, pdRepository SuratPerjalananDinasRepository, ruleRepository master.RuleApprovalRepository, pgwRepository master.PegawaiRepository, config *configs.Config) *PengajuanSppdHistoriServiceImpl {
	s := new(PengajuanSppdHistoriServiceImpl)
	s.PengajuanSppdHistoriRepository = repository
	s.SuratPerjalananDinasRepository = pdRepository
	s.RuleApprovalRepository = ruleRepository
	s.PegawaiRepository = pgwRepository
	s.Config = config
	return s
}

func (s *PengajuanSppdHistoriServiceImpl) GetTimeline(IdSuratPerjalananDinas string, idBpdPegawai string) (data []TimelineSppd, err error) {
	return s.PengajuanSppdHistoriRepository.GetTimeline(IdSuratPerjalananDinas, idBpdPegawai)
}

// Function Create digunakan untuk ketika melakukan pengajuan awal
func (s *PengajuanSppdHistoriServiceImpl) Create(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error) {
	reqPayload := PengajuanSppdHistoriRequest{
		IdSuratPerjalananDinas: reqFormat.IdSuratPerjalananDinas,
		IdRuleApprovalDetail:   reqFormat.IdRuleApprovalDetail,
		IdBpdPegawai:           reqFormat.IdBpdPegawai,
		Catatan:                reqFormat.Catatan,
		Keterangan:             reqFormat.Keterangan,
		Status:                 reqFormat.Status,
		Jenis:                  reqFormat.Jenis,
		TypeApproval:           reqFormat.TypeApproval,
		TenantId:               reqFormat.TenantId,
		IdBranch:               reqFormat.IdBranch,
	}

	// proses cek pegawai dari perjalanan dinas / bpd pegawai
	var idPegawaiStr string
	if reqFormat.TypeApproval == "PENYELESAIAN" {
		if reqFormat.IdBpdPegawai == nil {
			return PengajuanSppdHistori{}, errors.New("ID BPD Pegawai tidak ditemukan")
		}
		dinas, err := s.PerjalananDinasRepository.ResolveBpdPegawaiByID(*reqFormat.IdBpdPegawai)
		if err != nil {
			return PengajuanSppdHistori{}, errors.New("Data BPD Pegawai tidak ditemukan")
		}
		idPegawaiStr = dinas.IdPegawai
	} else {
		dinas, err := s.SuratPerjalananDinasRepository.ResolveDtoByID(reqFormat.IdSuratPerjalananDinas)
		if err != nil {
			return PengajuanSppdHistori{}, errors.New("Data perjalanan dinas tidak ditemukan")
		}
		idPegawaiStr = model.ParseString(&dinas.IdPegawai)
	}

	idPegawai, _ := uuid.FromString(idPegawaiStr)
	fmt.Println("pegwaiiid", idPegawai)
	pegawai, err := s.PegawaiRepository.ResolveByIDDTO(idPegawai)
	if err != nil {
		return PengajuanSppdHistori{}, errors.New("Data Pegawai tidak ditemukan")
	}

	reqPayload.IdPegawai = &idPegawaiStr
	// cek group rule
	ruleParam := master.RuleParams{
		Jenis:            reqFormat.Jenis,
		TypeApproval:     reqFormat.TypeApproval,
		IdPegawai:        model.ParseString(reqFormat.IdPegawai),
		IdBidang:         model.ParseString(pegawai.IdBidang),
		IdUnor:           model.ParseString(pegawai.IdUnor),
		IdFungsionalitas: model.ParseString(pegawai.IdFungsionalitas),
	}

	checkRule, existRule := s.CheckGroupRule(ruleParam)
	if !existRule {
		return PengajuanSppdHistori{}, errors.New("Rule Approval belum di setting")
	}

	ruleParam.GroupRule = checkRule.GroupRule
	fmt.Println("GroupRule", checkRule.GroupRule)
	rule, _ := s.RuleApprovalRepository.GetAllRuleApprovalDetailByKode(ruleParam)

	// find pegawai
	rule2 := rule[0]
	idUnor := ""
	idBidang := ""
	if rule2.IdUnor == nil && !model.ParseBoolean(rule2.IsHead) {
		idUnor = model.ParseString(pegawai.IdUnor)
	} else {
		idUnor = *rule2.IdUnor
	}

	if rule2.IdBidang != nil {
		idBidang = *rule2.IdBidang
	} else {
		idBidang = model.ParseString(pegawai.IdBidang)
	}

	RDetParamas := master.RuleDetailParams{
		IdPegawai:        model.ParseString(rule2.IdPegawai),
		IdApprovalLine:   model.ParseString(pegawai.IdApprovalLine),
		IdManager:        model.ParseString(pegawai.IdManager),
		ApprovalLine:     model.ParseInt(rule2.ApprovalLine),
		IdFungsionalitas: rule2.IdFungsionalitas,
		IdUnor:           idUnor,
		IdBidang:         idBidang,
		GroupApproval:    rule2.GroupApproval,
	}

	// fmt.Println("IdPegawai:", model.ParseString(rule2.IdPegawai))
	// fmt.Println("IdApprovalLine:", model.ParseString(pegawai.IdApprovalLine))
	// fmt.Println("IdManager:", model.ParseString(pegawai.IdManager))
	// fmt.Println("ApprovalLine:", model.ParseInt(rule2.ApprovalLine))
	// fmt.Println("IdFungsionalitas:", rule2.IdFungsionalitas)
	// fmt.Println("IdUnor:", idUnor)
	// fmt.Println("IdBidang:", idBidang)
	// fmt.Println("GroupApproval:", rule2.GroupApproval)

	pegawaiApp, err := s.GetPegawaiGroupApproval(RDetParamas)

	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	c, _ := json.Marshal(pegawaiApp)
	fmt.Println("next approval : ", string(c))

	if len(pegawaiApp) == 0 {
		return PengajuanSppdHistori{}, errors.New("Pegawai approval tidak ditemukan")
	}

	// insert pengajuan histori awal dari pengaju
	reqPayload.IdRuleApprovalDetail = rule2.ID.String()
	reqPayload.IdFungsionalitas = model.ParseString(pegawai.IdFungsionalitas)
	if reqFormat.Keterangan == nil {
		sk := "1"
		reqPayload.Keterangan = &sk
	}
	reqPayload.IdUnor = model.ParseString(pegawai.IdUnor)
	reqPayload.IdBidang = model.ParseString(pegawai.IdBidang)

	// Detail Employee
	pengajuanDetail := PengajuanSppdHistoriDetailRequest{
		IdPegawai: pegawai.ID.String(),
	}
	reqPayload.Detail = append(reqPayload.Detail, pengajuanDetail)

	data, _ = data.NewPengajuanSppdHistoriFormat(reqPayload, userID)
	err = s.PengajuanSppdHistoriRepository.Create(data)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	// insert next pengajuan ke user verifikator untuk approval
	reqPayload.Detail = make([]PengajuanSppdHistoriDetailRequest, 0)
	reqNext := reqPayload
	reqNext.IdPegawai = nil
	reqNext.Keterangan = nil
	reqNext.IdUnor = model.ParseString(rule2.IdUnor)
	reqNext.IdBidang = model.ParseString(rule2.IdBidang)
	reqNext.IdFungsionalitas = rule2.IdFungsionalitas
	reqNext.GroupApproval = &rule2.GroupApproval

	if rule2.GroupApproval == 1 && model.ParseInt(rule2.ApprovalLine) == 1 {
		reqNext.IdApprovalLine = pegawai.IdApprovalLine
	} else if rule2.GroupApproval == 1 && model.ParseInt(rule2.ApprovalLine) == 2 {
		reqNext.IdApprovalLine = pegawai.IdManager
	}

	for _, d := range pegawaiApp {
		pengajuanDetail := PengajuanSppdHistoriDetailRequest{
			IdPegawai: d.ID.String(),
		}
		reqNext.Detail = append(reqNext.Detail, pengajuanDetail)
	}

	nextPengajuan, _ := data.NewPengajuanSppdHistoriFormat(reqNext, userID)
	err = s.PengajuanSppdHistoriRepository.Create(nextPengajuan)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	// update status perjalanan dinas / bpd pegawai
	if reqFormat.TypeApproval == "PENYELESAIAN" && reqFormat.IdBpdPegawai != nil {
		err = s.PengajuanSppdHistoriRepository.UpdateStatusBpdPegawai(StatusSPPD{
			ID:     *reqFormat.IdBpdPegawai,
			Status: "1",
		})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err = s.PengajuanSppdHistoriRepository.UpdateStatusSPPD(StatusSPPD{
			ID:     reqFormat.IdSuratPerjalananDinas,
			Status: "1",
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return nextPengajuan, nil
}

// Function Create Next Approval
func (s *PengajuanSppdHistoriServiceImpl) CreateNextApproval(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, finish bool, err error) {
	reqPayload := PengajuanSppdHistoriRequest{
		IdSuratPerjalananDinas: reqFormat.IdSuratPerjalananDinas,
		IdBpdPegawai:           reqFormat.IdBpdPegawai,
		IdRuleApprovalDetail:   reqFormat.IdRuleApprovalDetail,
		Catatan:                reqFormat.Catatan,
		Keterangan:             reqFormat.Keterangan,
		Status:                 reqFormat.Status,
		Jenis:                  reqFormat.Jenis,
		TypeApproval:           reqFormat.TypeApproval,
	}

	// cek histori pengajuan approval
	isMaxPengajuan := false
	bpdHistori, _ := s.PengajuanSppdHistoriRepository.GetAll(reqFormat.IdSuratPerjalananDinas, model.ParseString(reqFormat.IdBpdPegawai), reqFormat.TypeApproval)
	countHistori := len(bpdHistori)

	if countHistori == 0 {
		return PengajuanSppdHistori{}, false, errors.New("Pengajuan histori tidak ditemukan")
	}

	// insert next approval
	var ruleApprovalDetail master.RuleApprovalDetailDTO
	if reqFormat.IdRuleApprovalDetail != "" {
		idRuleApprovalDetail := reqFormat.IdRuleApprovalDetail
		ruleApprovalDetail, err = s.RuleApprovalRepository.ResolveRuleApprovalDetailDTO(idRuleApprovalDetail)
		if err != nil {
			isMaxPengajuan = true
		}
	} else {
		idRuleApprovalDetail := bpdHistori[countHistori-1].IdRuleApprovalDetail
		ruleApprovalDetail, err = s.RuleApprovalRepository.GetNextRuleApprovalDetail(idRuleApprovalDetail, reqFormat.TypeApproval)
		if err != nil {
			isMaxPengajuan = true
		}
	}

	// proses cek pegawai dari perjalanan dinas / bpd pegawai
	var idPegawaiStr string
	if reqFormat.TypeApproval == "PENYELESAIAN" {
		if reqFormat.IdBpdPegawai == nil {
			return PengajuanSppdHistori{}, false, errors.New("ID BPD Pegawai tidak ditemukan")
		}
		dinas, err := s.PerjalananDinasRepository.ResolveBpdPegawaiByID(*reqFormat.IdBpdPegawai)
		if err != nil {
			return PengajuanSppdHistori{}, false, errors.New("Data BPD Pegawai tidak ditemukan")
		}
		idPegawaiStr = dinas.IdPegawai
	} else {
		dinas, err := s.SuratPerjalananDinasRepository.ResolveByID(reqFormat.IdSuratPerjalananDinas)
		if err != nil {
			return PengajuanSppdHistori{}, false, errors.New("Data perjalanan dinas tidak ditemukan")
		}
		idPegawaiStr = model.ParseString(&dinas.IdPegawai)
	}

	// dinas, err := s.SuratPerjalananDinasRepository.ResolveByIDDTO(reqFormat.IdSuratPerjalananDinas)
	// if err != nil {
	// 	return PengajuanSppdHistori{}, false, errors.New("Data perjalanan dinas tidak ditemukan")
	// }

	pegawaiID, _ := uuid.FromString(idPegawaiStr)
	pegawai, err := s.PegawaiRepository.ResolveByIDDTO(pegawaiID)
	if err != nil {
		return PengajuanSppdHistori{}, false, errors.New("Data Pegawai tidak ditemukan")
	}

	idUnor := ""
	idBidang := ""
	if ruleApprovalDetail.IdUnor == nil && !model.ParseBoolean(ruleApprovalDetail.IsHead) {
		idUnor = model.ParseString(pegawai.IdUnor)
	} else {
		idUnor = model.ParseString(ruleApprovalDetail.IdUnor)
	}

	if ruleApprovalDetail.IdBidang != nil {
		idBidang = model.ParseString(ruleApprovalDetail.IdBidang)
	} else {
		idBidang = model.ParseString(pegawai.IdBidang)
	}

	if !isMaxPengajuan {
		// Check Group Approval
		ruleParamsDet := master.RuleDetailParams{
			IdPegawai:        model.ParseString(ruleApprovalDetail.IdPegawai),
			IdApprovalLine:   model.ParseString(pegawai.IdApprovalLine),
			IdManager:        model.ParseString(pegawai.IdManager),
			ApprovalLine:     model.ParseInt(ruleApprovalDetail.ApprovalLine),
			IdFungsionalitas: ruleApprovalDetail.IdFungsionalitas,
			IdUnor:           idUnor,
			IdBidang:         idBidang,
			GroupApproval:    ruleApprovalDetail.GroupApproval,
		}

		pegawaiApp, err := s.GetPegawaiGroupApproval(ruleParamsDet)
		if err != nil {
			return PengajuanSppdHistori{}, false, err
		}

		c, _ := json.Marshal(pegawaiApp)
		fmt.Println("next approval : ", string(c))

		if len(pegawaiApp) == 0 {
			return PengajuanSppdHistori{}, false, errors.New("Pegawai approval tidak ditemukan")
		}

		reqPayload.IdRuleApprovalDetail = ruleApprovalDetail.ID.String()
		reqPayload.IdUnor = model.ParseString(ruleApprovalDetail.IdUnor)
		reqPayload.IdBidang = model.ParseString(ruleApprovalDetail.IdBidang)
		reqPayload.IdFungsionalitas = ruleApprovalDetail.IdFungsionalitas
		reqPayload.GroupApproval = &ruleApprovalDetail.GroupApproval

		if ruleApprovalDetail.GroupApproval == 1 && model.ParseInt(ruleApprovalDetail.ApprovalLine) == 1 {
			reqPayload.IdApprovalLine = reqFormat.IdApprovalLine
		} else if ruleApprovalDetail.GroupApproval == 1 && model.ParseInt(ruleApprovalDetail.ApprovalLine) == 2 {
			reqPayload.IdApprovalLine = reqFormat.IdManager
		}
		idPegawaiTes := ""
		for _, d := range pegawaiApp {
			pengajuanDetail := PengajuanSppdHistoriDetailRequest{
				IdPegawai: d.ID.String(),
			}
			reqPayload.Detail = append(reqPayload.Detail, pengajuanDetail)
			idPegawaiTes = d.ID.String()
		}

		data, _ = data.NewPengajuanSppdHistoriFormat(reqPayload, userID)

		if data.IdFungsionalitas == "" {
			id, err := uuid.FromString(idPegawaiTes)
			pegawaiDet, err := s.PegawaiRepository.ResolveByID(id)
			if err != nil {
				return PengajuanSppdHistori{}, false, errors.New("Data Pegawai tidak ditemukan")
			}
			data.IdFungsionalitas = *pegawaiDet.IdFungsionalitas
		}

		err = s.PengajuanSppdHistoriRepository.Create(data)
		if err != nil {
			return PengajuanSppdHistori{}, false, err
		}
	} else {
		// update status perjalanan dinas
		if reqFormat.TypeApproval == "PENGAJUAN_SPPD" {
			StatusSPPD := StatusSPPD{
				ID:     reqFormat.IdSuratPerjalananDinas,
				Status: "2", // pengajuan disetujui
			}
			err = s.PengajuanSppdHistoriRepository.UpdateStatusSPPD(StatusSPPD)
			if err != nil {
				fmt.Println(err)
			}
		} else if reqFormat.TypeApproval == "PENYELESAIAN" {
			StatusSPPD := StatusSPPD{
				ID:     *reqFormat.IdBpdPegawai,
				Status: "5", // proses penyelesaian
			}
			err = s.PengajuanSppdHistoriRepository.UpdateStatusBpdPegawai(StatusSPPD)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return data, isMaxPengajuan, nil
}

func (s *PengajuanSppdHistoriServiceImpl) CreatePrevious(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error) {
	// cek jika memiliki id bpd histori revisi sebelumnya
	idPengajuanHistori := ""
	if reqFormat.IdBpdHistoriRevisi != nil {
		idPengajuanHistori = *reqFormat.IdBpdHistoriRevisi
	} else {
		idPengajuanHistori = reqFormat.ID
	}

	bpdHistori, err := s.PengajuanSppdHistoriRepository.GetPreviousBpdHistori(reqFormat.IdSuratPerjalananDinas, model.ParseString(reqFormat.IdBpdPegawai), idPengajuanHistori)
	if err != nil {
		return PengajuanSppdHistori{}, errors.New("Pengajuan BPD tidak dapat direvisi")
	}

	IDHistori := bpdHistori.ID.String()
	reqPayload := PengajuanSppdHistoriRequest{
		ID:                     reqFormat.ID,
		IdSuratPerjalananDinas: reqFormat.IdSuratPerjalananDinas,
		IdBpdPegawai:           bpdHistori.IdBpdPegawai,
		IdRuleApprovalDetail:   bpdHistori.IdRuleApprovalDetail,
		IdFungsionalitas:       bpdHistori.IdFungsionalitas,
		IdUnor:                 bpdHistori.IdUnor,
		Catatan:                reqFormat.Catatan,
		Keterangan:             bpdHistori.Keterangan,
		Status:                 reqFormat.Status,
		TypeApproval:           reqFormat.TypeApproval,
		IdBpdHistoriRevisi:     &IDHistori,
		GroupApproval:          bpdHistori.GroupApproval,
		IdApprovalLine:         bpdHistori.IdApprovalLine,
	}

	pengajuanDetail := PengajuanSppdHistoriDetailRequest{
		IdPegawai: model.ParseString(bpdHistori.IdPegawai),
	}
	reqPayload.Detail = append(reqPayload.Detail, pengajuanDetail)

	// Save pengajuan
	data, _ = data.NewPengajuanSppdHistoriFormat(reqPayload, userID)
	err = s.PengajuanSppdHistoriRepository.Create(data)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	return data, nil
}

func (s *PengajuanSppdHistoriServiceImpl) Approve(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error) {
	// update status pengajuan bpd
	data, err = s.PengajuanSppdHistoriRepository.ResolveByID(reqFormat.ID)
	if err != nil {
		return PengajuanSppdHistori{}, errors.New("Data pengajuan bpd histori tidak ditemukan")
	}
	bpdHistori := data

	// keterangan jika tidak kosong menandakan yang mengajukan adalah user pengaju awal
	if bpdHistori.Keterangan != nil {
		reqFormat.Keterangan = bpdHistori.Keterangan
		reqFormat.IdRuleApprovalDetail = bpdHistori.IdRuleApprovalDetail
	}

	data.UpdatePengajuanSppdHistoriFormat(reqFormat, userID)
	err = s.PengajuanSppdHistoriRepository.UpdateApproval(data)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	// get rule approval detail
	ruleDetail, err := s.RuleApprovalRepository.ResolveRuleApprovalDetail(bpdHistori.IdRuleApprovalDetail)
	if err != nil {
		return PengajuanSppdHistori{}, errors.New("Data rule approval detail tidak ditemukan")
	}

	if ruleDetail.FeedbackTolak != nil {
		reqFormat.FeedbackTolak = *ruleDetail.FeedbackTolak
	}

	// next pengajuan
	reqFormat.IdSuratPerjalananDinas = data.IdSuratPerjalananDinas
	reqFormat.IdBpdPegawai = data.IdBpdPegawai
	reqFormat.Catatan = nil
	reqFormat.Keterangan = nil

	isUpdateBpd := false
	if reqFormat.Status == "2" { // disetujui
		reqFormat.Status = "1"
		// reqFormat.KodeUnor = dataDTO.KodeUnor
		_, _, err := s.CreateNextApproval(reqFormat, userID)
		if err != nil {
			return PengajuanSppdHistori{}, err
		}

		// update jika kembali ke user pengaju
		if bpdHistori.Keterangan != nil {
			pKet := *bpdHistori.Keterangan
			if pKet == "1" {
				isUpdateBpd = true
			}
		} else {
			isUpdateBpd = false
		}
	} else if reqFormat.Status == "3" { // revisi
		reqFormat.IdBpdHistoriRevisi = bpdHistori.IdBpdHistoriRevisi
		bpdPrevious := PengajuanSppdHistori{}
		if reqFormat.FeedbackTolak == "1" {
			bpdPrevious, err = s.CreatePrevious(reqFormat, userID)
			if err != nil {
				return PengajuanSppdHistori{}, err
			}
			// FEEBEEK TOLAK 2 mungkin perlu diskusikan
		} else {
			StatusSPPD := StatusSPPD{
				ID:     reqFormat.IdSuratPerjalananDinas,
				Status: reqFormat.Status,
			}
			err = s.PengajuanSppdHistoriRepository.UpdateStatusSPPD(StatusSPPD)
			if err != nil {
				return PengajuanSppdHistori{}, err
			}
		}

		if bpdPrevious.Keterangan != nil || bpdHistori.Keterangan != nil {
			isUpdateBpd = true
		} else {
			isUpdateBpd = false
		}
	}

	if isUpdateBpd {
		// update data jika diperlukan misal untuk tolak
		if reqFormat.TypeApproval == "PENYELESAIAN" {
			if reqFormat.IdBpdPegawai == nil {
				return PengajuanSppdHistori{}, errors.New("ID BPD Pegawai tidak ditemukan")
			}

			StatusSPPD := StatusSPPD{
				ID:     *reqFormat.IdBpdPegawai,
				Status: reqFormat.Status,
			}
			err = s.PengajuanSppdHistoriRepository.UpdateStatusBpdPegawai(StatusSPPD)
			if err != nil {
				return PengajuanSppdHistori{}, err
			}
		} else {
			StatusSPPD := StatusSPPD{
				ID:     reqFormat.IdSuratPerjalananDinas,
				Status: reqFormat.Status,
			}
			err = s.PengajuanSppdHistoriRepository.UpdateStatusSPPD(StatusSPPD)
			if err != nil {
				return PengajuanSppdHistori{}, err
			}
		}
	}

	return data, nil
}

// Function Create digunakan untuk ketika melakukan penyelesaian awal
func (s *PengajuanSppdHistoriServiceImpl) CreatePenyelesaian(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error) {
	reqPayload := PengajuanSppdHistoriRequest{
		IdSuratPerjalananDinas: reqFormat.IdSuratPerjalananDinas,
		IdRuleApprovalDetail:   reqFormat.IdRuleApprovalDetail,
		IdBpdPegawai:           reqFormat.IdBpdPegawai,
		Catatan:                reqFormat.Catatan,
		Keterangan:             reqFormat.Keterangan,
		Status:                 reqFormat.Status,
		// Jenis:                reqFormat.Jenis,
		TypeApproval: reqFormat.TypeApproval,
	}

	// Cek jenis approval
	jenis := "1" // default all employee
	jaSplit := strings.Split(reqFormat.Jenis, ",")
	if len(jaSplit) == 1 {
		jenis = reqFormat.Jenis
	} else if len(jaSplit) > 1 {
		for _, s := range jaSplit {
			if s == "1" {
				break
			} else {
				jenis = s
				break
			}
		}
	}
	reqPayload.Jenis = jenis
	fmt.Println("JenisApproval", jenis)

	// proses cek pegawai dari perjalanan dinas / bpd pegawai
	var idPegawaiStr string
	if reqFormat.IdBpdPegawai == nil {
		return PengajuanSppdHistori{}, errors.New("ID BPD Pegawai tidak ditemukan")
	}

	// dinas, err := s.SuratPerjalananDinasRepository.ResolveBpdPegawaiByID(*reqFormat.IdBpdPegawai)
	if err != nil {
		return PengajuanSppdHistori{}, errors.New("Data BPD Pegawai tidak ditemukan")
	}
	// idPegawaiStr = dinas.IdPegawai

	idPegawai, _ := uuid.FromString(idPegawaiStr)
	pegawai, err := s.PegawaiRepository.ResolveByIDDTO(idPegawai)
	if err != nil {
		return PengajuanSppdHistori{}, errors.New("Data Pegawai tidak ditemukan")
	}

	reqPayload.IdPegawai = &idPegawaiStr
	// cek group rule
	ruleParam := master.RuleParams{
		Jenis:            jenis,
		TypeApproval:     reqFormat.TypeApproval,
		IdPegawai:        model.ParseString(reqFormat.IdPegawai),
		IdBidang:         model.ParseString(pegawai.IdBidang),
		IdUnor:           model.ParseString(pegawai.IdUnor),
		IdFungsionalitas: model.ParseString(pegawai.IdFungsionalitas),
	}

	checkRule, existRule := s.CheckGroupRule(ruleParam)
	if !existRule {
		return PengajuanSppdHistori{}, errors.New("Rule Approval belum di setting")
	}

	ruleParam.GroupRule = checkRule.GroupRule
	fmt.Println("GroupRule", checkRule.GroupRule)
	rule, _ := s.RuleApprovalRepository.GetAllRuleApprovalDetailByKode(ruleParam)

	// find pegawai
	rule2 := rule[0]
	idUnor := ""
	idBidang := ""
	if rule2.IdUnor == nil && !model.ParseBoolean(rule2.IsHead) {
		idUnor = model.ParseString(pegawai.IdUnor)
	} else {
		idUnor = *rule2.IdUnor
	}

	if rule2.IdBidang != nil {
		idBidang = *rule2.IdBidang
	} else {
		idBidang = model.ParseString(pegawai.IdBidang)
	}

	RDetParamas := master.RuleDetailParams{
		IdPegawai:        model.ParseString(rule2.IdPegawai),
		IdApprovalLine:   model.ParseString(pegawai.IdApprovalLine),
		IdManager:        model.ParseString(pegawai.IdManager),
		ApprovalLine:     model.ParseInt(rule2.ApprovalLine),
		IdFungsionalitas: rule2.IdFungsionalitas,
		IdUnor:           idUnor,
		IdBidang:         idBidang,
		GroupApproval:    rule2.GroupApproval,
	}
	pegawaiApp, err := s.GetPegawaiGroupApproval(RDetParamas)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	c, _ := json.Marshal(pegawaiApp)
	fmt.Println("next approval : ", string(c))

	if len(pegawaiApp) == 0 {
		return PengajuanSppdHistori{}, errors.New("Pegawai approval tidak ditemukan")
	}

	// insert pengajuan histori awal dari pengaju
	reqPayload.IdRuleApprovalDetail = rule2.ID.String()
	reqPayload.IdFungsionalitas = model.ParseString(pegawai.IdFungsionalitas)
	if reqFormat.Keterangan == nil {
		sk := "1"
		reqPayload.Keterangan = &sk
	}
	reqPayload.IdUnor = model.ParseString(pegawai.IdUnor)
	reqPayload.IdBidang = model.ParseString(pegawai.IdBidang)

	// Detail Employee
	pengajuanDetail := PengajuanSppdHistoriDetailRequest{
		IdPegawai: pegawai.ID.String(),
	}
	reqPayload.Detail = append(reqPayload.Detail, pengajuanDetail)

	data, _ = data.NewPengajuanSppdHistoriFormat(reqPayload, userID)
	err = s.PengajuanSppdHistoriRepository.Create(data)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	// insert next pengajuan ke user verifikator untuk approval
	reqPayload.Detail = make([]PengajuanSppdHistoriDetailRequest, 0)
	reqNext := reqPayload
	reqNext.IdPegawai = nil
	reqNext.Keterangan = nil
	reqNext.IdUnor = model.ParseString(rule2.IdUnor)
	reqNext.IdBidang = model.ParseString(rule2.IdBidang)
	reqNext.IdFungsionalitas = rule2.IdFungsionalitas
	reqNext.GroupApproval = &rule2.GroupApproval

	if rule2.GroupApproval == 1 && model.ParseInt(rule2.ApprovalLine) == 1 {
		reqNext.IdApprovalLine = pegawai.IdApprovalLine
	} else if rule2.GroupApproval == 1 && model.ParseInt(rule2.ApprovalLine) == 2 {
		reqNext.IdApprovalLine = pegawai.IdManager
	}

	for _, d := range pegawaiApp {
		pengajuanDetail := PengajuanSppdHistoriDetailRequest{
			IdPegawai: d.ID.String(),
		}
		reqNext.Detail = append(reqNext.Detail, pengajuanDetail)
	}

	nextPengajuan, _ := data.NewPengajuanSppdHistoriFormat(reqNext, userID)
	err = s.PengajuanSppdHistoriRepository.Create(nextPengajuan)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	// update status / bpd pegawai
	err = s.PengajuanSppdHistoriRepository.UpdateStatusBpdPegawai(StatusSPPD{
		ID:     *reqFormat.IdBpdPegawai,
		Status: "1",
	})
	if err != nil {
		fmt.Println(err)
	}

	return nextPengajuan, nil
}

// Function Create digunakan untuk ketika melakukan pengajuan awal
func (s *PengajuanSppdHistoriServiceImpl) Batal(reqFormat PengajuanSppdHistoriInputRequest, userID string) (err error) {
	reqPayload := PengajuanSppdHistoriRequest{
		IdSuratPerjalananDinas: reqFormat.IdSuratPerjalananDinas,
		Status:                 "4",
		TypeApproval:           "PENGAJUAN_SPPD",
	}

	// proses cek pegawai dari perjalanan dinas / bpd pegawai
	var idPegawaiStr string
	dinas, err := s.SuratPerjalananDinasRepository.ResolveByID(reqFormat.IdSuratPerjalananDinas)
	if err != nil {
		return errors.New("Data perjalanan dinas tidak ditemukan")
	}
	idPegawaiStr = model.ParseString(&dinas.IdPegawai)

	idPegawai, _ := uuid.FromString(idPegawaiStr)
	pegawai, err := s.PegawaiRepository.ResolveByIDDTO(idPegawai)
	if err != nil {
		return errors.New("Data Pegawai tidak ditemukan")
	}

	// insert pengajuan histori awal dari pengaju
	reqPayload.IdPegawai = &idPegawaiStr
	reqPayload.IdFungsionalitas = model.ParseString(pegawai.IdFungsionalitas)
	if reqFormat.Keterangan == nil {
		sk := "1"
		reqPayload.Keterangan = &sk
	}
	reqPayload.IdUnor = model.ParseString(pegawai.IdUnor)
	reqPayload.IdBidang = model.ParseString(pegawai.IdBidang)

	// Detail Employee
	pengajuanDetail := PengajuanSppdHistoriDetailRequest{
		IdPegawai: pegawai.ID.String(),
	}
	reqPayload.Detail = append(reqPayload.Detail, pengajuanDetail)

	var data PengajuanSppdHistori
	data, _ = data.NewPengajuanSppdHistoriFormat(reqPayload, userID)
	err = s.PengajuanSppdHistoriRepository.Create(data)
	if err != nil {
		return err
	}

	err = s.PengajuanSppdHistoriRepository.UpdateStatusSPPD(StatusSPPD{
		ID:     reqFormat.IdSuratPerjalananDinas,
		Status: "4",
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PengajuanSppdHistoriServiceImpl) RevisiPenyelesaianBiaya(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error) {
	data, err = s.PengajuanSppdHistoriRepository.ResolveByID(reqFormat.ID)
	if err != nil {
		return PengajuanSppdHistori{}, errors.New("Data pengajuan bpd histori tidak ditemukan")
	}
	reqFormat.TypeApproval = "PENYELESAIAN"
	reqFormat.IdSuratPerjalananDinas = data.IdSuratPerjalananDinas
	reqFormat.IdBpdPegawai = data.IdBpdPegawai

	if reqFormat.IdBpdPegawai == nil {
		return PengajuanSppdHistori{}, errors.New("ID BPD Pegawai tidak ditemukan")
	}

	ket := "2"
	reqPayload := PengajuanSppdHistoriRequest{
		IdSuratPerjalananDinas: reqFormat.IdSuratPerjalananDinas,
		IdRuleApprovalDetail:   data.IdRuleApprovalDetail,
		IdBpdPegawai:           reqFormat.IdBpdPegawai,
		IdFungsionalitas:       data.IdFungsionalitas,
		IdPegawai:              data.IdPegawai,
		IdUnor:                 data.IdUnor,
		Keterangan:             &ket,
		Status:                 "2",
		TypeApproval:           reqFormat.TypeApproval,
	}

	// Detail Employee
	pengajuanDetail := PengajuanSppdHistoriDetailRequest{
		IdPegawai: *reqFormat.IdPegawai,
	}
	reqPayload.Detail = append(reqPayload.Detail, pengajuanDetail)

	data, _ = data.NewPengajuanSppdHistoriFormat(reqPayload, userID)
	err = s.PengajuanSppdHistoriRepository.Create(data)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	// next pengajuan
	_, _, err = s.CreateNextApproval(reqFormat, userID)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	StatusSPPD := StatusSPPD{
		ID:     *reqFormat.IdBpdPegawai,
		Status: reqFormat.Status,
	}
	err = s.PengajuanSppdHistoriRepository.UpdateStatusBpdPegawai(StatusSPPD)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	return data, nil
}

func (s *PengajuanSppdHistoriServiceImpl) CheckGroupRule(req master.RuleParams) (rule master.RuleApproval, exist bool) {
	// check imployee ID
	employee, err := s.RuleApprovalRepository.ResolveByKode(master.RuleParams{
		Jenis:     req.Jenis,
		IdPegawai: req.IdPegawai,
	})
	fmt.Println("Group employee ID", err)
	if err == nil {
		return employee, true
	}

	// check fungsionalitas
	fungsionalitas, err := s.RuleApprovalRepository.ResolveByKode(master.RuleParams{
		Jenis:            req.Jenis,
		IdFungsionalitas: req.IdFungsionalitas,
	})
	fmt.Println("Group fungsionalitas", err)
	if err == nil {
		return fungsionalitas, true
	}

	// check all employee
	allEmployee, err := s.RuleApprovalRepository.ResolveByKode(master.RuleParams{
		Jenis:     req.Jenis,
		GroupRule: 1,
	})
	fmt.Println("Group all employee", err)
	if err == nil {
		return allEmployee, true
	}

	return master.RuleApproval{}, false
}

func (s *PengajuanSppdHistoriServiceImpl) GetPegawaiGroupApproval(req master.RuleDetailParams) (pegawai []master.PegawaiDTO, err error) {
	groupApproval := req.GroupApproval
	payload := master.PegawaiParams{}

	fmt.Println("groupApproval", groupApproval)
	if groupApproval == 1 { // approval line
		if req.IdApprovalLine == "" && req.ApprovalLine == 1 {
			return []master.PegawaiDTO{}, errors.New("Pegawai approval line tidak ditemukan")
		}

		if req.IdManager == "" && req.ApprovalLine == 2 {
			return []master.PegawaiDTO{}, errors.New("Pegawai approval line 2 tidak ditemukan")
		}

		if req.ApprovalLine == 1 {
			payload.IdPegawai = req.IdApprovalLine
		} else if req.ApprovalLine == 2 {
			payload.IdPegawai = req.IdManager
		}
	} else if groupApproval == 2 { // employee ID
		payload.IdPegawai = req.IdPegawai
	} else if groupApproval == 3 { // Fungsionalitas
		payload.IdFungsionalitas = req.IdFungsionalitas
		payload.IdBidang = req.IdBidang
		payload.IdUnor = req.IdUnor
	}

	fmt.Println("payload", payload)

	pegawai, err = s.PegawaiRepository.GetAllPegawai(payload)
	return
}

// Function RevisiPengajuan digunakan untuk ketika melakukan pengajuan awal
func (s *PengajuanSppdHistoriServiceImpl) RevisiPengajuan(reqFormat PengajuanSppdHistoriInputRequest, userID string) (data PengajuanSppdHistori, err error) {
	reqPayload := PengajuanSppdHistoriRequest{
		IdSuratPerjalananDinas: reqFormat.IdSuratPerjalananDinas,
		IdRuleApprovalDetail:   reqFormat.IdRuleApprovalDetail,
		IdBpdPegawai:           reqFormat.IdBpdPegawai,
		Catatan:                reqFormat.Catatan,
		Keterangan:             reqFormat.Keterangan,
		Status:                 reqFormat.Status,
		Jenis:                  reqFormat.Jenis,
		TypeApproval:           reqFormat.TypeApproval,
		TenantId:               reqFormat.TenantId,
		IdBranch:               reqFormat.IdBranch,
	}

	// proses cek pegawai dari perjalanan dinas / bpd pegawai
	var idPegawaiStr string
	if reqFormat.TypeApproval == "PENYELESAIAN" {
		if reqFormat.IdBpdPegawai == nil {
			return PengajuanSppdHistori{}, errors.New("ID BPD Pegawai tidak ditemukan")
		}
		dinas, err := s.PerjalananDinasRepository.ResolveBpdPegawaiByID(*reqFormat.IdBpdPegawai)
		if err != nil {
			return PengajuanSppdHistori{}, errors.New("Data BPD Pegawai tidak ditemukan")
		}
		idPegawaiStr = dinas.IdPegawai
	} else {
		dinas, err := s.SuratPerjalananDinasRepository.ResolveDtoByID(reqFormat.IdSuratPerjalananDinas)
		if err != nil {
			return PengajuanSppdHistori{}, errors.New("Data perjalanan dinas tidak ditemukan")
		}
		idPegawaiStr = model.ParseString(&dinas.IdPegawai)
	}

	idPegawai, _ := uuid.FromString(idPegawaiStr)
	fmt.Println("pegwaiiid", idPegawai)
	pegawai, err := s.PegawaiRepository.ResolveByIDDTO(idPegawai)
	if err != nil {
		return PengajuanSppdHistori{}, errors.New("Data Pegawai tidak ditemukan")
	}

	reqPayload.IdPegawai = &idPegawaiStr
	// cek group rule
	ruleParam := master.RuleParams{
		Jenis:            reqFormat.Jenis,
		TypeApproval:     reqFormat.TypeApproval,
		IdPegawai:        model.ParseString(reqFormat.IdPegawai),
		IdBidang:         model.ParseString(pegawai.IdBidang),
		IdUnor:           model.ParseString(pegawai.IdUnor),
		IdFungsionalitas: model.ParseString(pegawai.IdFungsionalitas),
	}

	checkRule, existRule := s.CheckGroupRule(ruleParam)
	if !existRule {
		return PengajuanSppdHistori{}, errors.New("Rule Approval belum di setting")
	}

	ruleParam.GroupRule = checkRule.GroupRule
	fmt.Println("GroupRule", checkRule.GroupRule)
	rule, _ := s.RuleApprovalRepository.GetAllRuleApprovalDetailByKode(ruleParam)

	// find pegawai
	rule2 := rule[0]
	idUnor := ""
	idBidang := ""
	if rule2.IdUnor == nil && !model.ParseBoolean(rule2.IsHead) {
		idUnor = model.ParseString(pegawai.IdUnor)
	} else {
		idUnor = *rule2.IdUnor
	}

	if rule2.IdBidang != nil {
		idBidang = *rule2.IdBidang
	} else {
		idBidang = model.ParseString(pegawai.IdBidang)
	}

	RDetParamas := master.RuleDetailParams{
		IdPegawai:        model.ParseString(rule2.IdPegawai),
		IdApprovalLine:   model.ParseString(pegawai.IdApprovalLine),
		IdManager:        model.ParseString(pegawai.IdManager),
		ApprovalLine:     model.ParseInt(rule2.ApprovalLine),
		IdFungsionalitas: rule2.IdFungsionalitas,
		IdUnor:           idUnor,
		IdBidang:         idBidang,
		GroupApproval:    rule2.GroupApproval,
	}

	// fmt.Println("IdPegawai:", model.ParseString(rule2.IdPegawai))
	// fmt.Println("IdApprovalLine:", model.ParseString(pegawai.IdApprovalLine))
	// fmt.Println("IdManager:", model.ParseString(pegawai.IdManager))
	// fmt.Println("ApprovalLine:", model.ParseInt(rule2.ApprovalLine))
	// fmt.Println("IdFungsionalitas:", rule2.IdFungsionalitas)
	// fmt.Println("IdUnor:", idUnor)
	// fmt.Println("IdBidang:", idBidang)
	// fmt.Println("GroupApproval:", rule2.GroupApproval)

	pegawaiApp, err := s.GetPegawaiGroupApproval(RDetParamas)

	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	c, _ := json.Marshal(pegawaiApp)
	fmt.Println("next approval : ", string(c))

	if len(pegawaiApp) == 0 {
		return PengajuanSppdHistori{}, errors.New("Pegawai approval tidak ditemukan")
	}

	// insert pengajuan histori awal dari pengaju
	reqPayload.IdRuleApprovalDetail = rule2.ID.String()
	reqPayload.IdFungsionalitas = model.ParseString(pegawai.IdFungsionalitas)
	if reqFormat.Keterangan == nil {
		sk := "1"
		reqPayload.Keterangan = &sk
	}
	reqPayload.IdUnor = model.ParseString(pegawai.IdUnor)
	reqPayload.IdBidang = model.ParseString(pegawai.IdBidang)

	// Detail Employee
	pengajuanDetail := PengajuanSppdHistoriDetailRequest{
		IdPegawai: pegawai.ID.String(),
	}
	reqPayload.Detail = append(reqPayload.Detail, pengajuanDetail)

	data, _ = data.NewPengajuanSppdHistoriFormat(reqPayload, userID)
	err = s.PengajuanSppdHistoriRepository.Create(data)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	// insert next pengajuan ke user verifikator untuk approval
	reqPayload.Detail = make([]PengajuanSppdHistoriDetailRequest, 0)
	reqNext := reqPayload
	reqNext.IdPegawai = nil
	reqNext.Keterangan = nil
	reqNext.IdUnor = model.ParseString(rule2.IdUnor)
	reqNext.IdBidang = model.ParseString(rule2.IdBidang)
	reqNext.IdFungsionalitas = rule2.IdFungsionalitas
	reqNext.GroupApproval = &rule2.GroupApproval

	if rule2.GroupApproval == 1 && model.ParseInt(rule2.ApprovalLine) == 1 {
		reqNext.IdApprovalLine = pegawai.IdApprovalLine
	} else if rule2.GroupApproval == 1 && model.ParseInt(rule2.ApprovalLine) == 2 {
		reqNext.IdApprovalLine = pegawai.IdManager
	}

	for _, d := range pegawaiApp {
		pengajuanDetail := PengajuanSppdHistoriDetailRequest{
			IdPegawai: d.ID.String(),
		}
		reqNext.Detail = append(reqNext.Detail, pengajuanDetail)
	}

	nextPengajuan, _ := data.NewPengajuanSppdHistoriFormat(reqNext, userID)
	err = s.PengajuanSppdHistoriRepository.Create(nextPengajuan)
	if err != nil {
		return PengajuanSppdHistori{}, err
	}

	// update status perjalanan dinas / bpd pegawai
	if reqFormat.TypeApproval == "PENYELESAIAN" && reqFormat.IdBpdPegawai != nil {
		err = s.PengajuanSppdHistoriRepository.UpdateStatusBpdPegawai(StatusSPPD{
			ID:     *reqFormat.IdBpdPegawai,
			Status: "1",
		})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err = s.PengajuanSppdHistoriRepository.UpdateStatusSPPD(StatusSPPD{
			ID:     reqFormat.IdSuratPerjalananDinas,
			Status: "1",
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	return nextPengajuan, nil
}
