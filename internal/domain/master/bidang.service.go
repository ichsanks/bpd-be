package master

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

type BidangService interface {
	Create(req RequestBidang, userId uuid.UUID, tenantId uuid.UUID) (data Bidang, err error)
	Update(req RequestBidang, userId uuid.UUID, tenantId uuid.UUID) (data Bidang, err error)
	GetAll(req model.StandardRequest) (data []Bidang, err error)
	DeleteByID(id uuid.UUID) error
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Bidang, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type BidangServiceImpl struct {
	BidangRepository BidangRepository
	Config           *configs.Config
}

func ProvideBidangServiceImpl(repository BidangRepository, config *configs.Config) *BidangServiceImpl {
	s := new(BidangServiceImpl)
	s.BidangRepository = repository
	s.Config = config
	return s
}

func (s *BidangServiceImpl) Create(req RequestBidang, userId uuid.UUID, tenantId uuid.UUID) (data Bidang, err error) {
	existKode, err := s.BidangRepository.ExistByKode(req.Kode, "", *req.IdBranch)
	if existKode {
		return Bidang{}, errors.New("Kode Bidang sudah dipakai")
	}

	existNama, err := s.BidangRepository.ExistByNama(req.Nama, "", *req.IdBranch)
	if existNama {
		return Bidang{}, errors.New("Nama Bidang sudah dipakai")
	}

	data, _ = data.BidangFormatRequest(req, userId, tenantId)
	if err != nil {
		return
	}

	err = s.BidangRepository.Create(data)
	if err != nil {
		return Bidang{}, err
	}
	return data, nil
}

func (s *BidangServiceImpl) GetAll(req model.StandardRequest) (data []Bidang, err error) {
	return s.BidangRepository.GetAll(req)
}

func (s *BidangServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.BidangRepository.ResolveAll(request)
}

func (s *BidangServiceImpl) DeleteByID(id uuid.UUID) error {
	newBidang, err := s.BidangRepository.ResolveByID(id)

	if err != nil || (Bidang{}) == newBidang {
		return errors.New("Data Bidang dengan ID :" + id.String() + " tidak ditemukan")
	}

	err = s.BidangRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Bidang dengan ID: " + id.String())
	}
	return nil
}

func (s *BidangServiceImpl) Update(req RequestBidang, userId uuid.UUID, tenantId uuid.UUID) (data Bidang, err error) {
	existKode, _ := s.BidangRepository.ExistByKode(req.Kode, req.ID.String(), *req.IdBranch)
	if existKode {
		return Bidang{}, errors.New("Kode Bidang sudah dipakai")
	}

	existNama, _ := s.BidangRepository.ExistByNama(req.Nama, req.ID.String(), *req.IdBranch)
	if existNama {
		return Bidang{}, errors.New("Nama Bidang sudah dipakai")
	}

	data, _ = data.BidangFormatRequest(req, userId, tenantId)
	err = s.BidangRepository.Update(data)
	if err != nil {
		return Bidang{}, err
	}
	return data, nil
}

func (s *BidangServiceImpl) ResolveByID(id uuid.UUID) (data Bidang, err error) {
	return s.BidangRepository.ResolveByID(id)
}

func (s *BidangServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newBidang, err := s.BidangRepository.ResolveByID(id)

	if err != nil || (Bidang{}) == newBidang {
		return errors.New("Data Bidang dengan ID :" + id.String() + " tidak ditemukan")
	}
	exist := s.BidangRepository.ExistRelasiStatus(id)
	if exist {
		return errors.New(" Data Bidang dengan Nama :" + *&newBidang.Nama + " sedang digunakan")
	}

	newBidang.SoftDelete(userId)
	err = s.BidangRepository.Update(newBidang)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Bidang dengan ID: " + id.String())
	}
	return nil
}
