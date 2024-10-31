package master

import (
	"database/sql"
	"fmt"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/model"
)

var (
	jenisApprovalQuery = struct {
		Select string
	}{
		Select: `select id, nama, id_fungsionalitas from m_jenis_approval `,
	}
)

type JenisApprovalRepository interface {
	GetAll(ids string) (data []JenisApproval, err error)
}

type JenisApprovalRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJenisApprovalRepositoryPostgreSQL(db *infras.PostgresqlConn) *JenisApprovalRepositoryPostgreSQL {
	s := new(JenisApprovalRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JenisApprovalRepositoryPostgreSQL) GetAll(ids string) (data []JenisApproval, err error) {
	where := ""
	if ids != "" {
		id := model.ParseSplitString(ids)
		where += fmt.Sprintf(" where id in (%v) ", id)
	}
	where += " order by id asc"
	rows, err := r.DB.Read.Queryx(jenisApprovalQuery.Select + where)
	fmt.Println(err)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Jenis Approval NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JenisApproval
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}
