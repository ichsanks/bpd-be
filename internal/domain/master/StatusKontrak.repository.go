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
	StatusKontrakQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id,  nama, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_status_kontrak `,
		Insert: `insert into m_status_kontrak
				(id,  nama, tenant_id, id_branch, created_at, created_by)
				values
				(:id,  :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_status_kontrak set
				id=:id,
				nama=:nama,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_status_kontrak `,
		Count: `select count (id)
				from m_status_kontrak `,
		Exist: `select count(id)>0 from m_status_kontrak `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type StatusKontrakRepository interface {
	Create(data StatusKontrak) error
	GetAll(req model.StandardRequest) (data []StatusKontrak, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data StatusKontrak, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data StatusKontrak) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type StatusKontrakRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideStatusKontrakRepositoryPostgreSQL(db *infras.PostgresqlConn) *StatusKontrakRepositoryPostgreSQL {
	s := new(StatusKontrakRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *StatusKontrakRepositoryPostgreSQL) Create(data StatusKontrak) error {
	stmt, err := r.DB.Write.PrepareNamed(StatusKontrakQuery.Insert)
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

func (r *StatusKontrakRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(kode, nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(StatusKontrakQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappStatusKontrak[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchStatusKontrakQuery := searchRoleBuff.String()
	searchStatusKontrakQuery = r.DB.Read.Rebind(StatusKontrakQuery.Select + searchStatusKontrakQuery)
	fmt.Println("query", searchStatusKontrakQuery)
	rows, err := r.DB.Read.Queryx(searchStatusKontrakQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var StatusKontrak StatusKontrak
		err = rows.StructScan(&StatusKontrak)
		if err != nil {
			return
		}

		data.Items = append(data.Items, StatusKontrak)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *StatusKontrakRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []StatusKontrak, err error) {
	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(StatusKontrakQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("StatusKontrak NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList StatusKontrak
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *StatusKontrakRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data StatusKontrak, err error) {
	err = r.DB.Read.Get(&data, StatusKontrakQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *StatusKontrakRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(StatusKontrakQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *StatusKontrakRepositoryPostgreSQL) Update(data StatusKontrak) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *StatusKontrakRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data StatusKontrak) (err error) {
	stmt, err := tx.PrepareNamed(StatusKontrakQuery.Update + " WHERE id=:id")
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

func (r *StatusKontrakRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, StatusKontrakQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *StatusKontrakRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id = '%s' ", id)
	}

	err := r.DB.Read.Get(&exist, StatusKontrakQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *StatusKontrakRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, StatusKontrakQuery.ExistRelasi, id)

	return
}