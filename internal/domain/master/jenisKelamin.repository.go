package master

import (
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	jenisKelaminQuery = struct {
		SelectJenisKelamin string
	}{
		SelectJenisKelamin: `select jk.kode, jk.nama from public.m_jenis_kelamin jk `,
	}
)

type JenisKelaminRepository interface {
	GetAll() (dataJenisKelamin []JenisKelamin, err error)
}

type JenisKelaminRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideJenisKelaminRepositoryPostgreSQL(db *infras.PostgresqlConn) *JenisKelaminRepositoryPostgreSQL {
	s := new(JenisKelaminRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *JenisKelaminRepositoryPostgreSQL) GetAll() (dataJenisKelamin []JenisKelamin, err error) {
	err = r.DB.Read.Select(&dataJenisKelamin, jenisKelaminQuery.SelectJenisKelamin)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}
