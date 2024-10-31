package auth

import (
	"github.com/gofrs/uuid"
)

type Role struct {
	ID               string `json:"id" db:"id"`
	Nama             string `json:"nama" db:"nama"`
	Keterangan       string `json:"keterangan" db:"keterangan"`
	IsViewDataAll    bool   `json:"isViewDataAll" db:"is_view_data_all"`
	IsChoosePegawai  bool   `json:"isChoosePegawai" db:"is_choose_pegawai"`
	IsChooseTerbatas bool   `json:"isChooseTerbatas" db:"is_choose_terbatas"`
}

type RequestRole struct {
	Nama             string `json:"nama" db:"nama"`
	Keterangan       string `json:"keterangan" db:"keterangan"`
	IsViewDataAll    bool   `json:"isViewDataAll" db:"is_view_data_all"`
	IsChoosePegawai  bool   `json:"isChoosePegawai" db:"is_choose_pegawai"`
	IsChooseTerbatas bool   `json:"isChooseTerbatas" db:"is_choose_terbatas"`
}

type RoleMenu struct {
	RoleID string `json:"role_id" db:"role_id" validate:"required"`
	Menu   string `json:"menu" db:"menu" validate:"required"`
}

type RoleMenuFormat struct {
	Role  Role
	Menus []string
}

type RoleMenuRequestFormat struct {
	Menus []string
}

func (role *Role) NewRoleFormat(reqFormat Role) (newRole Role, err error) {
	newID, _ := uuid.NewV4()
	newRole = Role{
		ID:               newID.String(),
		Nama:             reqFormat.Nama,
		Keterangan:       reqFormat.Keterangan,
		IsViewDataAll:    reqFormat.IsViewDataAll,
		IsChoosePegawai:  reqFormat.IsChoosePegawai,
		IsChooseTerbatas: reqFormat.IsChooseTerbatas,
	}
	return
}
