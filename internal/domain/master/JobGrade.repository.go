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
	jobGradeQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, kode, nama, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_job_grade `,
		Insert: `insert into m_job_grade
				(id, kode, nama, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :kode, :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_job_grade set
				id=:id,
				kode=:kode,
				nama=:nama,
				tenant_id=:tenant_id,
				id_branch=:id_branch,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_job_grade `,
		Count: `select count (id)
				from m_job_grade `,
		Exist: `select count(id)>0 from m_job_grade `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where id = $1
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type JobGradeRepository interface {
	Create(data JobGrade) error
	GetAll(req model.StandardRequest) (data []JobGrade, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data JobGrade, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data JobGrade) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID) (exist bool)
}

type JobGradeRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJobGradeRepositoryPostgreSQL(db *infras.PostgresqlConn) *JobGradeRepositoryPostgreSQL {
	s := new(JobGradeRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JobGradeRepositoryPostgreSQL) Create(data JobGrade) error {
	stmt, err := r.DB.Write.PrepareNamed(jobGradeQuery.Insert)
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

func (r *JobGradeRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
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

	query := r.DB.Read.Rebind(jobGradeQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappJobGrade[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchjobGradeQuery := searchRoleBuff.String()
	searchjobGradeQuery = r.DB.Read.Rebind(jobGradeQuery.Select + searchjobGradeQuery)
	fmt.Println("query", searchjobGradeQuery)
	rows, err := r.DB.Read.Queryx(searchjobGradeQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var JobGrade JobGrade
		err = rows.StructScan(&JobGrade)
		if err != nil {
			return
		}

		data.Items = append(data.Items, JobGrade)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *JobGradeRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []JobGrade, err error) {
	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(jobGradeQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("JobGrade NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JobGrade
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *JobGradeRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data JobGrade, err error) {
	err = r.DB.Read.Get(&data, jobGradeQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JobGradeRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(jobGradeQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *JobGradeRepositoryPostgreSQL) Update(data JobGrade) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *JobGradeRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data JobGrade) (err error) {
	stmt, err := tx.PrepareNamed(jobGradeQuery.Update + " WHERE id=:id")
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

func (r *JobGradeRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	err := r.DB.Read.Get(&exist, jobGradeQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JobGradeRepositoryPostgreSQL) ExistByKode(kode string, id string, IdBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if IdBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", IdBranch)
	}

	err := r.DB.Read.Get(&exist, jobGradeQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *JobGradeRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID) (exist bool) {
	r.DB.Read.Get(&exist, jobGradeQuery.ExistRelasi, id)

	return
}
