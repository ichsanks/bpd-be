package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type SyaratDokumenService interface {
	Create(req RequestSyaratDokumen, userId uuid.UUID, tenantId uuid.UUID) (data SyaratDokumen, err error)
	Update(req RequestSyaratDokumen, userId uuid.UUID, tenantId uuid.UUID) (data SyaratDokumen, err error)
	GetAll(req model.StandardRequest) (data []SyaratDokumen, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data SyaratDokumen, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type SyaratDokumenServiceImpl struct {
	SyaratDokumenRepository SyaratDokumenRepository
	Config                  *configs.Config
}

func ProvideSyaratDokumenServiceImpl(repository SyaratDokumenRepository, config *configs.Config) *SyaratDokumenServiceImpl {
	s := new(SyaratDokumenServiceImpl)
	s.SyaratDokumenRepository = repository
	s.Config = config
	return s
}

func (s *SyaratDokumenServiceImpl) Create(req RequestSyaratDokumen, userId uuid.UUID, tenantId uuid.UUID) (data SyaratDokumen, err error) {

	// existNama, err := s.SyaratDokumenRepository.ExistByNama(req.IdDokumen, "", *req.IdBranch)
	// if existNama {
	// 	return SyaratDokumen{}, errors.New("Nama SyaratDokumen sudah dipakai")
	// }

	data, _ = data.SyaratDokumenFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.SyaratDokumenRepository.Create(data)
	if err != nil {
		return SyaratDokumen{}, err
	}
	return data, nil
}

func (s *SyaratDokumenServiceImpl) GetAll(req model.StandardRequest) (data []SyaratDokumen, err error) {
	return s.SyaratDokumenRepository.GetAll(req)
}

func (s *SyaratDokumenServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.SyaratDokumenRepository.ResolveAll(request)
}

func (s *SyaratDokumenServiceImpl) DeleteByID(id uuid.UUID) error {
	newSyaratDokumen, err := s.SyaratDokumenRepository.ResolveByID(id)

	if err != nil || (SyaratDokumen{}) == newSyaratDokumen {
		return errors.New("Data SyaratDokumen dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.SyaratDokumenRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data SyaratDokumen dengan ID: " + id.String())
	}
	return nil
}

func (s *SyaratDokumenServiceImpl) Update(req RequestSyaratDokumen, userId uuid.UUID, tenantId uuid.UUID) (data SyaratDokumen, err error) {

	// existNama, err := s.SyaratDokumenRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	// if existNama {
	// 	return SyaratDokumen{}, errors.New("Nama SyaratDokumen sudah dipakai")
	// }

	data, _ = data.SyaratDokumenFormatRequest(req, userId, tenantId)
	err = s.SyaratDokumenRepository.Update(data)
	if err != nil {
		return SyaratDokumen{}, err
	}
	return data, nil
}

func (s *SyaratDokumenServiceImpl) ResolveByID(id uuid.UUID) (data SyaratDokumen, err error) {
	return s.SyaratDokumenRepository.ResolveByID(id)
}

func (s *SyaratDokumenServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newSyaratDokumen, err := s.SyaratDokumenRepository.ResolveByID(id)

	if err != nil || (SyaratDokumen{}) == newSyaratDokumen {
		return errors.New("Data SyaratDokumen dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.SyaratDokumenRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data SyaratDokumen dengan Nama :" + *&newSyaratDokumen.Nama + " sedang digunakan")
	}

	newSyaratDokumen.SoftDelete(userId)
	err = s.SyaratDokumenRepository.Update(newSyaratDokumen)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data SyaratDokumen dengan ID: " + id.String())
	}
	return nil
}
