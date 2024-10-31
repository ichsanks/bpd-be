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
	jenisPerjalananDinasQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select mjpd.id, mjpd.nama, mjpd.tenant_id, mjpd.id_branch, mjpd.created_at, mjpd.created_by, mjpd.updated_at, mjpd.updated_by, mjpd.is_deleted from m_jenis_perjalanan_dinas mjpd `,
		Insert: `insert into m_jenis_perjalanan_dinas
				(id, nama, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_jenis_perjalanan_dinas set
				id=:id,
				nama=:nama,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_jenis_perjalanan_dinas mjpd `,
		Count: `select count (mjpd.id)
				from m_jenis_perjalanan_dinas mjpd `,
		Exist: `select count(id)>0 from m_jenis_perjalanan_dinas `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from perjalanan_dinas pd 
			where id_jenis_perjalanan_dinas = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type JenisPerjalananDinasRepository interface {
	Create(data JenisPerjalananDinas) error
	GetAll(req model.StandardRequest) (data []JenisPerjalananDinas, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data JenisPerjalananDinas, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data JenisPerjalananDinas) error
	ExistByNama(nama string, idBranch string) (bool, error)
	ExistByNamaID(id uuid.UUID, nama string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type JenisPerjalananDinasRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJenisPerjalananDinasRepositoryPostgreSQL(db *infras.PostgresqlConn) *JenisPerjalananDinasRepositoryPostgreSQL {
	s := new(JenisPerjalananDinasRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JenisPerjalananDinasRepositoryPostgreSQL) Create(data JenisPerjalananDinas) error {
	stmt, err := r.DB.Read.PrepareNamed(jenisPerjalananDinasQuery.Insert)
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

func (r *JenisPerjalananDinasRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE mjpd.is_deleted is false ")

	if req.IdBranch != "" {
		searchRoleBuff.WriteString(" AND mjpd.id_branch = ? ")
		searchParams = append(searchParams, req.IdBranch)
	}

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(mjpd.nama) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(jenisPerjalananDinasQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappJenisPerjalananDinas[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchjenisPerjalananDinasQuery := searchRoleBuff.String()
	searchjenisPerjalananDinasQuery = r.DB.Read.Rebind(jenisPerjalananDinasQuery.Select + searchjenisPerjalananDinasQuery)
	fmt.Println("query", searchjenisPerjalananDinasQuery)
	rows, err := r.DB.Read.Queryx(searchjenisPerjalananDinasQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var JenisPerjalananDinas JenisPerjalananDinas
		err = rows.StructScan(&JenisPerjalananDinas)
		if err != nil {
			return
		}

		data.Items = append(data.Items, JenisPerjalananDinas)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *JenisPerjalananDinasRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []JenisPerjalananDinas, err error) {

	where := " where coalesce(mjpd.is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and mjpd.id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(jenisPerjalananDinasQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("JenisPerjalananDinas NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JenisPerjalananDinas
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *JenisPerjalananDinasRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (JenisPerjalananDinas JenisPerjalananDinas, err error) {
	err = r.DB.Read.Get(&JenisPerjalananDinas, jenisPerjalananDinasQuery.Select+" WHERE mjpd.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JenisPerjalananDinasRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(jenisPerjalananDinasQuery.Delete+" WHERE mjpd.id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JenisPerjalananDinasRepositoryPostgreSQL) Update(data JenisPerjalananDinas) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := txUpdateJenisPerjalananDinas(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func txUpdateJenisPerjalananDinas(tx *sqlx.Tx, data JenisPerjalananDinas) (err error) {
	stmt, err := tx.PrepareNamed(jenisPerjalananDinasQuery.Update + " WHERE id=:id")
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

func (r *JenisPerjalananDinasRepositoryPostgreSQL) ExistByNama(nama string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, jenisPerjalananDinasQuery.Exist+" where upper(nama)=upper($1) and coalesce(is_deleted, false)=false and id_branch=$2 ", nama, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JenisPerjalananDinasRepositoryPostgreSQL) ExistByNamaID(id uuid.UUID, nama string, idBranch string) (bool, error) {
	var exist bool

	err := r.DB.Read.Get(&exist, jenisPerjalananDinasQuery.Exist+" where id <> $1 and upper(nama)=upper($2) and coalesce(is_deleted, false)=false and  id_branch=$3  ", id, nama, idBranch)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JenisPerjalananDinasRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, jenisPerjalananDinasQuery.ExistRelasi, id)

	return
}
