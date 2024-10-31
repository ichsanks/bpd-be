package master

import (
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
)

type JenisApprovalService interface {
	GetAll(ids string) (data []JenisApproval, err error)
}

type JenisApprovalServiceImpl struct {
	JenisApprovalRepository JenisApprovalRepository
	Config                  *configs.Config
}

func ProvideJenisApprovalServiceImpl(repository JenisApprovalRepository) *JenisApprovalServiceImpl {
	s := new(JenisApprovalServiceImpl)
	s.JenisApprovalRepository = repository
	return s
}

func (s *JenisApprovalServiceImpl) GetAll(ids string) (data []JenisApproval, err error) {
	return s.JenisApprovalRepository.GetAll(ids)
}
