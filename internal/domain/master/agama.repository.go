package master

import (
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	agamaQuery = struct {
		SelectAgama string
	}{
		SelectAgama: `select ma.kode, ma.nama from public.m_agama ma `,
	}
)

type AgamaRepository interface {
	GetAll() (dataAgama []Agama, err error)
}

type AgamaRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideAgamaRepositoryPostgreSQL(db *infras.PostgresqlConn) *AgamaRepositoryPostgreSQL {
	s := new(AgamaRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *AgamaRepositoryPostgreSQL) GetAll() (dataAgama []Agama, err error) {
	err = r.DB.Read.Select(&dataAgama, agamaQuery.SelectAgama)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}
