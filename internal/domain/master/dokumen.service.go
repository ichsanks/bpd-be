package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type DokumenService interface {
	Create(req RequestDokumen, userId uuid.UUID, tenantId uuid.UUID) (data Dokumen, err error)
	Update(req RequestDokumen, userId uuid.UUID, tenantId uuid.UUID) (data Dokumen, err error)
	GetAll(req model.StandardRequest) (data []Dokumen, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Dokumen, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type DokumenServiceImpl struct {
	DokumenRepository DokumenRepository
	Config            *configs.Config
}

func ProvideDokumenServiceImpl(repository DokumenRepository, config *configs.Config) *DokumenServiceImpl {
	s := new(DokumenServiceImpl)
	s.DokumenRepository = repository
	s.Config = config
	return s
}

func (s *DokumenServiceImpl) Create(req RequestDokumen, userId uuid.UUID, tenantId uuid.UUID) (data Dokumen, err error) {

	existNama, err := s.DokumenRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return Dokumen{}, errors.New("Nama Dokumen sudah dipakai")
	}

	data, _ = data.DokumenFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.DokumenRepository.Create(data)
	if err != nil {
		return Dokumen{}, err
	}
	return data, nil
}

func (s *DokumenServiceImpl) GetAll(req model.StandardRequest) (data []Dokumen, err error) {
	return s.DokumenRepository.GetAll(req)
}

func (s *DokumenServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.DokumenRepository.ResolveAll(request)
}

func (s *DokumenServiceImpl) DeleteByID(id uuid.UUID) error {
	newDokumen, err := s.DokumenRepository.ResolveByID(id)

	if err != nil || (Dokumen{}) == newDokumen {
		return errors.New("Data Dokumen dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.DokumenRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Dokumen dengan ID: " + id.String())
	}
	return nil
}

func (s *DokumenServiceImpl) Update(req RequestDokumen, userId uuid.UUID, tenantId uuid.UUID) (data Dokumen, err error) {

	existNama, err := s.DokumenRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return Dokumen{}, errors.New("Nama Dokumen sudah dipakai")
	}

	data, _ = data.DokumenFormatRequest(req, userId, tenantId)
	err = s.DokumenRepository.Update(data)
	if err != nil {
		return Dokumen{}, err
	}
	return data, nil
}

func (s *DokumenServiceImpl) ResolveByID(id uuid.UUID) (data Dokumen, err error) {
	return s.DokumenRepository.ResolveByID(id)
}

func (s *DokumenServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newDokumen, err := s.DokumenRepository.ResolveByID(id)

	if err != nil || (Dokumen{}) == newDokumen {
		return errors.New("Data Dokumen dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.DokumenRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data Dokumen dengan Nama :" + *&newDokumen.Nama + " sedang digunakan")
	}

	newDokumen.SoftDelete(userId)
	err = s.DokumenRepository.Update(newDokumen)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Dokumen dengan ID: " + id.String())
	}
	return nil
}
