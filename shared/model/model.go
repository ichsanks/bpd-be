package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid"
)

const RECORD_NOT_FOUND = "record not found"

// StandardRequest is a standard query string request
type StandardRequest struct {
	Keyword           string `json:"q" validate:"omitempty"`
	StartDate         string `json:"startDate" validate:"omitempty"`
	EndDate           string `json:"endDate" validate:"omitempty"`
	PageNumber        int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize          int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy            string `json:"sortBy" validate:"required"`
	SortType          string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	Status            string `json:"status" validate:"omitempty"`
	IgnorePaging      bool   `json:"ignorePaging" validate:"omitempty"`
	IdPerjalananDinas string `json:"idPerjalananDinas" validate:"omitempty"`
	IdPegawai         string `json:"idPegawai" validate:"omitempty"`
	IdPegawaiApproval string `json:"idPegawaiApproval" validate:"omitempty"`
	IdBidang          string `json:"idBidang" validate:"omitempty"`
	IdFungsionalitas  string `json:"idFungsionalitas" validate:"omitempty"`
	IdUnor            string `json:"idUnor" validate:"omitempty"`
	TypeApproval      string `json:"typeApproval" validate:"omitempty"`
	IsSppb            string `json:"isSppb" validate:"omitempty"`
	IdBranch          string `json:"idBranch" validate:"omitempty"`
	TenantID          string `json:"tenantId" validate:"omitempty"`
	IdBodLevel        string `json:"idBodLevel" validate:"omitempty"`
	IdTransaksi       string `json:"idTransaksi" validate:"omitempty"`
	IdJenisTujuan     string `json:"idJenisTujuan" validate:"omitempty"`
}
type StandardModel struct {
	ID     string `json:"id" db:"id"`
	Nama   string `json:"nama" db:"nama"`
	Alamat string `json:"alamat" db:"alamat"`
}
type ReportRequestParams struct {
	IDOpd     uuid.UUID `json:"id" db:"id"`
	IdItem    string    `json:"idItem" validate:"omitempty"`
	IdBidang  string    `json:"idBidang" validate:"omitempty"`
	StartDate string    `json:"startDate" validate:"omitempty"`
	EndDate   string    `json:"endDate" validate:"omitempty"`
	Status    string    `json:"status"`
}

type StandardRequestUser struct {
	Keyword    string `json:"q" validate:"omitempty"`
	StartDate  string `json:"startDate" validate:"omitempty"`
	EndDate    string `json:"endDate" validate:"omitempty"`
	PageNumber int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize   int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy     string `json:"sortBy" validate:"required"`
	SortType   string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	IdRole     string `json:"idRole" validate:"omitempty"`
	IdUnor     string `json:"idUnor" validate:"omitempty"`
	IdBidang   string `json:"idBidang" validate:"omitempty"`
	IdBranch   string `json:"idBranch" validate:"omitempty"`
	TenantID   string `json:"tenantId" validate:"omitempty"`
}

type StandardRequestMenu struct {
	Keyword    string `json:"q" validate:"omitempty"`
	PageNumber int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize   int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy     string `json:"sortBy" validate:"required"`
	SortType   string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	App        string `json:"app" validate:"omitempty"`
}

type StandardRequestKendaraan struct {
	Keyword          string `json:"q" validate:"omitempty"`
	PageNumber       int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize         int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy           string `json:"sortBy" validate:"required"`
	SortType         string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	IdJenisKendaraan string `json:"idJenisKendaraan" validate:"omitempty"`
}

type StandardRequestPegawai struct {
	Keyword    string `json:"q" validate:"omitempty"`
	IDPegawai  string `json:"idPegawai" validate:"omitempty"`
	StartDate  string `json:"startDate" validate:"omitempty"`
	EndDate    string `json:"endDate" validate:"omitempty"`
	PageNumber int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize   int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy     string `json:"sortBy" validate:"required"`
	SortType   string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	IdDivisi   string `json:"idDivisi" validate:"omitempty"`
}

type StandardRequestRuleApproval struct {
	Keyword    string `json:"q" validate:"omitempty"`
	StartDate  string `json:"startDate" validate:"omitempty"`
	EndDate    string `json:"endDate" validate:"omitempty"`
	PageNumber int    `json:"pageNumber" validate:"omitempty,gte=0"`
	PageSize   int    `json:"pageSize" validate:"omitempty,gte=0"`
	SortBy     string `json:"sortBy" validate:"required"`
	SortType   string `json:"sortType" validate:"required,oneof=asc ASC desc DESC"`
	Jenis      string `json:"jenis" validate:"omitempty"`
	IdBranch   string `json:"idBranch" validate:"omitempty"`
}

// JSONRaw ...
type JSONRaw json.RawMessage

// Value ...
func (j JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(j)

	return driver.Value(byteArr), nil
}

// Scan ...
func (j *JSONRaw) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	err := json.Unmarshal(asBytes, &j)
	if err != nil {
		return error(errors.New("Scan could not unmarshal to []string"))
	}

	return nil
}

// MarshalJSON ...
func (j *JSONRaw) MarshalJSON() ([]byte, error) {
	return *j, nil
}

// UnmarshalJSON ...
func (j *JSONRaw) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}
