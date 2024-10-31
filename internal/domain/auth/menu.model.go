package auth

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

/*
App Status :
1 => SILATURAHMI
2 => CPL
*/
type Menu struct {
	ID           string      `db:"id" json:"id"`
	NamaMenu     string      `db:"nama_menu" json:"namaMenu"`
	LinkMenu     string      `db:"link_menu" json:"linkMenu"`
	Keterangan   null.String `db:"keterangan" json:"keterangan"`
	ClassIcon    null.String `db:"class_icon" json:"classIcon"`
	Status       string      `db:"status" json:"status"`
	CreatedAt    time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt    null.Time   `db:"updated_at" json:"updatedAt"`
	IsPermission null.Bool   `db:"is_permission" json:"isPermission"`
	IsDeleted    bool        `db:"is_deleted" json:"isDeleted"`
	App          int         `db:"app" json:"app"`
}

type RequestMenuFormat struct {
	NamaMenu   string      `db:"nama_menu" json:"namaMenu"`
	LinkMenu   string      `db:"link_menu" json:"linkMenu"`
	Keterangan null.String `db:"keterangan" json:"keterangan"`
	ClassIcon  null.String `db:"class_icon" json:"classIcon"`
	App        int         `db:"app" json:"app"`
}

func (menu *Menu) NewMenuFormat(reqFormat RequestMenuFormat) (newMenu Menu, err error) {
	newID, _ := uuid.NewV4()

	newMenu = Menu{
		ID:         newID.String(),
		NamaMenu:   reqFormat.NamaMenu,
		LinkMenu:   reqFormat.LinkMenu,
		Keterangan: reqFormat.Keterangan,
		ClassIcon:  reqFormat.ClassIcon,
		Status:     "1",
		CreatedAt:  time.Now(),
		App:        reqFormat.App,
	}

	return
}

var ColumnMappMenu = map[string]interface{}{
	"id":           "id",
	"namaMenu":     "nama_menu",
	"linkMenu":     "link_menu",
	"keterangan":   "keterangan",
	"classIcon":    "class_icon",
	"createdAt":    "created_at",
	"updatedAt":    "updated_at",
	"isPermission": "is_permission",
	"isDeleted":    "is_deleted",
	"app":          "app",
}

func (menu *Menu) NewFormatUpdate(reqFormat RequestMenuFormat) (err error) {
	menu.NamaMenu = reqFormat.NamaMenu
	menu.LinkMenu = reqFormat.LinkMenu
	menu.Keterangan = reqFormat.Keterangan
	menu.ClassIcon = reqFormat.ClassIcon
	menu.App = reqFormat.App
	menu.UpdatedAt = null.TimeFrom(time.Now())

	return nil
}
func (menu *Menu) SoftDelete() {
	menu.Status = "0"
	menu.IsDeleted = true
	menu.UpdatedAt = null.TimeFrom(time.Now())
}

type MenuUser struct {
	ID           string      `db:"id" json:"id"`
	IDMenu       string      `db:"id_menu" json:"idMenu"`
	NamaMenu     string      `db:"nama_menu" json:"nama"`
	LinkMenu     string      `db:"link_menu" json:"link"`
	Keterangan   null.String `db:"keterangan" json:"keterangan"`
	ClassIcon    null.String `db:"class_icon" json:"classIcon"`
	Level        int         `db:"level" json:"level"`
	Urutan       int         `db:"urutan" json:"urutan"`
	Posisi       string      `db:"posisi" json:"posisi"`
	IsPermission null.Bool   `db:"is_permission" json:"isPermission"`
	Parent       null.String `db:"parent" json:"parent"`
	IdRole       string      `db:"id_role" json:"idRole"`
	Status       string      `db:"status" json:"status"`
	CreatedAt    time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt    null.Time   `db:"updated_at" json:"updatedAt"`
	IsDeleted    bool        `db:"is_deleted" json:"isDeleted"`
	App          int         `db:"app" json:"app"`
	IdBidang     *string     `db:"id_bidang" json:"idBidang"`
	TenantID     *uuid.UUID  `db:"tenant_id" json:"tenantId"`
	IDBranch     *string     `db:"id_branch" json:"idBranch"`
}

type MenuResponse struct {
	ID           string         `db:"id" json:"id"`
	IDMenu       string         `db:"id_menu" json:"idMenu"`
	NamaMenu     string         `db:"nama_menu" json:"nama"`
	LinkMenu     string         `db:"link_menu" json:"link"`
	Keterangan   null.String    `db:"keterangan" json:"keterangan"`
	ClassIcon    null.String    `db:"class_icon" json:"classIcon"`
	Level        int            `db:"level" json:"level"`
	Posisi       string         `db:"posisi" json:"posisi"`
	Urutan       int            `db:"urutan" json:"urutan"`
	IsPermission null.Bool      `db:"is_permission" json:"isPermission"`
	App          int            `db:"app" json:"app"`
	Children     []MenuResponse `json:"children"`
	KetPosisi    string         `db:"ket_posisi" json:"ketPosisi"`
	LinkParent   *string        `db:"link_parent" json:"linkParent"`
	IDBranch     *string        `db:"id_branch" json:"idBranch"`
	TenantID     *string        `db:"tenant_id" json:"tenantId"`
}

type RequestMenuUserFormat struct {
	IdMenu   []string    `db:"id_menu" json:"idMenu"`
	Posisi   string      `db:"posisi" json:"posisi"`
	Level    int         `db:"level" json:"level"`
	Parent   null.String `db:"parent" json:"parent"`
	IdRole   string      `db:"id_role" json:"idRole"`
	IdBidang *string     `db:"id_bidang" json:"idBidang"`
	IDBranch *string     `db:"id_branch" json:"idBranch"`
	TenantID *uuid.UUID  `db:"tenant_id" json:"tenantId"`
}

func (menuUser *MenuUser) NewMenuUserFormat(reqFormat RequestMenuUserFormat) (newMenuUser []MenuUser, err error) {
	for i := 0; i < len(reqFormat.IdMenu); i++ {
		newID, _ := uuid.NewV4()
		newMenu := MenuUser{
			ID:        newID.String(),
			IDMenu:    reqFormat.IdMenu[i],
			Level:     reqFormat.Level,
			Parent:    reqFormat.Parent,
			IdRole:    reqFormat.IdRole,
			Posisi:    reqFormat.Posisi,
			IdBidang:  reqFormat.IdBidang,
			Status:    "1",
			CreatedAt: time.Now(),
			TenantID:  reqFormat.TenantID,
			IDBranch:  reqFormat.IDBranch,
		}

		newMenuUser = append(newMenuUser, newMenu)
	}

	return
}

func (menuUser *MenuUser) SoftDeleteMenuUser() {
	menuUser.Status = "0"
	menuUser.IsDeleted = true
	menuUser.UpdatedAt = null.TimeFrom(time.Now())
}

type RequestMenuUserFilter struct {
	IdRole        string `json:"idRole"`
	IdMenu        string `json:"idMenu"`
	Posisi        string `json:"posisi"`
	PosisiSubMenu string `json:"posisiSubMenu"`
	Level         string `json:"level"`
	App           string `json:"app"`
	LinkParent    string `json:"linkParent"`
	IdBidang      string `json:"idBidang"`
	IDBranch      string `json:"idBranch"`
	TenantID      string `json:"tenantId"`
}
type RequestMenuSortFormat struct {
	Id        uuid.UUID `db:"id" json:"id"`
	JenisSort string    `json:"jenisSort"`
}
