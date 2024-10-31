package auth

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

var (
	roleQuery = struct {
		Select                        string
		Count                         string
		Insert                        string
		Update                        string
		Delete                        string
		Exist                         string
		DeleteMenuByRole              string
		InsertBulkRoleMenu            string
		InsertBulkRoleMenuPlaceHolder string
	}{
		Select: `select id, nama, keterangan, is_view_data_all, is_choose_pegawai, is_choose_terbatas from auth_role `,
		Count:  `select count(id) from auth_role `,
		Insert: `insert into auth_role(id, nama, keterangan, is_view_data_all, is_choose_pegawai, is_choose_terbatas)
							values(:id, :nama, :keterangan, :is_view_data_all, :is_choose_pegawai, :is_choose_terbatas) `,
		Update: `Update auth_role set
						 nama=:nama,
						keterangan=:keterangan,
						is_view_data_all=:is_view_data_all,
						is_choose_pegawai=:is_choose_pegawai,
						is_choose_terbatas=:is_choose_terbatas
						where id=:id `,
		Delete:                        `delete from auth_role `,
		Exist:                         `select count(id)>0 from auth_role `,
		DeleteMenuByRole:              `delete from menu_user `,
		InsertBulkRoleMenu:            `insert into role_menu (role_id, menu) VALUES `,
		InsertBulkRoleMenuPlaceHolder: `(:role_id, :menu) `,
	}
)

type RoleRepository interface {
	GetData(req model.StandardRequest) (data pagination.Response, err error)
	ResolveAll() (roles []Role, err error)
	ExistRoleByID(id string) (bool, error)
	ExistRoleByNama(id string) (bool, error)
	CreateRole(role Role) error
	ResolveRoleByID(id string) (Role, error)
	UpdateRole(role Role) error
	DeleteRoleByID(id string) error
	ResolveRoleMenuByRoleID(roleId string) ([]string, error)
	UpdateRoleMenu(roleId string, menus []string) error
}

type RoleRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideRoleRepositoryPostgreSQL(db *infras.PostgresqlConn) *RoleRepositoryPostgreSQL {
	s := new(RoleRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *RoleRepositoryPostgreSQL) GetData(req model.StandardRequest) (response pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" WHERE ")
		searchRoleBuff.WriteString("  concat (nama, keterangan) like ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := roleQuery.Count + searchRoleBuff.String()
	queryDTO := roleQuery.Select + searchRoleBuff.String()

	query = r.DB.Read.Rebind(query)
	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		response.Items = make([]interface{}, 0)
		return
	}

	offset := (req.PageNumber - 1) * req.PageSize
	queryDTO = r.DB.Read.Rebind(queryDTO + fmt.Sprintf("order by %s %s  limit %d offset %d", req.SortBy, req.SortType, req.PageSize, offset))

	rows, err := r.DB.Read.Queryx(queryDTO, searchParams...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	for rows.Next() {
		var role Role
		err = rows.StructScan(&role)
		if err != nil {
			return
		}

		response.Items = append(response.Items, role)
	}

	response.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *RoleRepositoryPostgreSQL) ResolveAll() (roles []Role, err error) {
	err = r.DB.Read.Select(&roles, roleQuery.Select+" order by id")
	return
}

func (r *RoleRepositoryPostgreSQL) ExistRoleByID(id string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, roleQuery.Exist+" where id = $1", id)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *RoleRepositoryPostgreSQL) ExistRoleByNama(nama string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, roleQuery.Exist+" where nama = $1", nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

// CreateRole adalah method untuk menambahkan Role Baru
func (r *RoleRepositoryPostgreSQL) CreateRole(role Role) error {
	stmt, err := r.DB.Read.PrepareNamed(roleQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(role)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

// ResolveByID adalah method yang digunakan untuk mendapatkan data Role berdasarkan ID
func (r *RoleRepositoryPostgreSQL) ResolveRoleByID(id string) (Role, error) {
	var role Role
	err := r.DB.Read.Get(&role, roleQuery.Select+" where id=$1", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return Role{}, err
	}
	return role, nil
}

// UpdateRole adalah method untuk mengubah Role yang sudah ada
func (r *RoleRepositoryPostgreSQL) UpdateRole(role Role) error {
	stmt, err := r.DB.Write.PrepareNamed(roleQuery.Update)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(role)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

// DeleteRoleByID adalah method untuk menghapus Role berdasarkan ID yang dikirimkan
func (r *RoleRepositoryPostgreSQL) DeleteRoleByID(id string) (err error) {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txDeleteRoleMenuByRoleID(id); err != nil {
			e <- err
			return
		}
		stmt, err := r.DB.Read.PrepareNamed(roleQuery.Delete + " where id=:id ")
		if err != nil {
			logger.ErrorWithStack(err)
			e <- err
			return
		}
		var params = map[string]interface{}{
			"id": id,
		}
		defer stmt.Close()
		res, err := stmt.Exec(params)
		count, err := res.RowsAffected()
		if err != nil {
			logger.ErrorWithStack(err)
			e <- err
			return
		} else {
			if count > 0 {
				if err := r.txDeleteRoleMenuByRoleID(id); err != nil {
					e <- err
					return
				}
				e <- nil
				return
			}
		}

		e <- nil
	})
}

// RoleMenus adalan method untuk menampilkan menu berdasarkanRoleID
func (r *RoleRepositoryPostgreSQL) ResolveRoleMenuByRoleID(roleId string) (menus []string, err error) {
	err = r.DB.Read.Select(&menus, "select menu from role_menu where role_id =$1", roleId)
	return
}

// UpdateRoleMenu composes a Bulk insert menu query given a slice of RoleMenu
func (r *RoleRepositoryPostgreSQL) UpdateRoleMenu(roleId string, menus []string) (err error) {
	err = r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdateRoleMenu(tx, roleId, menus); err != nil {
			e <- err
			return
		}
		e <- nil
	})
	if err != nil {
		return
	}
	return
}

func (r *RoleRepositoryPostgreSQL) txUpdateRoleMenu(tx *sqlx.Tx, roleId string, menus []string) (err error) {
	if len(menus) == 0 {
		return errors.New("Data menu kosong")
	}
	if err = r.txDeleteRoleMenuByRoleID(roleId); err != nil {
		return
	}

	query, args, err := r.composeBulkInsertRoleMenuQuery(roleId, menus)
	if err != nil {
		return
	}
	// Secara default sqlx akan binding Named parameter dengan ?
	// Contoh: insert into role_menu (role_id, menu) VALUES  (?, ?) ,(?, ?) ,(?, ?) ,(?, ?)
	// Maka kemudian akan muncul error gini:
	// syntax error at or near "," at character 54
	// Oleh karena itu kita harus rebind Named sesuai dengan database yang sudah kita setting
	query = tx.Rebind(query)
	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}

	defer stmt.Close()
	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return nil
}

func (r *RoleRepositoryPostgreSQL) composeBulkInsertRoleMenuQuery(roleId string, menus []string) (query string, params []interface{}, err error) {
	values := []string{}
	for _, menu := range menus {
		param := map[string]interface{}{
			"role_id": roleId,
			"menu":    menu,
		}
		q, args, err := sqlx.Named(roleQuery.InsertBulkRoleMenuPlaceHolder, param)
		if err != nil {
			return query, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	query = fmt.Sprintf("%v %v", roleQuery.InsertBulkRoleMenu, strings.Join(values, ","))

	log.Info().Msg(query)
	return
}

func (r *RoleRepositoryPostgreSQL) txDeleteRoleMenuByRoleID(roleId string) (err error) {
	stmt, err := r.DB.Read.PrepareNamed(roleQuery.DeleteMenuByRole + " where id_role=:roleId")
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	arg := map[string]interface{}{
		"roleId": roleId,
	}
	defer stmt.Close()
	_, err = stmt.Exec(arg)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}
