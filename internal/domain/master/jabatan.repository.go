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
	jabatanQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select mj.id, mj.nama, mj.tenant_id, mj.id_branch, mj.created_at, mj.created_by, mj.updated_at, mj.updated_by, mj.is_deleted from m_jabatan mj`,
		Insert: `insert into m_jabatan
				(id, nama, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_jabatan set
				id=:id,
				nama=:nama,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_jabatan mj `,
		Count: `select count (mj.id)
				from m_jabatan mj `,
		Exist: `select count(id)>0 from m_jabatan `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id_jabatan = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type JabatanRepository interface {
	Create(data Jabatan) error
	GetAll(req model.StandardRequest) (data []Jabatan, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Jabatan, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data Jabatan) error
	ExistByNama(nama string, IdBranch string) (bool, error)
	ExistByNamaID(id uuid.UUID, nama string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type JabatanRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJabatanRepositoryPostgreSQL(db *infras.PostgresqlConn) *JabatanRepositoryPostgreSQL {
	s := new(JabatanRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JabatanRepositoryPostgreSQL) Create(data Jabatan) error {
	stmt, err := r.DB.Read.PrepareNamed(jabatanQuery.Insert)
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

func (r *JabatanRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE mj.is_deleted is false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND mj.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(mj.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(jabatanQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappJabatan[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchJabatanQuery := searchRoleBuff.String()
	searchJabatanQuery = r.DB.Read.Rebind(jabatanQuery.Select + searchJabatanQuery)
	fmt.Println("query", searchJabatanQuery)
	rows, err := r.DB.Read.Queryx(searchJabatanQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var Jabatan Jabatan
		err = rows.StructScan(&Jabatan)
		if err != nil {
			return
		}

		data.Items = append(data.Items, Jabatan)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *JabatanRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []Jabatan, err error) {
	where := " where coalesce(mj.is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and mj.id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(jabatanQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Jabatan NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList Jabatan
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *JabatanRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (Jabatan Jabatan, err error) {
	err = r.DB.Read.Get(&Jabatan, jabatanQuery.Select+" WHERE mj.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JabatanRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(jabatanQuery.Delete+" WHERE mj.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JabatanRepositoryPostgreSQL) Update(data Jabatan) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateJabatan(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateJabatan(tx *sqlx.Tx, data Jabatan) (err error) {
	stmt, err := tx.PrepareNamed(jabatanQuery.Update + " WHERE id=:id")
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

func (r *JabatanRepositoryPostgreSQL) ExistByNama(nama string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, jabatanQuery.Exist+" where upper(nama)=upper($1) and coalesce(is_deleted, false)=false and id_branch=$2 ", nama, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JabatanRepositoryPostgreSQL) ExistByNamaID(id uuid.UUID, nama string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, jabatanQuery.Exist+" where id <> $1 and upper(nama)=upper($2) and coalesce(is_deleted, false)=false and id_branch=$3 ", id, nama, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JabatanRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, jabatanQuery.ExistRelasi, id)

	return
}
