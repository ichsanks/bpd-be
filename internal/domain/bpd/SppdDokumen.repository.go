package bpd

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
	sppdDokumenQuery = struct {
		Select string
		Insert string
		Update string
		Delete string
		Count  string
		Exist  string
	}{
		Select: `select 
					id, id_sppd, file, keterangan, id_dokumen, created_at, created_by, updated_at, updated_by, is_deleted 
				from 
					sppd_dokumen `,
		Insert: `insert into sppd_dokumen
					(id, id_sppd, file, keterangan, id_dokumen, created_at, created_by) 
				values 
					(:id, :id_sppd, :file, :keterangan, :id_dokumen, :created_at, :created_by) `,
		Update: `update 
					sppd_dokumen 
				 set
					id=:id,
					id_sppd=:id_sppd,
					file=:file,
					keterangan=:keterangan,
					id_dokumen=:id_dokumen,
					updated_at=:updated_at,
					updated_by=:updated_by,
					is_deleted=:is_deleted `,
		Delete: `delete from sppd_dokumen `,
		Count: `select 
					count (id)
				from 
					sppd_dokumen `,
		Exist: `select 
					count(id)>0 
				from 
					sppd_dokumen `,
	}
)

type SppdDokumenRepository interface {
	Create(data SppdDokumen) error
	Update(data SppdDokumen) error
	GetAll() (data []SppdDokumen, err error)
	ResolveAll(req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(id uuid.UUID) (data SppdDokumen, err error)
	DeleteByID(id uuid.UUID) (err error)
}

type SppdDokumenRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideSppdDokumenRepositoryPostgreSQL(db *infras.PostgresqlConn) *SppdDokumenRepositoryPostgreSQL {
	s := new(SppdDokumenRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *SppdDokumenRepositoryPostgreSQL) Create(data SppdDokumen) error {
	stmt, err := r.DB.Write.PrepareNamed(sppdDokumenQuery.Insert)
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

func (r *SppdDokumenRepositoryPostgreSQL) ResolveAll(req model.StandardRequest) (data pagination.Response, err error) {
	var searchParams []interface{}
	var searchRoleBuff bytes.Buffer
	searchRoleBuff.WriteString(" WHERE coalesce(is_deleted, false) = false ")

	if req.Keyword != "" {
		searchRoleBuff.WriteString(" AND ")
		searchRoleBuff.WriteString(" concat(keterangan) ilike ?  ")
		searchParams = append(searchParams, "%"+req.Keyword+"%")
	}

	query := r.DB.Read.Rebind(sppdDokumenQuery.Count + searchRoleBuff.String())

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

	searchRoleBuff.WriteString("order by " + ColumnMappSppdDokumen[req.SortBy].(string) + " " + req.SortType + " ")

	offset := (req.PageNumber - 1) * req.PageSize
	searchRoleBuff.WriteString("limit ? offset ? ")
	searchParams = append(searchParams, req.PageSize)
	searchParams = append(searchParams, offset)

	searchbidangQuery := searchRoleBuff.String()
	searchbidangQuery = r.DB.Read.Rebind(sppdDokumenQuery.Select + searchbidangQuery)
	fmt.Println("query", searchbidangQuery)
	rows, err := r.DB.Read.Queryx(searchbidangQuery, searchParams...)
	if err != nil {
		return
	}
	for rows.Next() {
		var item SppdDokumen
		err = rows.StructScan(&item)
		if err != nil {
			return
		}

		data.Items = append(data.Items, item)
	}

	data.Meta = pagination.CreateMeta(totalData, req.PageSize, req.PageNumber)
	return
}

func (r *SppdDokumenRepositoryPostgreSQL) GetAll() (data []SppdDokumen, err error) {
	rows, err := r.DB.Read.Queryx(sppdDokumenQuery.Select + " where coalesce(is_deleted, false) = false")
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Sppd Dokumen NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList SppdDokumen
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}

func (r *SppdDokumenRepositoryPostgreSQL) ResolveByID(id uuid.UUID) (data SppdDokumen, err error) {
	err = r.DB.Read.Get(&data, sppdDokumenQuery.Select+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *SppdDokumenRepositoryPostgreSQL) DeleteByID(id uuid.UUID) (err error) {
	_, err = r.DB.Read.Query(sppdDokumenQuery.Delete+" WHERE id=$1  ", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *SppdDokumenRepositoryPostgreSQL) Update(data SppdDokumen) error {
	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, data); err != nil {
			e <- err
			return
		}
		e <- nil
	})
}

func (r *SppdDokumenRepositoryPostgreSQL) txUpdate(tx *sqlx.Tx, data SppdDokumen) (err error) {
	stmt, err := tx.PrepareNamed(sppdDokumenQuery.Update + " WHERE id=:id")
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
