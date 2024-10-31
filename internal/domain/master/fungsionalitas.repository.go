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
	fungsionalitasQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, nama, level, is_head, jenis_approval, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_fungsionalitas `,
		Insert: `insert into m_fungsionalitas
				(id, nama, level, is_head, jenis_approval, tenant_id, id_branch, created_at, created_by) 
				values 
				(:id, :nama, :level, :is_head, :jenis_approval, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_fungsionalitas set
				id=:id,
				nama=:nama,
				level=:level,
				is_head=:is_head,
				jenis_approval=:jenis_approval,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_fungsionalitas `,
		Count:  `select count (id) from m_fungsionalitas `,
		Exist:  `select count(id)>0 from m_fungsionalitas `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id_fungsionalitas = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type FungsionalitasRepository interface {
	Create(data Fungsionalitas) error
	GetAll() (data []Fungsionalitas, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Fungsionalitas, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data Fungsionalitas) error
	ExistByNama(nama string, id string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type FungsionalitasRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideFungsionalitasRepositoryPostgreSQL(db *infras.PostgresqlConn) *FungsionalitasRepositoryPostgreSQL {
	s := new(FungsionalitasRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *FungsionalitasRepositoryPostgreSQL) Create(data Fungsionalitas) error {
	stmt, err := r.DB.Write.PrepareNamed(fungsionalitasQuery.Insert)
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

func (r *FungsionalitasRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(fungsionalitasQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappFungsionalitas[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchfungsionalitasQuery := searchRoleBuff.String()
	searchfungsionalitasQuery = r.DB.Read.Rebind(fungsionalitasQuery.Select + searchfungsionalitasQuery)
	fmt.Println("query", searchfungsionalitasQuery)
	rows, err := r.DB.Read.Queryx(searchfungsionalitasQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var items Fungsionalitas
		err = rows.StructScan(&items)
		if err != nil {
			return
		}

		data.Items = append(data.Items, items)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *FungsionalitasRepositoryPostgreSQL) GetAll() (data []Fungsionalitas, err error) {
	rows, err := r.DB.Read.Queryx(fungsionalitasQuery.Select + " where coalesce(is_deleted, false) = false")
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Fungsionalitas NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var items Fungsionalitas
		err = rows.StructScan(&items)

		if err != nil {
			return
		}

		data = append(data, items)
	}
	return
}

func (r *FungsionalitasRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data Fungsionalitas, err error) {
	err = r.DB.Read.Get(&data, fungsionalitasQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *FungsionalitasRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(fungsionalitasQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *FungsionalitasRepositoryPostgreSQL) Update(data Fungsionalitas) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *FungsionalitasRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data Fungsionalitas) (err error) {
	stmt, err := tx.PrepareNamed(fungsionalitasQuery.Update + " WHERE id=:id")
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

func (r *FungsionalitasRepositoryPostgreSQL) ExistByNama(nama string, id string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	err := r.DB.Read.Get(&exist, fungsionalitasQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *FungsionalitasRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, fungsionalitasQuery.ExistRelasi, id)

	return
}
