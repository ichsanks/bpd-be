package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type JobGradeService interface {
	Create(req RequestJobGrade, userId uuid.UUID, tenantId uuid.UUID) (data JobGrade, err error)
	Update(req RequestJobGrade, userId uuid.UUID, tenantId uuid.UUID) (data JobGrade, err error)
	GetAll(req model.StandardRequest) (data []JobGrade, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data JobGrade, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type JobGradeServiceImpl struct {
	JobGradeRepository JobGradeRepository
	Config             *configs.Config
}

func ProvideJobGradeServiceImpl(repository JobGradeRepository, config *configs.Config) *JobGradeServiceImpl {
	s := new(JobGradeServiceImpl)
	s.JobGradeRepository = repository
	s.Config = config
	return s
}

func (s *JobGradeServiceImpl) Create(req RequestJobGrade, userId uuid.UUID, tenantId uuid.UUID) (data JobGrade, err error) {
	existKode, err := s.JobGradeRepository.ExistByKode(req.Kode, "", *req.IdBranch)
	if existKode {
		return JobGrade{}, errors.New("Kode JobGrade sudah dipakai")
	}
	if err != nil {
		return
	}

	existNama, err := s.JobGradeRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return JobGrade{}, errors.New("Nama JobGrade sudah dipakai")
	}
	if err != nil {
		return
	}

	data, _ = data.JobGradeFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.JobGradeRepository.Create(data)
	if err != nil {
		return JobGrade{}, err
	}
	return data, nil
}

func (s *JobGradeServiceImpl) GetAll(req model.StandardRequest) (data []JobGrade, err error) {
	return s.JobGradeRepository.GetAll(req)
}

func (s *JobGradeServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.JobGradeRepository.ResolveAll(request)
}

func (s *JobGradeServiceImpl) DeleteByID(id uuid.UUID) error {
	newJobGrade, err := s.JobGradeRepository.ResolveByID(id)

	if err != nil || (JobGrade{}) == newJobGrade {
		return errors.New("Data JobGrade dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.JobGradeRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JobGrade dengan ID: " + id.String())
	}
	return nil
}

func (s *JobGradeServiceImpl) Update(req RequestJobGrade, userId uuid.UUID, tenantId uuid.UUID) (data JobGrade, err error) {
	existKode, err := s.JobGradeRepository.ExistByKode(req.Kode, req.ID.String(), *req.IdBranch)
	if existKode {
		return JobGrade{}, errors.New("Kode JobGrade sudah dipakai")
	}

	existNama, err := s.JobGradeRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return JobGrade{}, errors.New("Nama JobGrade sudah dipakai")
	}

	data, _ = data.JobGradeFormatRequest(req, userId, tenantId)
	err = s.JobGradeRepository.Update(data)
	if err != nil {
		return JobGrade{}, err
	}
	return data, nil
}

func (s *JobGradeServiceImpl) ResolveByID(id uuid.UUID) (data JobGrade, err error) {
	return s.JobGradeRepository.ResolveByID(id)
}

func (s *JobGradeServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newJobGrade, err := s.JobGradeRepository.ResolveByID(id)

	if err != nil || (JobGrade{}) == newJobGrade {
		return errors.New("Data JobGrade dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.JobGradeRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data JobGrade dengan Nama :" + *&newJobGrade.Nama + " sedang digunakan")
	}

	newJobGrade.SoftDelete(userId)
	err = s.JobGradeRepository.Update(newJobGrade)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JobGrade dengan ID: " + id.String())
	}
	return nil
}
