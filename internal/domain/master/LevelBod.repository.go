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
	levelBodQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, kode, nama, level, tenant_id, id_branch,  created_at, created_by, updated_at, updated_by, is_deleted from m_level_bod `,
		Insert: `insert into m_level_bod
				(id, kode, nama, level, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :kode, :nama, :level, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_level_bod set
				id=:id,
				kode=:kode,
				nama=:nama,
				level=:level,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_level_bod `,
		Count: `select count (id)
				from m_level_bod `,
		Exist: `select count(id)>0 from m_level_bod `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type LevelBodRepository interface {
	Create(data LevelBod) error
	GetAll(req model.StandardRequest) (data []LevelBod, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data LevelBod, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data LevelBod) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type LevelBodRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideLevelBodRepositoryPostgreSQL(db *infras.PostgresqlConn) *LevelBodRepositoryPostgreSQL {
	s := new(LevelBodRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *LevelBodRepositoryPostgreSQL) Create(data LevelBod) error {
	stmt, err := r.DB.Write.PrepareNamed(levelBodQuery.Insert)
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

func (r *LevelBodRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
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

	query := r.DB.Read.Rebind(levelBodQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappLevelBod[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchlevelBodQuery := searchRoleBuff.String()
	searchlevelBodQuery = r.DB.Read.Rebind(levelBodQuery.Select + searchlevelBodQuery)
	fmt.Println("query", searchlevelBodQuery)
	rows, err := r.DB.Read.Queryx(searchlevelBodQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var LevelBod LevelBod
		err = rows.StructScan(&LevelBod)
		if err != nil {
			return
		}

		data.Items = append(data.Items, LevelBod)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *LevelBodRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []LevelBod, err error) {
	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(levelBodQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("LevelBod NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList LevelBod
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *LevelBodRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data LevelBod, err error) {
	err = r.DB.Read.Get(&data, levelBodQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *LevelBodRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(levelBodQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *LevelBodRepositoryPostgreSQL) Update(data LevelBod) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *LevelBodRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data LevelBod) (err error) {
	stmt, err := tx.PrepareNamed(levelBodQuery.Update + " WHERE id=:id")
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

func (r *LevelBodRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, levelBodQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *LevelBodRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, levelBodQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *LevelBodRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, levelBodQuery.ExistRelasi, id)

	return
}
