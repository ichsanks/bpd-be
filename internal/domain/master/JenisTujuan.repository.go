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
	jenisTujuanQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, nama, keterangan, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_jenis_tujuan `,
		Insert: `insert into m_jenis_tujuan
				(id, nama, keterangan, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :nama, :keterangan, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_jenis_tujuan set
				id=:id,
				nama=:nama,
				keterangan=:keterangan,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_jenis_tujuan `,
		Count: `select count (id)
				from m_jenis_tujuan `,
		Exist: `select count(id)>0 from m_jenis_tujuan `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type JenisTujuanRepository interface {
	Create(data JenisTujuan) error
	GetAll(req model.StandardRequest) (data []JenisTujuan, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data JenisTujuan, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data JenisTujuan) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type JenisTujuanRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJenisTujuanRepositoryPostgreSQL(db *infras.PostgresqlConn) *JenisTujuanRepositoryPostgreSQL {
	s := new(JenisTujuanRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JenisTujuanRepositoryPostgreSQL) Create(data JenisTujuan) error {
	stmt, err := r.DB.Write.PrepareNamed(jenisTujuanQuery.Insert)
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

func (r *JenisTujuanRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
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

	query := r.DB.Read.Rebind(jenisTujuanQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappJenisTujuan[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchjenisTujuanQuery := searchRoleBuff.String()
	searchjenisTujuanQuery = r.DB.Read.Rebind(jenisTujuanQuery.Select + searchjenisTujuanQuery)
	fmt.Println("query", searchjenisTujuanQuery)
	rows, err := r.DB.Read.Queryx(searchjenisTujuanQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var JenisTujuan JenisTujuan
		err = rows.StructScan(&JenisTujuan)
		if err != nil {
			return
		}

		data.Items = append(data.Items, JenisTujuan)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *JenisTujuanRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []JenisTujuan, err error) {

	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(jenisTujuanQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("JenisTujuan NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JenisTujuan
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *JenisTujuanRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data JenisTujuan, err error) {
	err = r.DB.Read.Get(&data, jenisTujuanQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JenisTujuanRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(jenisTujuanQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JenisTujuanRepositoryPostgreSQL) Update(data JenisTujuan) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *JenisTujuanRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data JenisTujuan) (err error) {
	stmt, err := tx.PrepareNamed(jenisTujuanQuery.Update + " WHERE id=:id")
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

func (r *JenisTujuanRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}
	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, jenisTujuanQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JenisTujuanRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, jenisTujuanQuery.ExistRelasi, id)

	return
}
