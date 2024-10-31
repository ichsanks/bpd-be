package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type SettingBiayaService interface {
	Create(req SettingBiayaFormat, userId uuid.UUID) (newSettingBiaya SettingBiaya, err error)
	GetAll(req model.StandardRequest) (data []SettingBiaya, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	Update(req SettingBiayaUpdateFormat, userId uuid.UUID) (newSettingBiaya SettingBiaya, err error)
	ResolveByID(id uuid.UUID) (data SettingBiaya, err error)
	DeleteSoft(id uuid.UUID, userId uuid.UUID) error
}

type SettingBiayaServiceImpl struct {
	SettingBiayaRepository SettingBiayaRepository
	Config                 *configs.Config
}

func ProvideSettingBiayaServiceImpl(repository SettingBiayaRepository) *SettingBiayaServiceImpl {
	s := new(SettingBiayaServiceImpl)
	s.SettingBiayaRepository = repository
	return s
}

func (s *SettingBiayaServiceImpl) Create(req SettingBiayaFormat, userId uuid.UUID) (newSettingBiaya SettingBiaya, err error) {

	exist, err := s.SettingBiayaRepository.ExistByNama("", *req.IdBranch, *req.IdBodLevel, *req.IdJenisTujuan, req.IdJenisBiaya)
	if exist {
		x := errors.New("Nama SettingBiaya sudah dipakai")
		return SettingBiaya{}, x
	}
	if err != nil {
		return SettingBiaya{}, err
	}

	newSettingBiaya, err = newSettingBiaya.SettingBiayaNewFormat(req, userId)

	err = s.SettingBiayaRepository.Create(newSettingBiaya)

	if err != nil {
		return SettingBiaya{}, err
	}

	return newSettingBiaya, nil
}

func (s *SettingBiayaServiceImpl) GetAll(req model.StandardRequest) (data []SettingBiaya, err error) {
	return s.SettingBiayaRepository.GetAll(req)
}

func (s *SettingBiayaServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.SettingBiayaRepository.ResolveAll(request)
}

func (s *SettingBiayaServiceImpl) DeleteByID(id uuid.UUID) error {
	_, err := s.SettingBiayaRepository.ResolveByID(id)
	if err != nil {
		return errors.New("Data SettingBiaya dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.SettingBiayaRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data SettingBiaya dengan ID: " + id.String())
	}
	return nil
}

func (s *SettingBiayaServiceImpl) Update(req SettingBiayaUpdateFormat, userId uuid.UUID) (newSettingBiaya SettingBiaya, err error) {
	exist, err := s.SettingBiayaRepository.ExistByNama(req.ID.String(), *req.IdBranch, *req.IdBodLevel, *req.IdJenisTujuan, req.IdJenisBiaya)
	if exist {
		x := errors.New("Nama SettingBiaya sudah dipakai")
		return SettingBiaya{}, x
	}
	if err != nil {
		return SettingBiaya{}, err
	}
	newSettingBiaya, _ = newSettingBiaya.SettingBiayaUpdateFormat(req, userId)
	err = s.SettingBiayaRepository.Update(newSettingBiaya)
	if err != nil {
		return SettingBiaya{}, err
	}
	return newSettingBiaya, nil
}

func (s *SettingBiayaServiceImpl) ResolveByID(id uuid.UUID) (newSettingBiaya SettingBiaya, err error) {
	return s.SettingBiayaRepository.ResolveByID(id)
}

func (s *SettingBiayaServiceImpl) DeleteSoft(id uuid.UUID, userId uuid.UUID) error {
	newSettingBiaya, err := s.SettingBiayaRepository.ResolveByID(id)

	if err != nil {
		return errors.New("Data SettingBiaya dengan ID :" + id.String() + " tidak ditemukan")
	}

	// exist := s.SettingBiayaRepository.ExistRelasiStatus(id)
	// if exist {
	// 	return errors.New(" Data Jenis Biaya dengan Nama :" + *&newSettingBiaya.Nama + " sedang digunakan")
	// }

	newSettingBiaya.SoftDelete(userId)
	err = s.SettingBiayaRepository.Update(newSettingBiaya)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data SettingBiaya dengan ID: " + id.String())
	}
	return nil
}
