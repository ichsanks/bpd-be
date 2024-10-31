package auth

import (
	"errors"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/configs"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

// RoleService adalah interface RoleService untuk entity Role
type RoleService interface {
	GetData(request model.StandardRequest) (orders pagination.Response, err error)
	ResolveAll() ([]Role, error)
	CreateRole(role Role) (bool, error)
	UpdateRole(id string, role Role) error
	DeleteRole(id string) error
	ResolveByRoleID(id string) (RoleMenuFormat, error)
	ResolveMenuByRoleID(id string) ([]string, error)
	CreateOrUpdateMenuByRoleID(roleMenuFormat RoleMenuFormat) error
}

// RoleServiceImpl adalah implementasi dari service yang digunakan untuk entity Role
type RoleServiceImpl struct {
	RoleRepository RoleRepository
	Config         *configs.Config
}

// ProvideRoleServiceImpl adalah provider untuk service RoleService
func ProvideRoleServiceImpl(roleRepository RoleRepository, config *configs.Config) *RoleServiceImpl {
	s := new(RoleServiceImpl)
	s.RoleRepository = roleRepository
	s.Config = config
	return s
}

func (s *RoleServiceImpl) GetData(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.RoleRepository.GetData(request)
}

// ResolveAll get all Role data
func (s *RoleServiceImpl) ResolveAll() ([]Role, error) {
	return s.RoleRepository.ResolveAll()
}

// CreateRole is the service to create Role entity
func (r *RoleServiceImpl) CreateRole(role Role) (bool, error) {

	exist, err := r.RoleRepository.ExistRoleByNama(role.Nama)
	if exist {
		return exist, errors.New("Nama Role sudah dipakai")
	}
	if err != nil {
		return exist, err
	}
	newRole, _ := role.NewRoleFormat(role)
	err = r.RoleRepository.CreateRole(newRole)
	if err != nil {
		return exist, err
	}
	return exist, nil
}

// UpdateRole aalah service yang digunakan untuk mengubah data Role
func (r *RoleServiceImpl) UpdateRole(id string, newRole Role) error {
	role, err := r.RoleRepository.ResolveRoleByID(id)
	if err != nil || (Role{}) == role {
		return errors.New("Role dengan nama :" + role.Nama + " tidak ditemukan")
	}

	return r.RoleRepository.UpdateRole(newRole)
}

// DeleteRole adalah service yang digunakan untuk menghapus data Role
func (r *RoleServiceImpl) DeleteRole(id string) error {
	role, err := r.RoleRepository.ResolveRoleByID(id)
	if err != nil || (Role{}) == role {
		return errors.New("Role dengan nama :" + role.Nama + " tidak ditemukan")
	}
	err = r.RoleRepository.DeleteRoleByID(id)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data Role dengan nama: " + role.Nama)
	}
	return nil
}

// ResolveByRoleID adalah service yang digunakan untuk mendapatkan Role berdasarkan RoleID
func (r *RoleServiceImpl) ResolveByRoleID(roleId string) (roleMenu RoleMenuFormat, err error) {
	role, err := r.RoleRepository.ResolveRoleByID(roleId)
	if err != nil {
		return RoleMenuFormat{}, errors.New("Role dengan nama: " + role.Nama + " tidak ditemukan")
	}
	roleMenu.Role = role
	menus, err := r.RoleRepository.ResolveRoleMenuByRoleID(roleId)
	if err != nil {
		return RoleMenuFormat{}, errors.New("Ada kesalahan waktu get menu berdasarkan roleID: " + roleId)
	}
	roleMenu.Menus = make([]string, 0)
	for _, menu := range menus {
		roleMenu.Menus = append(roleMenu.Menus, menu)
	}
	return
}

// ResolveMenuByRoleID adalah service yang digunakan untuk mendapatkan menu berdasarkan RoleID
func (r *RoleServiceImpl) ResolveMenuByRoleID(roleId string) (menus []string, err error) {
	menus, err = r.RoleRepository.ResolveRoleMenuByRoleID(roleId)
	if err != nil {
		return make([]string, 0), errors.New("Ada kesalahan waktu get menu berdasarkan roleID: " + roleId)
	}
	return
}

// CreateOrUpdateMenuByRoleID digunakan untuk menambah atau update Role beserta menu yang diperbolehkan
func (r *RoleServiceImpl) CreateOrUpdateMenuByRoleID(roleMenuFormat RoleMenuFormat) (err error) {
	existRole, err := r.RoleRepository.ExistRoleByID(roleMenuFormat.Role.ID)
	if err != nil {
		return
	}
	if existRole == true {
		err = r.RoleRepository.UpdateRole(roleMenuFormat.Role)
		if err != nil {
			return
		}
	} else {
		_, err = r.CreateRole(roleMenuFormat.Role)
		if err != nil {
			return
		}
	}
	return r.RoleRepository.UpdateRoleMenu(roleMenuFormat.Role.ID, roleMenuFormat.Menus)
}
