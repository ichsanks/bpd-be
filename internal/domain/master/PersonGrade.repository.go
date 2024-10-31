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
	personGradeQuery = struct {
		Select      string
		Insert      string
		Update      string
		Delete      string
		Count       string
		Exist       string
		ExistRelasi string
	}{
		Select: `select id, kode, nama, tenant_id, id_branch, created_at, created_by, updated_at, updated_by, is_deleted from m_person_grade `,
		Insert: `insert into m_person_grade
				(id, kode, nama, tenant_id, id_branch, created_at, created_by)
				values
				(:id, :kode, :nama, :tenant_id, :id_branch, :created_at, :created_by) `,
		Update: `update m_person_grade set
				id=:id,
				kode=:kode,
				nama=:nama,
				tenant_id=:tenant_id,
				updated_at=:updated_at,
				updated_by=:updated_by,
				is_deleted=:is_deleted `,
		Delete: `delete from m_person_grade `,
		Count: `select count (id)
				from m_person_grade `,
		Exist: `select count(id)>0 from m_person_grade `,
		ExistRelasi: `select count(1) > 0 exist from (
			select id from m_pegawai pd 
			where pd.id_person_grade = $1
			and   pd.id_branch=$2
			and coalesce(is_deleted, false) is false 
			for update 
		) x  `,
	}
)

type PersonGradeRepository interface {
	Create(data PersonGrade) error
	GetAll(req model.StandardRequest) (data []PersonGrade, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data PersonGrade, err error)
	DeleteByID(id uuid.UUID) (err error)
	Update(data PersonGrade) error
	ExistByNama(nama string, id string, idBranch string) (bool, error)
	ExistByKode(kode string, id string, idBranch string) (bool, error)
	ExistRelasiStatus(id uuid.UUID, idBranch string) (exist bool)
}

type PersonGradeRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvidePersonGradeRepositoryPostgreSQL(db *infras.PostgresqlConn) *PersonGradeRepositoryPostgreSQL {
	s := new(PersonGradeRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *PersonGradeRepositoryPostgreSQL) Create(data PersonGrade) error {
	stmt, err := r.DB.Write.PrepareNamed(personGradeQuery.Insert)
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

func (r *PersonGradeRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
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

	query := r.DB.Read.Rebind(personGradeQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappPersonGrade[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchpersonGradeQuery := searchRoleBuff.String()
	searchpersonGradeQuery = r.DB.Read.Rebind(personGradeQuery.Select + searchpersonGradeQuery)
	fmt.Println("query", searchpersonGradeQuery)
	rows, err := r.DB.Read.Queryx(searchpersonGradeQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var PersonGrade PersonGrade
		err = rows.StructScan(&PersonGrade)
		if err != nil {
			return
		}

		data.Items = append(data.Items, PersonGrade)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *PersonGradeRepositoryPostgreSQL) GetAll(req model.StandardRequest) (data []PersonGrade, err error) {
	where := " where coalesce(is_deleted, false) = false "
	if req.IdBranch != "" {
		where += fmt.Sprintf(" and id_branch='%s' ", req.IdBranch)
	}

	rows, err := r.DB.Read.Queryx(personGradeQuery.Select + where)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("PersonGrade NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList PersonGrade
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *PersonGradeRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data PersonGrade, err error) {
	err = r.DB.Read.Get(&data, personGradeQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PersonGradeRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(personGradeQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *PersonGradeRepositoryPostgreSQL) Update(data PersonGrade) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *PersonGradeRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data PersonGrade) (err error) {
	stmt, err := tx.PrepareNamed(personGradeQuery.Update + " WHERE id=:id")
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

func (r *PersonGradeRepositoryPostgreSQL) ExistByNama(nama string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(nama)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}
	err := r.DB.Read.Get(&exist, personGradeQuery.Exist+criteria, nama)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *PersonGradeRepositoryPostgreSQL) ExistByKode(kode string, id string, idBranch string) (bool, error) {
	var exist bool

	criteria := ` where upper(kode)=upper($1) and coalesce(is_deleted, false)=false `
	if id != "" {
		criteria += fmt.Sprintf(" and id <> '%s' ", id)
	}

	if idBranch != "" {
		criteria += fmt.Sprintf(" and id_branch = '%s' ", idBranch)
	}

	err := r.DB.Read.Get(&exist, personGradeQuery.Exist+criteria, kode)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return exist, err
}

func (r *PersonGradeRepositoryPostgreSQL) ExistRelasiStatus(id uuid.UUID, idBranch string) (exist bool) {
	r.DB.Read.Get(&exist, personGradeQuery.ExistRelasi, id, idBranch)

	return
}
