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
	golonganQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, kode, nama, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_golongan `,
		Insert: `insert into m_golongan
				(id, kode, nama, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :kode, :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_golongan set
				id=:id,
				kode=:kode,
				nama=:nama,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_golongan `,
		Count: `select count (id)
				from m_golongan `,
		Exist: `select count(id)>0 from m_golongan `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id_golongan = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type GolonganRepository interface {
	Create(data Golongan) error
	GetAll(req model.StandardRequest) (data []Golongan, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data Golongan, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data Golongan) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type GolonganRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideGolonganRepositoryPostgreSQL(db *infras.PostgresqlConn) *GolonganRepositoryPostgreSQL {
	s := new(GolonganRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *GolonganRepositoryPostgreSQL) Create(data Golongan) error {
	stmt, err := r.DB.Write.PrepareNamed(golonganQuery.Insert)
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

func (r *GolonganRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
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

	query := r.DB.Read.Rebind(golonganQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappGolongan[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchgolonganQuery := searchRoleBuff.String()
	searchgolonganQuery = r.DB.Read.Rebind(golonganQuery.Select + searchgolonganQuery)
	fmt.Println("query", searchgolonganQuery)
	rows, err := r.DB.Read.Queryx(searchgolonganQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var Golongan Golongan
		err = rows.StructScan(&Golongan)
		if err != nil {
			return
		}

		data.Items = append(data.Items, Golongan)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *GolonganRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []Golongan, err error) {
	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}
	rows, err := r.DB.Read.Queryx(golonganQuery.Select + where)

	if err == sql.ErrNoRows {
		_ = failure.NotFound("Golongan NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList Golongan
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *GolonganRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data Golongan, err error) {
	err = r.DB.Read.Get(&data, golonganQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *GolonganRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(golonganQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *GolonganRepositoryPostgreSQL) Update(data Golongan) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *GolonganRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data Golongan) (err error) {
	stmt, err := tx.PrepareNamed(golonganQuery.Update + " WHERE id=:id")
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

func (r *GolonganRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, golonganQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *GolonganRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, golonganQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *GolonganRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, golonganQuery.ExistRelasi, id)

	return
}
