package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type JenisTujuanService interface {
	Create(req RequestJenisTujuan, userId uuid.UUID, tenanId uuid.UUID) (data JenisTujuan, err error)
	Update(req RequestJenisTujuan, userId uuid.UUID, tenanId uuid.UUID) (data JenisTujuan, err error)
	GetAll(req model.StandardRequest) (data []JenisTujuan, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data JenisTujuan, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type JenisTujuanServiceImpl struct {
	JenisTujuanRepository JenisTujuanRepository
	Config                *configs.Config
}

func ProvideJenisTujuanServiceImpl(repository JenisTujuanRepository, config *configs.Config) *JenisTujuanServiceImpl {
	s := new(JenisTujuanServiceImpl)
	s.JenisTujuanRepository = repository
	s.Config = config
	return s
}

func (s *JenisTujuanServiceImpl) Create(req RequestJenisTujuan, userId uuid.UUID, tenanId uuid.UUID) (data JenisTujuan, err error) {
	existNama, err := s.JenisTujuanRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return JenisTujuan{}, errors.New("Nama JenisTujuan sudah dipakai")
	}

	data, _ = data.JenisTujuanFormatRequest(req, userId, tenanId)
	if err != nil {
		return
	}

	err = s.JenisTujuanRepository.Create(data)
	if err != nil {
		return JenisTujuan{}, err
	}
	return data, nil
}

func (s *JenisTujuanServiceImpl) GetAll(req model.StandardRequest) (data []JenisTujuan, err error) {
	return s.JenisTujuanRepository.GetAll(req)
}

func (s *JenisTujuanServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.JenisTujuanRepository.ResolveAll(request)
}

func (s *JenisTujuanServiceImpl) DeleteByID(id uuid.UUID) error {
	newJenisTujuan, err := s.JenisTujuanRepository.ResolveByID(id)

	if err != nil || (JenisTujuan{}) == newJenisTujuan {
		return errors.New("Data JenisTujuan dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.JenisTujuanRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JenisTujuan dengan ID: " + id.String())
	}
	return nil
}

func (s *JenisTujuanServiceImpl) Update(req RequestJenisTujuan, userId uuid.UUID, tenanId uuid.UUID) (data JenisTujuan, err error) {
	existNama, err := s.JenisTujuanRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return JenisTujuan{}, errors.New("Nama JenisTujuan sudah dipakai")
	}

	data, _ = data.JenisTujuanFormatRequest(req, userId, tenanId)
	err = s.JenisTujuanRepository.Update(data)
	if err != nil {
		return JenisTujuan{}, err
	}
	return data, nil
}

func (s *JenisTujuanServiceImpl) ResolveByID(id uuid.UUID) (data JenisTujuan, err error) {
	return s.JenisTujuanRepository.ResolveByID(id)
}

func (s *JenisTujuanServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newJenisTujuan, err := s.JenisTujuanRepository.ResolveByID(id)

	if err != nil || (JenisTujuan{}) == newJenisTujuan {
		return errors.New("Data JenisTujuan dengan ID :" + id.String() + " tidak ditemukan")
	}

	exist := s.JenisTujuanRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data JenisTujuan dengan Nama :" + *&newJenisTujuan.Nama + " sedang digunakan")
	}

	newJenisTujuan.SoftDelete(userId)
	err = s.JenisTujuanRepository.Update(newJenisTujuan)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data JenisTujuan dengan ID: " + id.String())
	}
	return nil
}
