package bpd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type PerjalananDinasBiayaService interface {
	CreateBulk(reqFormat RequestPerjalananDinasBiaya, userID uuid.UUID) (data []PerjalananDinasBiaya, err error)
	CreateBulkUm(reqFormat RequestPerjalananDinasBiaya, userID uuid.UUID) (data []PerjalananDinasBiaya, err error)
	GetAllData(idBpdPegawai string) (data []PerjalananDinasBiayaDTO, err error)
	GetAllDataUm(idBpdPegawai string) (data []PerjalananDinasBiayaDTO, err error)
	UploadDocPenyelesaianBpd(req DocPenyelesaianBpdPegawai, userID string) (data DocPenyelesaianBpdPegawai, err error)
	UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error)
	Create(req PerjalananDinasBiayaDetail, userId uuid.UUID, jenis string) (data PerjalananDinasBiaya, error error)
	GetBiayaDto(idBpdPegawai string, idPegawai string, isReimbursement string) (data []BiayaPerjalananDinasDto, err error)
	DeleteByID(id uuid.UUID) error
	Update(req PerjalananDinasBiayaDetail, userId uuid.UUID, jenis string) (data PerjalananDinasBiaya, error error)
	ResolveByID(id uuid.UUID) (data BiayaPerjalananDinasDto, err error)
	GetHistoriBiaya(idBpdPegawai string, idPegawai string, idJenisBiaya string, isReimbursement string) (data []HistoriPerjalananDinas, err error)
	SoftDelete(id uuid.UUID, userId uuid.UUID) error
}

type PerjalananDinasBiayaServiceImpl struct {
	PerjalananDinasBiayaRepository PerjalananDinasBiayaRepository
	PerjalananDinasRepository      PerjalananDinasRepository
	Config                         *configs.Config
}

func ProvidePerjalananDinasBiayaServiceImpl(repository PerjalananDinasBiayaRepository, pdRepository PerjalananDinasRepository, config *configs.Config) *PerjalananDinasBiayaServiceImpl {
	s := new(PerjalananDinasBiayaServiceImpl)
	s.PerjalananDinasBiayaRepository = repository
	s.PerjalananDinasRepository = pdRepository
	s.Config = config
	return s
}

func (s *PerjalananDinasBiayaServiceImpl) CreateBulk(reqFormat RequestPerjalananDinasBiaya, userID uuid.UUID) (data []PerjalananDinasBiaya, err error) {
	var pdBiaya PerjalananDinasBiaya
	data, _ = pdBiaya.PerjalananDinasBiayaFormatRequest(reqFormat, userID)
	err = s.PerjalananDinasBiayaRepository.CreateBulk(data)
	if err != nil {
		return []PerjalananDinasBiaya{}, err
	}

	// update is_revisi when true
	if reqFormat.IsRevisi && len(reqFormat.Data) > 0 {
		idBpdPegawai := reqFormat.Data[0].IDBpdPegawai
		pd, err := s.PerjalananDinasRepository.ResolveBpdPegawaiByID(idBpdPegawai)
		if err != nil {
			return []PerjalananDinasBiaya{}, errors.New("Data Perjalanan Dinas Pegawai dengan ID :" + idBpdPegawai + " tidak ditemukan")
		}

		pd.IsRevisi = &reqFormat.IsRevisi
		err = s.PerjalananDinasRepository.UpdateBpdPegawai(pd)
		if err != nil {
			return []PerjalananDinasBiaya{}, err
		}
	}

	// update uang muka
	um := UangMukaBpd{
		ID:             reqFormat.IdBpdPegawai,
		IsUm:           reqFormat.IsUm,
		PersentaseUm:   reqFormat.PersentaseUm,
		PersentaseSisa: reqFormat.PersentaseSisa,
		ShowUm:         reqFormat.ShowUm,
		ShowSisa:       reqFormat.ShowSisa,
		TotalUm:        reqFormat.TotalUm,
		SisaUm:         reqFormat.SisaUm,
	}

	err = s.PerjalananDinasBiayaRepository.UpdateUangMuka(um)
	if err != nil {
		return []PerjalananDinasBiaya{}, err
	}

	return data, nil
}

func (s *PerjalananDinasBiayaServiceImpl) CreateBulkUm(reqFormat RequestPerjalananDinasBiaya, userID uuid.UUID) (data []PerjalananDinasBiaya, err error) {
	var pdBiaya PerjalananDinasBiaya
	data, _ = pdBiaya.PerjalananDinasBiayaFormatRequest(reqFormat, userID)
	err = s.PerjalananDinasBiayaRepository.CreateBulkUm(data)
	if err != nil {
		return []PerjalananDinasBiaya{}, err
	}

	// if true auto insert biaya
	if reqFormat.InsertBiaya {
		err = s.PerjalananDinasBiayaRepository.CreateBulk(data)
		if err != nil {
			return []PerjalananDinasBiaya{}, err
		}
	}

	// update is_revisi when true
	if reqFormat.IsRevisi && len(reqFormat.Data) > 0 {
		idBpdPegawai := reqFormat.Data[0].IDBpdPegawai
		pd, err := s.PerjalananDinasRepository.ResolveBpdPegawaiByID(idBpdPegawai)
		if err != nil {
			return []PerjalananDinasBiaya{}, errors.New("Data Perjalanan Dinas Pegawai dengan ID :" + idBpdPegawai + " tidak ditemukan")
		}

		pd.IsRevisi = &reqFormat.IsRevisi
		err = s.PerjalananDinasRepository.UpdateBpdPegawai(pd)
		if err != nil {
			return []PerjalananDinasBiaya{}, err
		}
	}

	// update uang muka
	um := UangMukaBpd{
		ID:             reqFormat.IdBpdPegawai,
		IsUm:           reqFormat.IsUm,
		PersentaseUm:   reqFormat.PersentaseUm,
		PersentaseSisa: reqFormat.PersentaseSisa,
		ShowUm:         reqFormat.ShowUm,
		ShowSisa:       reqFormat.ShowSisa,
		TotalUm:        reqFormat.TotalUm,
		SisaUm:         reqFormat.SisaUm,
	}

	err = s.PerjalananDinasBiayaRepository.UpdateUangMuka(um)
	if err != nil {
		return []PerjalananDinasBiaya{}, err
	}

	return data, nil
}

func (s *PerjalananDinasBiayaServiceImpl) GetAllData(idBpdPegawai string) (data []PerjalananDinasBiayaDTO, err error) {
	return s.PerjalananDinasBiayaRepository.GetAllData(idBpdPegawai)
}

func (s *PerjalananDinasBiayaServiceImpl) GetAllDataUm(idBpdPegawai string) (data []PerjalananDinasBiayaDTO, err error) {
	return s.PerjalananDinasBiayaRepository.GetAllDataUm(idBpdPegawai)
}

func (s *PerjalananDinasBiayaServiceImpl) UploadDocPenyelesaianBpd(req DocPenyelesaianBpdPegawai, userID string) (data DocPenyelesaianBpdPegawai, err error) {
	err = s.PerjalananDinasBiayaRepository.UploadDocPenyelesaianBpd(req)
	if err != nil {
		return DocPenyelesaianBpdPegawai{}, err
	}

	data = req
	return data, nil
}

func (s *PerjalananDinasBiayaServiceImpl) UploadFile(w http.ResponseWriter, r *http.Request, formValue string, pathFile string) (path string, err error) {
	if err = r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadedFile, handler, err := r.FormFile(formValue)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	newID, _ := uuid.NewV4()
	filename := fmt.Sprintf("%s%s", "penyelesaian_bpd_"+newID.String(), filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir
	DokumenDir := s.Config.App.File.Bpd

	if pathFile == "" {
		path = filepath.Join(DokumenDir, filename)
	} else {
		path = pathFile
	}
	fileLocation := filepath.Join(dir, path)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("ERROR FILE:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err = io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("ERROR COPY FILE:", err)
		return
	}
	return
}

func (s *PerjalananDinasBiayaServiceImpl) Create(req PerjalananDinasBiayaDetail, userId uuid.UUID, jenis string) (data PerjalananDinasBiaya, err error) {
	if jenis == "AK" {
		existNama, _ := s.PerjalananDinasBiayaRepository.ExistAkomodasi(req.IDJenisBiaya, req.IDBpdPegawai, "")
		if existNama {
			return PerjalananDinasBiaya{}, errors.New("Jenis Akomodasi sudah dipakai")
		}
	}
	data, err = data.BiayaRequest(req, userId)
	if err != nil {
		return PerjalananDinasBiaya{}, err
	}

	err = s.PerjalananDinasBiayaRepository.Create(data, jenis)
	if err != nil {
		return PerjalananDinasBiaya{}, err
	}
	return data, nil
}

func (s *PerjalananDinasBiayaServiceImpl) GetBiayaDto(idBpdPegawai string, idPegawai string, isReimbursement string) (data []BiayaPerjalananDinasDto, err error) {
	return s.PerjalananDinasBiayaRepository.GetBiayaDto(idBpdPegawai, idPegawai, isReimbursement)
}

func (s *PerjalananDinasBiayaServiceImpl) DeleteByID(id uuid.UUID) error {
	err := s.PerjalananDinasBiayaRepository.DeleteByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Perjalanan Dinas Biaya dengan ID: " + id.String())
	}
	return nil
}

func (s *PerjalananDinasBiayaServiceImpl) Update(req PerjalananDinasBiayaDetail, userId uuid.UUID, jenis string) (data PerjalananDinasBiaya, err error) {

	if jenis == "AK" {
		existNama, _ := s.PerjalananDinasBiayaRepository.ExistAkomodasi(req.IDJenisBiaya, req.IDBpdPegawai, req.ID)
		if existNama {
			return PerjalananDinasBiaya{}, errors.New("Jenis Akomodasi sudah dipakai")
		}
		fmt.Println("existnama", existNama)
	}
	data, err = data.BiayaRequest(req, userId)
	if err != nil {
		return PerjalananDinasBiaya{}, err
	}
	err = s.PerjalananDinasBiayaRepository.Update(data, jenis)
	if err != nil {
		return PerjalananDinasBiaya{}, err
	}
	return data, nil
}

func (s *PerjalananDinasBiayaServiceImpl) ResolveByID(id uuid.UUID) (data BiayaPerjalananDinasDto, err error) {
	return s.PerjalananDinasBiayaRepository.ResolveByID(id)
}

func (s *PerjalananDinasBiayaServiceImpl) GetHistoriBiaya(idBpdPegawai string, idPegawai string, idJenisBiaya string, isReimbursement string) (data []HistoriPerjalananDinas, err error) {
	return s.PerjalananDinasBiayaRepository.GetHistoriBiaya(idBpdPegawai, idPegawai, idJenisBiaya, isReimbursement)
}

func (s *PerjalananDinasBiayaServiceImpl) SoftDelete(id uuid.UUID, userId uuid.UUID) error {
	newPerjalananDinasBiaya, err := s.PerjalananDinasBiayaRepository.ResolveByIdBiaya(id)

	if err != nil || (PerjalananDinasBiaya{}) == newPerjalananDinasBiaya {
		return errors.New("Data Perjalanan Dins Biaya dengan ID :" + id.String() + " tidak ditemukan")
	}

	newPerjalananDinasBiaya.SoftDeleteBiaya(userId.String())
	err = s.PerjalananDinasBiayaRepository.Update(newPerjalananDinasBiaya, "REM")
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Perjalanan Dinas Biaya dengan ID: " + id.String())
	}
	return nil
}
