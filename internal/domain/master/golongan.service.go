package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type GolonganService interface {
	Create(req RequestGolongan, userId uuid.UUID, tenantId uuid.UUID) (data Golongan, err error)
	Update(req RequestGolongan, userId uuid.UUID, tenantId uuid.UUID) (data Golongan, err error)
	GetAll(req model.StandardRequest) (data []Golongan, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Golongan, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type GolonganServiceImpl struct {
	GolonganRepository GolonganRepository
	Config             *configs.Config
}

func ProvideGolonganServiceImpl(repository GolonganRepository, config *configs.Config) *GolonganServiceImpl {
	s := new(GolonganServiceImpl)
	s.GolonganRepository = repository
	s.Config = config
	return s
}

func (s *GolonganServiceImpl) Create(req RequestGolongan, userId uuid.UUID, tenantId uuid.UUID) (data Golongan, err error) {
	existKode, err := s.GolonganRepository.ExistByKode(req.Kode, "", *req.IdBranch)
	if existKode {
		return Golongan{}, errors.New("Kode Golongan sudah dipakai")
	}

	existNama, err := s.GolonganRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return Golongan{}, errors.New("Nama Golongan sudah dipakai")
	}

	data, _ = data.GolonganFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.GolonganRepository.Create(data)
	if err != nil {
		return Golongan{}, err
	}
	return data, nil
}

func (s *GolonganServiceImpl) GetAll(req model.StandardRequest) (data []Golongan, err error) {
	return s.GolonganRepository.GetAll(req)
}

func (s *GolonganServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.GolonganRepository.ResolveAll(request)
}

func (s *GolonganServiceImpl) DeleteByID(id uuid.UUID) error {
	newGolongan, err := s.GolonganRepository.ResolveByID(id)

	if err != nil || (Golongan{}) == newGolongan {
		return errors.New("Data Golongan dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.GolonganRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Golongan dengan ID: " + id.String())
	}
	return nil
}

func (s *GolonganServiceImpl) Update(req RequestGolongan, userId uuid.UUID, tenantId uuid.UUID) (data Golongan, err error) {
	existKode, err := s.GolonganRepository.ExistByKode(req.Kode, req.ID.String(), *req.IdBranch)
	if existKode {
		return Golongan{}, errors.New("Kode Golongan sudah dipakai")
	}

	existNama, err := s.GolonganRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return Golongan{}, errors.New("Nama Golongan sudah dipakai")
	}

	data, _ = data.GolonganFormatRequest(req, userId, tenantId)
	err = s.GolonganRepository.Update(data)
	if err != nil {
		return Golongan{}, err
	}
	return data, nil
}

func (s *GolonganServiceImpl) ResolveByID(id uuid.UUID) (data Golongan, err error) {
	return s.GolonganRepository.ResolveByID(id)
}

func (s *GolonganServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newGolongan, err := s.GolonganRepository.ResolveByID(id)

	if err != nil || (Golongan{}) == newGolongan {
		return errors.New("Data Golongan dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.GolonganRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data Golongan dengan Nama :" + *&newGolongan.Nama + " sedang digunakan")
	}

	newGolongan.SoftDelete(userId)
	err = s.GolonganRepository.Update(newGolongan)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Golongan dengan ID: " + id.String())
	}
	return nil
}
