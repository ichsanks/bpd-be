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
	kategoriBiayaQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, nama, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_kategori_biaya `,
		Insert: `insert into m_kategori_biaya
				(id,  nama, tenant_id, id_branch, created_at, created_by)
				values
				(:id,  :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_kategori_biaya set
				id=:id,
				nama=:nama,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_kategori_biaya `,
		Count: `select count (id)
				from m_kategori_biaya `,
		Exist: `select count(id)>0 from m_kategori_biaya `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type KategoriBiayaRepository interface {
	Create(data KategoriBiaya) error
	GetAll(req model.StandardRequest) (data []KategoriBiaya, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data KategoriBiaya, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data KategoriBiaya) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type KategoriBiayaRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideKategoriBiayaRepositoryPostgreSQL(db *infras.PostgresqlConn) *KategoriBiayaRepositoryPostgreSQL {
	s := new(KategoriBiayaRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *KategoriBiayaRepositoryPostgreSQL) Create(data KategoriBiaya) error {
	stmt, err := r.DB.Write.PrepareNamed(kategoriBiayaQuery.Insert)
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

func (r *KategoriBiayaRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(kategoriBiayaQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappKategoriBiaya[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchkategoriBiayaQuery := searchRoleBuff.String()
	searchkategoriBiayaQuery = r.DB.Read.Rebind(kategoriBiayaQuery.Select + searchkategoriBiayaQuery)
	fmt.Println("query", searchkategoriBiayaQuery)
	rows, err := r.DB.Read.Queryx(searchkategoriBiayaQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var KategoriBiaya KategoriBiaya
		err = rows.StructScan(&KategoriBiaya)
		if err != nil {
			return
		}

		data.Items = append(data.Items, KategoriBiaya)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *KategoriBiayaRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []KategoriBiaya, err error) {
	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(kategoriBiayaQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("KategoriBiaya NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList KategoriBiaya
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *KategoriBiayaRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data KategoriBiaya, err error) {
	err = r.DB.Read.Get(&data, kategoriBiayaQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *KategoriBiayaRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(kategoriBiayaQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *KategoriBiayaRepositoryPostgreSQL) Update(data KategoriBiaya) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *KategoriBiayaRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data KategoriBiaya) (err error) {
	stmt, err := tx.PrepareNamed(kategoriBiayaQuery.Update + " WHERE id=:id")
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

func (r *KategoriBiayaRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, kategoriBiayaQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *KategoriBiayaRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, kategoriBiayaQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *KategoriBiayaRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, kategoriBiayaQuery.ExistRelasi, id)

	return
}
