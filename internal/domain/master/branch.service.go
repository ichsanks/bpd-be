package master

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/transport/http/response"
)

type BranchService interface {
	Create(reqFormat RequestBranchFormat, userID uuid.UUID) (data Branch, err error)
	Update(reqFormat RequestBranchFormat, userID uuid.UUID) (data Branch, err error)
	GetAllData() (data []Branch, err error)
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Branch, err error)
	ResolveByIDDTO(id uuid.UUID) (data BranchDTO, err error)
	DeleteByID(id uuid.UUID, userID uuid.UUID) error
	UploadFile(w http.ResponseWriter, r *http.Request, path_file string) (path string, err error)
	DeleteFile(path string) (err error)
}

type BranchServiceImpl struct {
	BranchRepository BranchRepository
	Config           *configs.Config
}

func ProvideBranchServiceImpl(repository BranchRepository, config *configs.Config) *BranchServiceImpl {
	s := new(BranchServiceImpl)
	s.BranchRepository = repository
	s.Config = config
	return s
}

func (s *BranchServiceImpl) Create(reqFormat RequestBranchFormat, userID uuid.UUID) (data Branch, err error) {
	existCode, err := s.BranchRepository.ExistKode(reqFormat.Kode, "")
	if existCode {
		x := errors.New("Kode Sudah Ada")
		return Branch{}, x
	}

	exist, err := s.BranchRepository.ExistName(reqFormat.Nama, "")
	if exist {
		x := errors.New("Nama Sudah Ada")
		return Branch{}, x
	}
	data, _ = data.BranchFormat(reqFormat, userID)
	err = s.BranchRepository.Create(data)
	if err != nil {
		return Branch{}, err
	}

	return data, nil
}

func (s *BranchServiceImpl) Update(reqFormat RequestBranchFormat, userID uuid.UUID) (data Branch, err error) {
	existCode, err := s.BranchRepository.ExistKode(reqFormat.Kode, reqFormat.Id.String())
	if existCode {
		x := errors.New("Kode Sudah Ada")
		return Branch{}, x
	}

	exist, err := s.BranchRepository.ExistName(reqFormat.Nama, reqFormat.Id.String())
	if exist {
		x := errors.New("Nama Sudah Ada")
		return Branch{}, x
	}
	np, _ := data.BranchFormat(reqFormat, userID)
	err = s.BranchRepository.Update(np)
	if err != nil {
		log.Error().Msgf("service.UpdateBranch error", err)
	}
	return np, nil
}

func (s *BranchServiceImpl) GetAllData() (data []Branch, err error) {
	return s.BranchRepository.GetAllData()
}

func (s *BranchServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.BranchRepository.ResolveAll(request)
}

func (s *BranchServiceImpl) ResolveByID(id uuid.UUID) (data Branch, err error) {
	return s.BranchRepository.ResolveByID(id)
}

func (s *BranchServiceImpl) ResolveByIDDTO(id uuid.UUID) (data BranchDTO, err error) {
	return s.BranchRepository.ResolveByIDDTO(id)
}

func (s *BranchServiceImpl) DeleteByID(id uuid.UUID, userID uuid.UUID) error {
	branch, err := s.BranchRepository.ResolveByID(id)

	if err != nil || (Branch{}) == branch {
		return errors.New(" ID :" + id.String() + " Branch data  not found")
	}

	branch.SoftDelete(userID)
	err = s.BranchRepository.Update(branch)
	if err != nil {
		return errors.New("There is an error in delete branch data with ID: " + id.String())
	}
	return nil
}

func (s *BranchServiceImpl) UploadFile(w http.ResponseWriter, r *http.Request, path_file string) (path string, err error) {
	if err = r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uploadedFile, handler, err := r.FormFile("file")
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
	filename := fmt.Sprintf("%s%s", "image_branch_"+newID.String(), filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir
	DokumenDir := s.Config.App.File.ImageBranch

	if path_file == "" {
		path = filepath.Join(DokumenDir, filename)
	} else {
		path = path_file
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

func (s *BranchServiceImpl) DeleteFile(path string) (err error) {
	if path == "" {
		return nil
	}

	dir := s.Config.App.File.Dir
	fileLocation := filepath.Join(dir, path)
	err = os.Remove(fileLocation)
	return
}
