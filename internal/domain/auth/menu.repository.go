package auth

import (
	"bytes"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

var (
	menuQuery = struct {
		Select, Insert, Update, Count string
	}{
		Select: `SELECT id, nama_menu, link_menu, keterangan, class_icon, status, created_at, updated_at, is_permission,coalesce(is_deleted, false) as is_deleted, app from menu `,
		Insert: `Insert into menu(id, nama_menu, link_menu, keterangan, class_icon, status, created_at, app) values (:id, :nama_menu, :link_menu, :keterangan, :class_icon, :status, :created_at, :app)`,
		Update: `Update menu set id=:id, nama_menu=:nama_menu, link_menu=:link_menu, keterangan=:keterangan, class_icon=:class_icon, updated_at=:updated_at, is_deleted=:is_deleted`,
		Count:  `Select count(id) from menu`,
	}

	menuUserQuery = struct {
		Select,
		SelectDTO,
		SelectDTO2,
		Insert,
		Update string
	}{
		Select: `SELECT id, id_menu, posisi, level, urutan, status, created_at, updated_at, parent, id_role from menu_user `,
		SelectDTO: ` SELECT mu.id as id, m.id as id_menu, m.nama_menu, m.link_menu, m.keterangan, m.class_icon, m.is_permission, mu.level, mu.urutan, mu.posisi, m.app, m2.link_menu link_parent,
					(CASE WHEN mu.posisi='1' THEN 'SIDEBAR'
						WHEN mu.posisi='2' THEN 'CARD MENU'
						WHEN mu.posisi='3' THEN 'DETAIL MENU'
						ELSE '' END) AS ket_posisi
					FROM menu_user mu 
					JOIN menu m on mu.id_menu = m.id 
					LEFT JOIN menu m2 on mu.parent = m2.id
					`,
		SelectDTO2: `
			SELECT x.id, x.id_menu, x.nama_menu, x.link_menu, x.keterangan, x.class_icon, x.is_permission, x.level, x.urutan, x.posisi, x.id_branch, x.tenant_id,
				x.app, x.link_parent, (CASE WHEN x.posisi='1' THEN 'SIDEBAR'
				WHEN x.posisi='2' THEN 'CARD MENU'
				WHEN x.posisi='3' THEN 'DETAIL MENU'
				ELSE '' END) AS ket_posisi 
			FROM (
				SELECT mu.id as id, m.id as id_menu, m.nama_menu, m.link_menu, m.keterangan, m.class_icon, m.is_permission, mu.level, mu.urutan, mu.posisi, 
				m.app, m2.link_menu link_parent, mu.id_bidang, mu.parent, mu.id_branch, mu.tenant_id
				FROM menu_user mu 
				JOIN menu m on mu.id_menu = m.id 
				LEFT JOIN menu m2 on mu.parent = m2.id
				WHERE m.status = '1' 
				and coalesce(mu.is_deleted, false)= false
				and mu.id_role = $1 
				UNION
				SELECT mu.id as id, m.id as id_menu, m.nama_menu, m.link_menu, m.keterangan, m.class_icon, m.is_permission, mu.level, mu.urutan, mu.posisi, 
				m.app, m2.link_menu link_parent, mu.id_bidang, mu.parent, mu.id_branch, mu.tenant_id
				FROM menu_user mu 
				JOIN menu m on mu.id_menu = m.id 
				LEFT JOIN menu m2 on mu.parent = m2.id
				WHERE m.status = '1' 
				and coalesce(mu.is_deleted, false)= false
				and mu.id_role = $1 
				)x
				`,
		// and mu.id_bidang is null
		// and mu.id_bidang = $2
		// and mu.id_branch is null
		// and mu.id_branch = $3
		Insert: `Insert into menu_user(id, id_menu, posisi, level, urutan, status, created_at, parent, id_role, id_bidang, tenant_id, id_branch) values
					 (:id, :id_menu, :posisi, :level, :urutan, :status, :created_at, :parent, :id_role, :id_bidang, :tenant_id, :id_branch)`,
		Update: `Update menu_user set id=:id, is_deleted=:is_deleted, urutan=:urutan, updated_at=:updated_at`,
	}
)

type MenuRepository interface {
	GetAllMenu(app string) (dataMenu []Menu, err error)
	ResolveAll(req model.StandardRequestMenu) (dataMenu pagination.Response, err error)
	ResolveMenuByRoleID(req RequestMenuUserFilter) (data []MenuResponse, err error)
	ResolveMenuByParentID(req RequestMenuUserFilter) (data []MenuResponse, err error)
	CreateMenu(menu Menu) error
	UpdateMenu(menu Menu) error
	ResolveMenuByID(id uuid.UUID) (menu Menu, err error)
	CreateMenuUser(menuUser MenuUser) error
	ResolveMenuUserByID(id uuid.UUID) (menuUser MenuUser, err error)
	UpdateMenuUser(menuUser MenuUser) error
	GetMenuUp(data MenuUser) (menu MenuUser, err error)
	GetMenuDown(data MenuUser) (menu MenuUser, err error)
}

type MenuRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideMenuRepositoryPostgreSQL(db *infras.PostgresqlConn) *MenuRepositoryPostgreSQL {
	s := new(MenuRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuByRoleID(req RequestMenuUserFilter) (data []MenuResponse, err error) {
	criteria := ` where x.id is not null `

	// if req.IdBidang != "" {
	// 	criteria += fmt.Sprintf(` and x.id_bidang = '%s' `, req.IdBidang)
	// }

	if req.App != "" {
		criteria += fmt.Sprintf(` and x.app = '%s' `, req.App)
	}

	if req.Posisi != "" {
		criteria += fmt.Sprintf(` and x.posisi = '%s' `, req.Posisi)
	}

	if req.Level != "" {
		criteria += fmt.Sprintf(` and x.level = '%s' `, req.Level)
	}

	if req.LinkParent != "" {
		criteria += fmt.Sprintf(` and x.link_menu = '%s' `, req.LinkParent)
	}

	if req.IDBranch != "" {
		criteria += fmt.Sprintf(` and x.id_branch = '%s' `, req.IDBranch)
	}

	if req.IdRole != "HA01" {
		criteria += fmt.Sprintf(` and x.tenant_id = '%s' `, req.TenantID)
	}

	criteria += ` order by x.urutan asc`

	err = r.DB.Read.Select(&data, menuUserQuery.SelectDTO2+criteria, req.IdRole)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuByParentID(req RequestMenuUserFilter) (data []MenuResponse, err error) {
	criteria := ` WHERE x.id is not null and x.level = 2 and x.parent = $2 `

	// if req.IdBidang != "" {
	// 	criteria += fmt.Sprintf(` and x.id_bidang = '%s' `, req.IdBidang)
	// }

	if req.App != "" {
		criteria += fmt.Sprintf(` and x.app = '%s' `, req.App)
	}

	if req.PosisiSubMenu != "" {
		criteria += fmt.Sprintf(` and x.posisi = '%s' `, req.PosisiSubMenu)
	}

	if req.IDBranch != "" {
		criteria += fmt.Sprintf(` AND x.id_branch = '%s' `, req.IDBranch)
	}

	if req.IdRole != "HA01" {
		criteria += fmt.Sprintf(` and x.tenant_id = '%s' `, req.TenantID)
	}

	criteria += ` order by x.urutan asc`
	err = r.DB.Read.Select(&data, menuUserQuery.SelectDTO2+criteria, req.IdRole, req.IdMenu)

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *MenuRepositoryPostgreSQL) GetAllMenu(app string) (dataMenu []Menu, err error) {
	criteria := ` WHERE COALESCE(is_deleted, false) <> true `
	if app != "" {
		criteria += fmt.Sprintf(` and app = '%s' `, app)
	}

	err = r.DB.Read.Select(&dataMenu, menuQuery.Select+criteria)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *MenuRepositoryPostgreSQL) ResolveAll(req model.StandardRequestMenu) (dataMenu pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" where coalesce(is_deleted, false) = ? ")
	searchParams = append(searchParams, false)

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND concat(nama_menu, link_menu) like ? ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	if req.App != "" {
		searchRoleBuff.WriteString(" AND app = ? ")
		searchParams = append(searchParams, req.App)
	}

	query := r.DB.Read.Rebind(menuQuery.Count + searchRoleBuff.String())

	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		dataMenu.Items = make([]interface{}, 0)
		return
	}

	searchRoleBuff.WriteString("order by " + ColumnMappMenu[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchMenuQuery := searchRoleBuff.String()
	searchMenuQuery = r.DB.Read.Rebind(menuQuery.Select + searchMenuQuery)
	rows, err := r.DB.Read.Queryx(searchMenuQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var menu Menu
		err = rows.StructScan(&menu)
		if err != nil {
			return
		}

		dataMenu.Items = append(dataMenu.Items, menu)
	}

	dataMenu.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *MenuRepositoryPostgreSQL) CreateMenu(menu Menu) error {
	stmt, err := r.DB.Read.PrepareNamed(menuQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(menu)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

// ResolveMenuByID digunakan untuk mendapatkan data Menu berdasarkan ID
func (r *MenuRepositoryPostgreSQL) ResolveMenuByID(id uuid.UUID) (menu Menu, err error) {
	err = r.DB.Read.Get(&menu, menuQuery.Select+" WHERE id=$1 AND coalesce(is_deleted, false) = false ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) UpdateMenu(menu Menu) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateMenu(tx, menu); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateMenu(tx *sqlx.Tx, menu Menu) (err error) {
	stmt, err := tx.PrepareNamed(menuQuery.Update + " WHERE id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(menu)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *MenuRepositoryPostgreSQL) ResolveMenuUserByID(id uuid.UUID) (menuUser MenuUser, err error) {
	err = r.DB.Read.Get(&menuUser, menuUserQuery.Select+" WHERE id=$1 AND coalesce(is_deleted, false) = false ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) CreateMenuUser(menuUser MenuUser) error {
	var urut []MenuUser
	if menuUser.Parent.String != "" {
		err := r.DB.Read.Select(&urut, menuUserQuery.Select+` where posisi::varchar = $1 and level = $2 and parent = $3 and id_role::varchar = $4 and coalesce(is_deleted, false) = false and status = '1'
							order by urutan desc limit 1`, menuUser.Posisi, menuUser.Level, menuUser.Parent, menuUser.IdRole)
		if err != nil {
			logger.ErrorWithStack(err)
			return err
		}
	} else {
		err := r.DB.Read.Select(&urut, menuUserQuery.Select+` where posisi::varchar = $1 and level = $2 and id_role::varchar = $3 and coalesce(is_deleted, false) = false and status = '1'
							order by urutan desc limit 1`, menuUser.Posisi, menuUser.Level, menuUser.IdRole)
		if err != nil {
			logger.ErrorWithStack(err)
			return err
		}
	}

	if len(urut) > 0 {
		menuUser.Urutan = urut[0].Urutan + 1
	} else {
		menuUser.Urutan = 1
	}

	stmt, err := r.DB.Read.PrepareNamed(menuUserQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(menuUser)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

func (r *MenuRepositoryPostgreSQL) UpdateMenuUser(menuUser MenuUser) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateMenuUser(tx, menuUser); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateMenuUser(tx *sqlx.Tx, menuUser MenuUser) (err error) {
	stmt, err := tx.PrepareNamed(menuUserQuery.Update + " WHERE id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(menuUser)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *MenuRepositoryPostgreSQL) GetMenuUp(data MenuUser) (menu MenuUser, err error) {
	err = r.DB.Read.Get(&menu, menuUserQuery.Select+" where coalesce(is_deleted, false) <> true and posisi = $1 and level = $2 and coalesce(parent, '') = $3 and urutan <= $4 and id_role = $5 and id <> $6 order by urutan desc limit 1 ", data.Posisi, data.Level, data.Parent, data.Urutan, data.IdRole, data.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *MenuRepositoryPostgreSQL) GetMenuDown(data MenuUser) (menu MenuUser, err error) {
	err = r.DB.Read.Get(&menu, menuUserQuery.Select+" where coalesce(is_deleted, false) <> true and posisi = $1 and level = $2 and coalesce(parent, '') = $3 and urutan >= $4 and id_role = $5 and id <> $6 order by urutan asc limit 1 ", data.Posisi, data.Level, data.Parent, data.Urutan, data.IdRole, data.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}
