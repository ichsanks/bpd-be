package master

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type RuleApprovalService interface {
	Create(reqFormat RuleApprovalRequest, userID string, tenantId uuid.UUID) (data RuleApproval, err error)
	Update(reqFormat RuleApprovalRequest, userID string, tenantId uuid.UUID) (data RuleApproval, err error)
	ResolveByIDDTO(id string) (data RuleApprovalDTO, err error)
	ResolveAll(req model.StandardRequestRuleApproval) (data pagination.Response, err error)
	DeleteByID(id string, userID string) error
	GetAll(idFungsionalitas string) (data []RuleApprovalDTO, err error)
	ResolveTtd(id string) (data RuleApprovalTtd, err error)
}

type RuleApprovalServiceImpl struct {
	RuleApprovalRepository RuleApprovalRepository
	Config                 *configs.Config
}

func ProvideRuleApprovalServiceImpl(repository RuleApprovalRepository, config *configs.Config) *RuleApprovalServiceImpl {
	s := new(RuleApprovalServiceImpl)
	s.RuleApprovalRepository = repository
	s.Config = config
	return s
}

func (s *RuleApprovalServiceImpl) ResolveAll(req model.StandardRequestRuleApproval) (data pagination.Response, err error) {
	return s.RuleApprovalRepository.ResolveAll(req)
}

func (s *RuleApprovalServiceImpl) GetAll(idFungsionalitas string) (data []RuleApprovalDTO, err error) {
	return s.RuleApprovalRepository.GetAll(idFungsionalitas)
}

func (s *RuleApprovalServiceImpl) Create(reqFormat RuleApprovalRequest, userID string, tenantId uuid.UUID) (data RuleApproval, err error) {
	data, _ = data.NewRuleApprovalFormat(reqFormat, userID, tenantId)
	err = s.RuleApprovalRepository.Create(data)
	if err != nil {
		return RuleApproval{}, err
	}
	return data, nil
}

func (s *RuleApprovalServiceImpl) Update(reqFormat RuleApprovalRequest, userID string, tenantId uuid.UUID) (data RuleApproval, err error) {
	data, _ = data.NewRuleApprovalFormat(reqFormat, userID, tenantId)
	err = s.RuleApprovalRepository.UpdateRuleApproval(data)
	if err != nil {
		return RuleApproval{}, err
	}
	return data, nil
}

func (s *RuleApprovalServiceImpl) ResolveByIDDTO(id string) (data RuleApprovalDTO, err error) {
	return s.RuleApprovalRepository.ResolveByIDDTO(id)
}

func (s *RuleApprovalServiceImpl) DeleteByID(id string, userID string) error {
	rule, err := s.RuleApprovalRepository.ResolveByID(id)

	if err != nil {
		return errors.New("Data Rule Approval dengan ID :" + id + " tidak ditemukan")
	}

	now := time.Now()
	rule.IsDeleted = true
	rule.UpdatedBy = &userID
	rule.UpdatedAt = &now
	err = s.RuleApprovalRepository.Update(rule)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data rule approval dengan ID: " + id)
	}

	return nil
}

func (s *RuleApprovalServiceImpl) ResolveTtd(id string) (data RuleApprovalTtd, err error) {
	return s.RuleApprovalRepository.ResolveTtd(id)
}
