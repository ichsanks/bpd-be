package master

import (
	"database/sql"
	"fmt"

	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/failure"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	jenisSppdQuery = struct {
		Select string
	}{
		Select: `select id, nama from m_jenis_sppd `,
	}
)

type JenisSppdRepository interface {
	GetAll() (data []JenisSppd, err error)
}

type JenisSppdRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJenisSppdRepositoryPostgreSQL(db *infras.PostgresqlConn) *JenisSppdRepositoryPostgreSQL {
	s := new(JenisSppdRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JenisSppdRepositoryPostgreSQL) GetAll() (data []JenisSppd, err error) {
	where := ""
	where += " order by id asc"
	rows, err := r.DB.Read.Queryx(jenisSppdQuery.Select + where)
	fmt.Println(err)
	if err == sql.ErrNoRows {
		_ = failure.NotFound("Jenis Sppd NotFound")
		return
	}

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	for rows.Next() {
		var dataList JenisSppd
		err = rows.StructScan(&dataList)

		if err != nil {
			return
		}

		data = append(data, dataList)
	}
	return
}
