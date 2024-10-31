package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"github.com/rs/zerolog/log"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

// MenuService adalah interface MenuService untuk entity menu
type MenuService interface {
	GetAllMenu(app string) ([]Menu, error)
	ResolveAll(request model.StandardRequestMenu) (orders pagination.Response, err error)
	ResolveMenuByRoleID(req RequestMenuUserFilter) (menu []MenuResponse, err error)
	ResolveAllMenuUser(req RequestMenuUserFilter) (menu []MenuResponse, err error)
	CreateMenu(reqFormat RequestMenuFormat) (menu Menu, error error)
	UpdateMenu(id uuid.UUID, newMenu RequestMenuFormat) (menu Menu, error error)
	ResolveMenuByID(id uuid.UUID) (menu Menu, error error)
	DeleteMenuByID(id uuid.UUID) error
	CreateMenuUser(reqFormat RequestMenuUserFormat) (newMenuUser []MenuUser, err error)
	ResolveMenuUserByID(id uuid.UUID) (menuUser MenuUser, error error)
	DeleteMenuUserByID(id uuid.UUID) error
	SortMenu(sortMenu RequestMenuSortFormat) (menuUser MenuUser, err error)
}

// MenuServiceImpl adalah implementasi dari service yang digunakan untuk entity Role
type MenuServiceImpl struct {
	MenuRepository MenuRepository
}

// ProvideServiceImpl adalah provider untuk service MenuService
func ProvideMenuServiceImpl(MenuRepository MenuRepository) *MenuServiceImpl {
	s := new(MenuServiceImpl)
	s.MenuRepository = MenuRepository
	return s
}

// ResolveByMenuRoleID adalah service yang digunakan untuk mendapatkan menu berdasarkan RoleID
func (r *MenuServiceImpl) ResolveMenuByRoleID(req RequestMenuUserFilter) (menu []MenuResponse, err error) {
	menu, err = r.MenuRepository.ResolveMenuByRoleID(req)
	if err != nil {
		return []MenuResponse{}, errors.New("Ada kesalahan waktu get menu berdasarkan roleID: " + req.IdRole)
	}

	for i := 0; i < len(menu); i++ {
		var reqDetail = RequestMenuUserFilter{
			IdRole:        req.IdRole,
			IdMenu:        menu[i].IDMenu,
			App:           req.App,
			PosisiSubMenu: req.PosisiSubMenu,
			IDBranch:      req.IDBranch,
			TenantID:      req.TenantID,
		}
		var CMenu []MenuResponse
		CMenu, err = r.MenuRepository.ResolveMenuByParentID(reqDetail)
		if err != nil {
			logger.ErrorWithStack(err)
			return
		}
		menu[i].Children = CMenu
	}

	return
}

// ResolveByMenuRoleID adalah service yang digunakan untuk mendapatkan menu berdasarkan RoleID
func (r *MenuServiceImpl) ResolveAllMenuUser(req RequestMenuUserFilter) (menu []MenuResponse, err error) {
	menu, err = r.MenuRepository.ResolveMenuByRoleID(req)
	if err != nil {
		return []MenuResponse{}, errors.New("Ada kesalahan waktu get menu berdasarkan roleID: " + req.IdRole)
	}

	return
}

func (s *MenuServiceImpl) GetAllMenu(app string) (data []Menu, err error) {
	return s.MenuRepository.GetAllMenu(app)
}

func (s *MenuServiceImpl) ResolveAll(request model.StandardRequestMenu) (orders pagination.Response, err error) {
	return s.MenuRepository.ResolveAll(request)
}

// CreateMenu is the service to create Menu entity
func (s *MenuServiceImpl) CreateMenu(reqFormat RequestMenuFormat) (newMenu Menu, err error) {
	if err != nil {
		return Menu{}, err
	}
	newMenu, err = newMenu.NewMenuFormat(reqFormat)

	err = s.MenuRepository.CreateMenu(newMenu)
	if err != nil {
		return Menu{}, err
	}
	return newMenu, nil
}

func (s *MenuServiceImpl) ResolveMenuByID(id uuid.UUID) (menu Menu, err error) {
	menu, err = s.MenuRepository.ResolveMenuByID(id)
	if err != nil {
		return
	}
	return
}

// UpdateMenu adalah service yang digunakan untuk mengupdate Menu
func (s *MenuServiceImpl) UpdateMenu(id uuid.UUID, newMenu RequestMenuFormat) (menu Menu, err error) {
	menu, err = s.MenuRepository.ResolveMenuByID(id)
	if err != nil {
		return Menu{}, errors.New("Menu dengan ID :" + id.String() + " tidak ditemukan")
	}
	menu.NewFormatUpdate(newMenu)
	log.Info().Msgf("service.UpdateMenu %s", menu)

	err = s.MenuRepository.UpdateMenu(menu)

	if err != nil {
		log.Error().Msgf("service.UpdateMenu error", err)
	}
	return
}

func (s *MenuServiceImpl) DeleteMenuByID(id uuid.UUID) error {
	menu, err := s.MenuRepository.ResolveMenuByID(id)
	if err != nil || (Menu{}) == menu {
		return errors.New("menu dengan ID :" + id.String() + " tidak ditemukan")
	}

	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data menu dengan ID: " + id.String())
	}
	menu.SoftDelete()

	err = s.MenuRepository.UpdateMenu(menu)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data menu dengan ID: " + id.String())
	}
	return nil
}

func (s *MenuServiceImpl) CreateMenuUser(reqFormat RequestMenuUserFormat) (newMenuUser []MenuUser, err error) {
	var menuUser MenuUser
	newMenuUser, err = menuUser.NewMenuUserFormat(reqFormat)
	for _, v := range newMenuUser {
		err = s.MenuRepository.CreateMenuUser(v)
	}

	if err != nil {
		return []MenuUser{}, err
	}

	return newMenuUser, nil
}

func (s *MenuServiceImpl) ResolveMenuUserByID(id uuid.UUID) (menuUser MenuUser, err error) {
	menuUser, err = s.MenuRepository.ResolveMenuUserByID(id)
	if err != nil {
		return
	}
	return
}

func (s *MenuServiceImpl) DeleteMenuUserByID(id uuid.UUID) error {
	menuUser, err := s.MenuRepository.ResolveMenuUserByID(id)
	fmt.Println("menuUser:", menuUser)
	if err != nil {
		return errors.New("menu user dengan ID :" + id.String() + " tidak ditemukan")
	}

	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data menu user dengan ID: " + id.String())
	}
	menuUser.SoftDeleteMenuUser()

	err = s.MenuRepository.UpdateMenuUser(menuUser)
	if err != nil {
		return errors.New("Ada kesalahan dalam menghapus data menu user dengan ID: " + id.String())
	}
	return nil
}

func (s *MenuServiceImpl) SortMenu(sortMenu RequestMenuSortFormat) (menuUser MenuUser, err error) {
	menuUser, err = s.MenuRepository.ResolveMenuUserByID(sortMenu.Id)
	if err != nil {
		return MenuUser{}, errors.New("Menu dengan ID :" + sortMenu.Id.String() + " tidak ditemukan")
	}

	menuUserPost := menuUser

	menuCh := MenuUser{}
	if sortMenu.JenisSort == "UP" {
		menuCh, err = s.MenuRepository.GetMenuUp(menuUser)
		if err != nil {
			return MenuUser{}, errors.New("Menu dengan ID :" + sortMenu.Id.String() + " berada pada posisi teratas")
		}
	}
	if sortMenu.JenisSort == "DOWN" {
		menuCh, err = s.MenuRepository.GetMenuDown(menuUserPost)
		if err != nil {
			return MenuUser{}, errors.New("Menu dengan ID :" + sortMenu.Id.String() + " berada pada posisi terbawah")
		}
	}

	// diturunkan
	menuUserPost.ID = menuUser.ID
	menuUserPost.Urutan = menuCh.Urutan
	menuUserPost.UpdatedAt = null.TimeFrom(time.Now())
	err = s.MenuRepository.UpdateMenuUser(menuUserPost)
	if err != nil {
		log.Error().Msgf("service.UpdateMenu error", err)
	}

	// dinaikan
	menuUserPost.ID = menuCh.ID
	menuUserPost.Urutan = menuUser.Urutan
	menuUserPost.UpdatedAt = null.TimeFrom(time.Now())
	err = s.MenuRepository.UpdateMenuUser(menuUserPost)
	if err != nil {
		log.Error().Msgf("service.UpdateMenu error", err)
	}
	return
}
