package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type JenisPerjalananDinasService interface {
	Create(req JenisPerjalananDinasFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisPerjalananDinas JenisPerjalananDinas, err error)
	GetAll(req model.StandardRequest) (data []JenisPerjalananDinas, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	Update(req JenisPerjalananDinasFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisPerjalananDinas JenisPerjalananDinas, err error)
	ResolveByID(id uuid.UUID) (data JenisPerjalananDinas, err error)
	DeleteSoft(id uuid.UUID, userId uuid.UUID) error
}

type JenisPerjalananDinasServiceImpl struct {
	JenisPerjalananDinasRepository JenisPerjalananDinasRepository
	Config                         *configs.Config
}

func ProvideJenisPerjalananDinasServiceImpl(repository JenisPerjalananDinasRepository) *JenisPerjalananDinasServiceImpl {
	s := new(JenisPerjalananDinasServiceImpl)
	s.JenisPerjalananDinasRepository = repository
	return s
}

func (s *JenisPerjalananDinasServiceImpl) Create(req JenisPerjalananDinasFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisPerjalananDinas JenisPerjalananDinas, err error) {
	exist, err := s.JenisPerjalananDinasRepository.ExistByNama(req.Nama, *req.IdBranch)
	if exist {
		x := errors.New("Nama JenisPerjalananDinas sudah dipakai")
		return JenisPerjalananDinas{}, x
	}
	if err != nil {
		return JenisPerjalananDinas{}, err
	}
	newJenisPerjalananDinas, _ = newJenisPerjalananDinas.JenisPerjalananDinasFormat(req, userId, tenantId)
	err = s.JenisPerjalananDinasRepository.Create(newJenisPerjalananDinas)
	if err != nil {
		return JenisPerjalananDinas{}, err
	}
	return newJenisPerjalananDinas, nil
}

func (s *JenisPerjalananDinasServiceImpl) GetAll(req model.StandardRequest) (data []JenisPerjalananDinas, err error) {
	return s.JenisPerjalananDinasRepository.GetAll(req)
}

func (s *JenisPerjalananDinasServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.JenisPerjalananDinasRepository.ResolveAll(request)
}

func (s *JenisPerjalananDinasServiceImpl) DeleteByID(id uuid.UUID) error {
	jenisPerjalananDinas, err := s.JenisPerjalananDinasRepository.ResolveByID(id)

	if err != nil || (JenisPerjalananDinas{}) == jenisPerjalananDinas {
		return errors.New("Data JenisPerjalananDinas dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.JenisPerjalananDinasRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JenisPerjalananDinas dengan ID: " + id.String())
	}
	return nil
}

func (s *JenisPerjalananDinasServiceImpl) Update(req JenisPerjalananDinasFormat, userId uuid.UUID, tenantId uuid.UUID) (newJenisPerjalananDinas JenisPerjalananDinas, err error) {
	exist, err := s.JenisPerjalananDinasRepository.ExistByNamaID(req.ID, req.Nama, *req.IdBranch)
	if exist {
		x := errors.New("Nama JenisPerjalananDinas sudah dipakai")
		return JenisPerjalananDinas{}, x
	}
	if err != nil {
		return JenisPerjalananDinas{}, err
	}
	newJenisPerjalananDinas, _ = newJenisPerjalananDinas.JenisPerjalananDinasFormat(req, userId, tenantId)
	err = s.JenisPerjalananDinasRepository.Update(newJenisPerjalananDinas)
	if err != nil {
		return JenisPerjalananDinas{}, err
	}
	return newJenisPerjalananDinas, nil
}

func (s *JenisPerjalananDinasServiceImpl) ResolveByID(id uuid.UUID) (newJenisPerjalananDinas JenisPerjalananDinas, err error) {
	return s.JenisPerjalananDinasRepository.ResolveByID(id)
}

func (s *JenisPerjalananDinasServiceImpl) DeleteSoft(id uuid.UUID, userId uuid.UUID) error {
	newJenisPerjalananDinas, err := s.JenisPerjalananDinasRepository.ResolveByID(id)

	if err != nil || (JenisPerjalananDinas{}) == newJenisPerjalananDinas {
		return errors.New("Data JenisPerjalananDinas dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.JenisPerjalananDinasRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New("Data Jenis Perjalanan Dinas dengan Nama :" + *&newJenisPerjalananDinas.Nama + " sedang digunakan")
	}

	newJenisPerjalananDinas.SoftDelete(userId)
	err = s.JenisPerjalananDinasRepository.Update(newJenisPerjalananDinas)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JenisPerjalananDinas dengan ID: " + id.String())
	}
	return nil
}
