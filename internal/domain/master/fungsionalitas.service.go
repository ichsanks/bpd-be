package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type FungsionalitasService interface {
	Create(req RequestFungsionalitas, userId uuid.UUID, tenantId uuid.UUID) (data Fungsionalitas, err error)
	Update(req RequestFungsionalitas, userId uuid.UUID, tenantId uuid.UUID) (data Fungsionalitas, err error)
	GetAll() (data []Fungsionalitas, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Fungsionalitas, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type FungsionalitasServiceImpl struct {
	FungsionalitasRepository FungsionalitasRepository
	Config                   *configs.Config
}

func ProvideFungsionalitasServiceImpl(repository FungsionalitasRepository, config *configs.Config) *FungsionalitasServiceImpl {
	s := new(FungsionalitasServiceImpl)
	s.FungsionalitasRepository = repository
	s.Config = config
	return s
}

func (s *FungsionalitasServiceImpl) Create(req RequestFungsionalitas, userId uuid.UUID, tenantId uuid.UUID) (data Fungsionalitas, err error) {
	existNama, err := s.FungsionalitasRepository.ExistByNama(req.Nama, "")
	if existNama {
		return Fungsionalitas{}, errors.New("Nama Fungsionalitas sudah dipakai")
	}

	data, _ = data.FungsionalitasFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.FungsionalitasRepository.Create(data)
	if err != nil {
		return Fungsionalitas{}, err
	}
	return data, nil
}

func (s *FungsionalitasServiceImpl) GetAll() (data []Fungsionalitas, err error) {
	return s.FungsionalitasRepository.GetAll()
}

func (s *FungsionalitasServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.FungsionalitasRepository.ResolveAll(request)
}

func (s *FungsionalitasServiceImpl) DeleteByID(id uuid.UUID) error {
	newLayanan, err := s.FungsionalitasRepository.ResolveByID(id)

	if err != nil || (Fungsionalitas{}) == newLayanan {
		return errors.New("Data Fungsionalitas dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.FungsionalitasRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Fungsionalitas dengan ID: " + id.String())
	}
	return nil
}

func (s *FungsionalitasServiceImpl) Update(req RequestFungsionalitas, userId uuid.UUID, tenantId uuid.UUID) (data Fungsionalitas, err error) {
	existNama, err := s.FungsionalitasRepository.ExistByNama(req.Nama, req.ID.String())
	if existNama {
		return Fungsionalitas{}, errors.New("Nama Fungsionalitas sudah dipakai")
	}

	if err != nil {
		return Fungsionalitas{}, err
	}

	data, _ = data.FungsionalitasFormatRequest(req, userId, tenantId)
	err = s.FungsionalitasRepository.Update(data)
	if err != nil {
		return Fungsionalitas{}, err
	}
	return data, nil
}

func (s *FungsionalitasServiceImpl) ResolveByID(id uuid.UUID) (data Fungsionalitas, err error) {
	return s.FungsionalitasRepository.ResolveByID(id)
}

func (s *FungsionalitasServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newFungsionalitas, err := s.FungsionalitasRepository.ResolveByID(id)

	if err != nil || (Fungsionalitas{}) == newFungsionalitas {
		return errors.New("Data Fungsionalitas dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.FungsionalitasRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data Fungsionalitas dengan Nama :" + *&newFungsionalitas.Nama + " sedang digunakan")
	}

	newFungsionalitas.SoftDelete(userId)
	err = s.FungsionalitasRepository.Update(newFungsionalitas)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Fungsionalitas dengan ID: " + id.String())
	}
	return nil
}
