package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type JabatanService interface {
	Create(req JabatanFormat, userId uuid.UUID, tenantId uuid.UUID) (newJabatan Jabatan, err error)
	GetAll(req model.StandardRequest) (data []Jabatan, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	Update(req JabatanFormat, userId uuid.UUID, tenantId uuid.UUID) (newJabatan Jabatan, err error)
	ResolveByID(id uuid.UUID) (data Jabatan, err error)
	DeleteSoft(id uuid.UUID, userId uuid.UUID) error
}

type JabatanServiceImpl struct {
	JabatanRepository JabatanRepository
	Config            *configs.Config
}

func ProvideJabatanServiceImpl(repository JabatanRepository) *JabatanServiceImpl {
	s := new(JabatanServiceImpl)
	s.JabatanRepository = repository
	return s
}

func (s *JabatanServiceImpl) Create(req JabatanFormat, userId uuid.UUID, tenantId uuid.UUID) (newJabatan Jabatan, err error) {
	exist, err := s.JabatanRepository.ExistByNama(req.Nama, *req.IdBranch)
	if exist {
		x := errors.New("Nama Jabatan sudah dipakai")
		return Jabatan{}, x
	}
	if err != nil {
		return Jabatan{}, err
	}
	newJabatan, _ = newJabatan.JabatanFormat(req, userId, tenantId)
	err = s.JabatanRepository.Create(newJabatan)
	if err != nil {
		return Jabatan{}, err
	}
	return newJabatan, nil
}

func (s *JabatanServiceImpl) GetAll(req model.StandardRequest) (data []Jabatan, err error) {
	return s.JabatanRepository.GetAll(req)
}

func (s *JabatanServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.JabatanRepository.ResolveAll(request)
}

func (s *JabatanServiceImpl) DeleteByID(id uuid.UUID) error {
	jabatan, err := s.JabatanRepository.ResolveByID(id)

	if err != nil || (Jabatan{}) == jabatan {
		return errors.New("Data Jabatan dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.JabatanRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Jabatan dengan ID: " + id.String())
	}
	return nil
}

func (s *JabatanServiceImpl) Update(req JabatanFormat, userId uuid.UUID, tenantId uuid.UUID) (newJabatan Jabatan, err error) {
	exist, err := s.JabatanRepository.ExistByNamaID(req.ID, req.Nama, *req.IdBranch)
	if exist {
		x := errors.New("Nama Jabatan sudah dipakai")
		return Jabatan{}, x
	}
	if err != nil {
		return Jabatan{}, err
	}
	newJabatan, _ = newJabatan.JabatanFormat(req, userId, tenantId)
	err = s.JabatanRepository.Update(newJabatan)
	if err != nil {
		return Jabatan{}, err
	}
	return newJabatan, nil
}

func (s *JabatanServiceImpl) ResolveByID(id uuid.UUID) (newJabatan Jabatan, err error) {
	return s.JabatanRepository.ResolveByID(id)
}

func (s *JabatanServiceImpl) DeleteSoft(id uuid.UUID, userId uuid.UUID) error {
	newJabatan, err := s.JabatanRepository.ResolveByID(id)

	if err != nil || (Jabatan{}) == newJabatan {
		return errors.New("Data Jabatan dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.JabatanRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data Jabatan dengan Nama :" + *&newJabatan.Nama + " sedang digunakan")
	}

	newJabatan.SoftDelete(userId)
	err = s.JabatanRepository.Update(newJabatan)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Jabatan dengan ID: " + id.String())
	}
	return nil
}
