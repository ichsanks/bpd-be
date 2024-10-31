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
	dokumenQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, nama, keterangan, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_dokumen `,
		Insert: `insert into m_dokumen
				(id, nama, keterangan, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :nama, :keterangan, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_dokumen set
				id=:id,
				nama=:nama,
				keterangan=:keterangan,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_dokumen `,
		Count: `select count (id)
				from m_dokumen `,
		Exist: `select count(id)>0 from m_dokumen `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type DokumenRepository interface {
	Create(data Dokumen) error
	GetAll(req model.StandardRequest) (data []Dokumen, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Dokumen, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data Dokumen) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type DokumenRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideDokumenRepositoryPostgreSQL(db *infras.PostgresqlConn) *DokumenRepositoryPostgreSQL {
	s := new(DokumenRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *DokumenRepositoryPostgreSQL) Create(data Dokumen) error {
	stmt, err := r.DB.Write.PrepareNamed(dokumenQuery.Insert)
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

func (r *DokumenRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(nama, keterangan) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(dokumenQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappDokumen[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchdokumenQuery := searchRoleBuff.String()
	searchdokumenQuery = r.DB.Read.Rebind(dokumenQuery.Select + searchdokumenQuery)
	fmt.Println("query", searchdokumenQuery)
	rows, err := r.DB.Read.Queryx(searchdokumenQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var Dokumen Dokumen
		err = rows.StructScan(&Dokumen)
		if err != nil {
			return
		}

		data.Items = append(data.Items, Dokumen)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *DokumenRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []Dokumen, err error) {
	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(dokumenQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Dokumen NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList Dokumen
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *DokumenRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data Dokumen, err error) {
	err = r.DB.Read.Get(&data, dokumenQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *DokumenRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(dokumenQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *DokumenRepositoryPostgreSQL) Update(data Dokumen) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *DokumenRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data Dokumen) (err error) {
	stmt, err := tx.PrepareNamed(dokumenQuery.Update + " WHERE id=:id")
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

func (r *DokumenRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, dokumenQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *DokumenRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, dokumenQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *DokumenRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, dokumenQuery.ExistRelasi, id)

	return
}
