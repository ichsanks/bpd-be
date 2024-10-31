package auth

import (
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/infras"
	"gitlab.com/GSI-SBY-Team/ptpn-bpd-pusat-be/shared/logger"
)

var (
	logSystemQuery = struct {
		Insert string
	}{

		Insert: `Insert into log_activity(id, actions, jam, keterangan, id_user, platform, ip_address, user_agent, kode) 
		values (:id, :actions, :jam, :keterangan, :id_user, :platform, :ip_address, :user_agent, :kode)`,
	}
)

type LogSystemRepository interface {
	CreateLogSystem(logSystem LogSystem) error
}

type LogSystemRepositoryPostgreSQL struct {
	DB *infras.PostgresqlConn
}

func ProvideLogSystemRepositoryPostgreSQL(db *infras.PostgresqlConn) *LogSystemRepositoryPostgreSQL {
	s := new(LogSystemRepositoryPostgreSQL)
	s.DB = db
	return s
}

func (r *LogSystemRepositoryPostgreSQL) CreateLogSystem(logSystem LogSystem) error {
	stmt, err := r.DB.Read.PrepareNamed(logSystemQuery.Insert)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(logSystem)
	if err != nil {
		logger.ErrorWithStack(err)
		return err
	}
	return nil
}
