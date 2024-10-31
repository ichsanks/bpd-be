package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type LevelBodService interface {
	Create(req RequestLevelBod, userId uuid.UUID, tenantId uuid.UUID) (data LevelBod, err error)
	Update(req RequestLevelBod, userId uuid.UUID, tenantId uuid.UUID) (data LevelBod, err error)
	GetAll(req model.StandardRequest) (data []LevelBod, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data LevelBod, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type LevelBodServiceImpl struct {
	LevelBodRepository LevelBodRepository
	Config             *configs.Config
}

func ProvideLevelBodServiceImpl(repository LevelBodRepository, config *configs.Config) *LevelBodServiceImpl {
	s := new(LevelBodServiceImpl)
	s.LevelBodRepository = repository
	s.Config = config
	return s
}

func (s *LevelBodServiceImpl) Create(req RequestLevelBod, userId uuid.UUID, tenantId uuid.UUID) (data LevelBod, err error) {
	existKode, err := s.LevelBodRepository.ExistByKode(req.Kode, "", *req.IdBranch)
	if existKode {
		return LevelBod{}, errors.New("Kode LevelBod sudah dipakai")
	}

	existNama, err := s.LevelBodRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return LevelBod{}, errors.New("Nama LevelBod sudah dipakai")
	}

	data, _ = data.LevelBodFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.LevelBodRepository.Create(data)
	if err != nil {
		return LevelBod{}, err
	}
	return data, nil
}

func (s *LevelBodServiceImpl) GetAll(req model.StandardRequest) (data []LevelBod, err error) {
	return s.LevelBodRepository.GetAll(req)
}

func (s *LevelBodServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.LevelBodRepository.ResolveAll(request)
}

func (s *LevelBodServiceImpl) DeleteByID(id uuid.UUID) error {
	newLevelBod, err := s.LevelBodRepository.ResolveByID(id)

	if err != nil || (LevelBod{}) == newLevelBod {
		return errors.New("Data LevelBod dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.LevelBodRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data LevelBod dengan ID: " + id.String())
	}
	return nil
}

func (s *LevelBodServiceImpl) Update(req RequestLevelBod, userId uuid.UUID, tenantId uuid.UUID) (data LevelBod, err error) {
	existKode, err := s.LevelBodRepository.ExistByKode(req.Kode, req.ID.String(), *req.IdBranch)
	if existKode {
		return LevelBod{}, errors.New("Kode LevelBod sudah dipakai")
	}

	existNama, err := s.LevelBodRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return LevelBod{}, errors.New("Nama LevelBod sudah dipakai")
	}

	data, _ = data.LevelBodFormatRequest(req, userId, tenantId)
	err = s.LevelBodRepository.Update(data)
	if err != nil {
		return LevelBod{}, err
	}
	return data, nil
}

func (s *LevelBodServiceImpl) ResolveByID(id uuid.UUID) (data LevelBod, err error) {
	return s.LevelBodRepository.ResolveByID(id)
}

func (s *LevelBodServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newLevelBod, err := s.LevelBodRepository.ResolveByID(id)

	if err != nil || (LevelBod{}) == newLevelBod {
		return errors.New("Data LevelBod dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.LevelBodRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data LevelBod dengan Nama :" + *&newLevelBod.Nama + " sedang digunakan")
	}

	newLevelBod.SoftDelete(userId)
	err = s.LevelBodRepository.Update(newLevelBod)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data LevelBod dengan ID: " + id.String())
	}
	return nil
}
