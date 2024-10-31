package master

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/pagination"
)

var (
	branchQuery = struct {
		Select    string
		SelectDTO string
		Insert    string
		Update    string
		Delete    string
		Exist     string
		Count     string
	}{
		Select: `SELECT id, kode, nama, email, address, city, contact, phone, website, image, tenant_id, created_by, created_at, updated_by, updated_at, is_deleted, color, is_dark
		FROM m_branch`,
		SelectDTO: `SELECT 
						id, kode, nama, email, address, city, contact, phone, website, image, tenant_id, 
						created_by, created_at, updated_by, updated_at, is_deleted, color, is_dark
						FROM m_branch`,
		Insert: `INSERT INTO 
						m_branch 
							(id, kode, nama, email, address, city, contact, phone, website, image, 
							tenant_id, created_by, created_at,  color, is_dark) 
		                values
							(:id, :kode, :nama, :email, :address, :city, :contact, :phone, :website, :image, 
							:tenant_id, :created_by, :created_at,  :color, :is_dark) `,
		Update: `UPDATE m_branch SET 
				id=:id, 
				kode=:kode, 
				nama=:nama, 
				email=:email, 
				address=:address, 
				city=:city, 
				contact=:contact, 
				phone=:phone,
				website=:website,
				image=:image,
				color=:color,
				is_dark=:is_dark,
				tenant_id=:tenant_id,
				updated_at=:updated_at,
				updated_by=:updated_by, 
				is_deleted=:is_deleted`,
		Delete: `delete from m_branch `,
		Exist:  `select count(id)>0 from m_branch `,
		Count:  `select count(id) from m_branch `,
	}
)

type BranchRepository interface {
	Create(data Branch) error
	Update(data Branch) error
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	GetAllData() (data []Branch, err error)
	ResolveByID(id uuid.UUID) (data Branch, err error)
	ResolveByIDDTO(id uuid.UUID) (data BranchDTO, err error)
	ExistData(kode string, nama string, id string) (bool, error)
	ExistKode(kode string, id string) (bool, error)
	ExistName(name string, id string) (bool, error)
}

type BranchRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideBranchRepositoryPostgreSQL(db *infras.PostgresqlConn) *BranchRepositoryPostgreSQL {
	s := new(BranchRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *BranchRepositoryPostgreSQL) Create(data Branch) error {
	stmt, err := r.DB.Read.PrepareNamed(branchQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}

func (r *BranchRepositoryPostgreSQL) Update(data Branch) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateBranch(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateBranch(tx *sqlx.Tx, data Branch) (err error) {
	stmt, err := tx.PrepareNamed(branchQuery.Update + " WHERE id=:id")
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(data)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (r *BranchRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false)=false ")

	if req.TenantID != "" {
		searchRoleBuff.WriteString(" AND tenant_id = ? ")
		searchParams = append(searchParams, req.TenantID)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(nama, email, address, city, contact, phone, website) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(branchQuery.Count + searchRoleBuff.String())

	var totalData int
	err = r.DB.Read.QueryRow(query, searchParams...).Scan(&totalData)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if totalData < 1 {
		data.Items = make([]interface{}, 0)
		return
	}

	searchRoleBuff.WriteString("order by " + ColumnMappBranch[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchbranchQuery := searchRoleBuff.String()
	searchbranchQuery = r.DB.Read.Rebind(branchQuery.Select + searchbranchQuery)
	rows, err := r.DB.Read.Queryx(searchbranchQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var items Branch
		err = rows.StructScan(&items)
		if err != nil {
			return
		}

		data.Items = append(data.Items, items)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)

	return
}

func (r *BranchRepositoryPostgreSQL) GetAllData() (data []Branch, err error) {
	criteria := ` where coalesce(is_deleted,false)=false `
	// if tenantID != "" {
	// 	criteria += fmt.Sprintf(" and tenant_id='%s' ", tenantID)
	// }
	rows, err := r.DB.Read.Queryx(branchQuery.Select + criteria)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Branch")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items Branch
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}
	return
}

func (r *BranchRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data Branch, err error) {
	err = r.DB.Read.Get(&data, branchQuery.Select+" WHERE id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *BranchRepositoryPostgreSQL) ResolveByIDDTO(id uuid.UUID) (data BranchDTO, err error) {
	err = r.DB.Read.Get(&data, branchQuery.SelectDTO+" WHERE id=$1 ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *BranchRepositoryPostgreSQL) ExistData(kode string, nama string, id string) (bool, error) {
	var exist bool
	query := ` where coalesce(is_deleted, false)=false and kode= $1  and nama = $2 `
	if id != "" {
		query += fmt.Sprintf(" and id != '%s'", id)
	}

	err := r.DB.Read.Get(&exist, branchQuery.Exist+query, kode, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *BranchRepositoryPostgreSQL) ExistKode(kode string, id string) (bool, error) {
	var exist bool
	query := ` where coalesce(is_deleted, false)=false and upper(kode)=upper($1) `
	if id != "" {
		query += fmt.Sprintf(" and id != '%s'", id)
	}

	err := r.DB.Read.Get(&exist, branchQuery.Exist+query, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *BranchRepositoryPostgreSQL) ExistName(name string, id string) (bool, error) {
	var exist bool
	query := ` where coalesce(is_deleted, false)=false and upper(nama)=upper($1) `
	if id != "" {
		query += fmt.Sprintf(" and id != '%s'", id)
	}

	err := r.DB.Read.Get(&exist, branchQuery.Exist+query, name)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}
