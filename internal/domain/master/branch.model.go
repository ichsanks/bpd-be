package master

import (
	"time"

	"github.com/gofrs/uuid"
)

type Branch struct {
	Id        uuid.UUID  `db:"id" json:"id"`
	Kode      string     `db:"kode" json:"kode"`
	Nama      string     `db:"nama" json:"nama"`
	Email     *string    `db:"email" json:"email"`
	Address   *string    `db:"address" json:"address"`
	City      *string    `db:"city" json:"city"`
	Contact   *string    `db:"contact" json:"contact"`
	Phone     *string    `db:"phone" json:"phone"`
	Website   *string    `db:"website" json:"website"`
	Image     *string    `db:"image" json:"image"`
	Color     *string    `db:"color" json:"color"`
	IsDark    bool       `db:"is_dark" json:"isDark"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy" validate:"required"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt" validate:"required"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
	TenantID  uuid.UUID  `db:"tenant_id" json:"tenantId"`
}
type BranchDTO struct {
	Id        uuid.UUID  `db:"id" json:"id"`
	Kode      string     `db:"kode" json:"kode"`
	Nama      string     `db:"nama" json:"nama"`
	Email     *string    `db:"email" json:"email"`
	Address   *string    `db:"address" json:"address"`
	City      *string    `db:"city" json:"city"`
	Contact   *string    `db:"contact" json:"contact"`
	Phone     *string    `db:"phone" json:"phone"`
	Website   *string    `db:"website" json:"website"`
	Image     *string    `db:"image" json:"image"`
	Color     *string    `db:"color" json:"color"`
	IsDark    bool       `db:"is_dark" json:"isDark"`
	CreatedBy *uuid.UUID `db:"created_by" json:"createdBy" validate:"required"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt" validate:"required"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updatedBy,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
	TenantID  uuid.UUID  `db:"tenant_id" json:"tenantId"`
}
type RequestBranchFormat struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Kode     string    `db:"kode" json:"kode"`
	Nama     string    `db:"nama" json:"nama"`
	Email    *string   `db:"email" json:"email"`
	Address  *string   `db:"address" json:"address"`
	City     *string   `db:"city" json:"city"`
	Contact  *string   `db:"contact" json:"contact"`
	Phone    *string   `db:"phone" json:"phone"`
	Website  *string   `db:"website" json:"website"`
	Image    *string   `db:"image" json:"image"`
	Color    *string   `db:"color" json:"color"`
	IsDark   bool      `db:"is_dark" json:"isDark"`
	TenantID uuid.UUID `db:"tenant_id" json:"-"`
}

func (b *Branch) BranchFormat(reqFormat RequestBranchFormat, userID uuid.UUID) (newBranch Branch, err error) {
	newID, _ := uuid.NewV4()
	now := time.Now()

	if reqFormat.Id == uuid.Nil {
		newBranch = Branch{
			Id:        newID,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			Email:     reqFormat.Email,
			Address:   reqFormat.Address,
			City:      reqFormat.City,
			Contact:   reqFormat.Contact,
			Phone:     reqFormat.Phone,
			Website:   reqFormat.Website,
			Image:     reqFormat.Image,
			Color:     reqFormat.Color,
			IsDark:    reqFormat.IsDark,
			CreatedBy: &userID,
			CreatedAt: now,
			IsDeleted: false,
			TenantID:  reqFormat.TenantID,
		}
	} else {
		newBranch = Branch{
			Id:        reqFormat.Id,
			Kode:      reqFormat.Kode,
			Nama:      reqFormat.Nama,
			Email:     reqFormat.Email,
			Address:   reqFormat.Address,
			City:      reqFormat.City,
			Contact:   reqFormat.Contact,
			Phone:     reqFormat.Phone,
			Website:   reqFormat.Website,
			Image:     reqFormat.Image,
			Color:     reqFormat.Color,
			IsDark:    reqFormat.IsDark,
			UpdatedBy: &userID,
			UpdatedAt: &now,
			TenantID:  reqFormat.TenantID,
		}
	}
	return
}

func (b *Branch) SoftDelete(userId uuid.UUID) {
	now := time.Now()
	b.IsDeleted = true
	b.UpdatedBy = &userId
	b.UpdatedAt = &now
}

var ColumnMappBranch = map[string]interface{}{
	"id":        "id",
	"kode":      "kode",
	"nama":      "nama",
	"email":     "email",
	"address":   "address",
	"city":      "city",
	"contact":   "contact",
	"phone":     "phone",
	"website":   "website",
	"createdBy": "created_by",
	"createdAt": "created_at",
	"updatedBy": "updated_by",
	"updatedAt": "updated_at",
	"isDeleted": "is_deleted",
}
