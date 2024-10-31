package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type StatusKontrakService interface {
	Create(req RequestStatusKontrak, userId uuid.UUID, tenantId uuid.UUID) (data StatusKontrak, err error)
	Update(req RequestStatusKontrak, userId uuid.UUID, tenantId uuid.UUID) (data StatusKontrak, err error)
	GetAll(req model.StandardRequest) (data []StatusKontrak, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data StatusKontrak, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type StatusKontrakServiceImpl struct {
	StatusKontrakRepository StatusKontrakRepository
	Config                  *configs.Config
}

func ProvideStatusKontrakServiceImpl(repository StatusKontrakRepository, config *configs.Config) *StatusKontrakServiceImpl {
	s := new(StatusKontrakServiceImpl)
	s.StatusKontrakRepository = repository
	s.Config = config
	return s
}

func (s *StatusKontrakServiceImpl) Create(req RequestStatusKontrak, userId uuid.UUID, tenantId uuid.UUID) (data StatusKontrak, err error) {

	existNama, err := s.StatusKontrakRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return StatusKontrak{}, errors.New("Nama StatusKontrak sudah dipakai")
	}

	data, _ = data.StatusKontrakFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.StatusKontrakRepository.Create(data)
	if err != nil {
		return StatusKontrak{}, err
	}
	return data, nil
}

func (s *StatusKontrakServiceImpl) GetAll(req model.StandardRequest) (data []StatusKontrak, err error) {
	return s.StatusKontrakRepository.GetAll(req)
}

func (s *StatusKontrakServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.StatusKontrakRepository.ResolveAll(request)
}

func (s *StatusKontrakServiceImpl) DeleteByID(id uuid.UUID) error {
	newStatusKontrak, err := s.StatusKontrakRepository.ResolveByID(id)

	if err != nil || (StatusKontrak{}) == newStatusKontrak {
		return errors.New("Data StatusKontrak dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.StatusKontrakRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data StatusKontrak dengan ID: " + id.String())
	}
	return nil
}

func (s *StatusKontrakServiceImpl) Update(req RequestStatusKontrak, userId uuid.UUID, tenantId uuid.UUID) (data StatusKontrak, err error) {

	existNama, err := s.StatusKontrakRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return StatusKontrak{}, errors.New("Nama StatusKontrak sudah dipakai")
	}

	data, _ = data.StatusKontrakFormatRequest(req, userId, tenantId)
	err = s.StatusKontrakRepository.Update(data)
	if err != nil {
		return StatusKontrak{}, err
	}
	return data, nil
}

func (s *StatusKontrakServiceImpl) ResolveByID(id uuid.UUID) (data StatusKontrak, err error) {
	return s.StatusKontrakRepository.ResolveByID(id)
}

func (s *StatusKontrakServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newStatusKontrak, err := s.StatusKontrakRepository.ResolveByID(id)

	if err != nil || (StatusKontrak{}) == newStatusKontrak {
		return errors.New("Data StatusKontrak dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.StatusKontrakRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data StatusKontrak dengan Nama :" + *&newStatusKontrak.Nama + " sedang digunakan")
	}

	newStatusKontrak.SoftDelete(userId)
	err = s.StatusKontrakRepository.Update(newStatusKontrak)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data StatusKontrak dengan ID: " + id.String())
	}
	return nil
}
