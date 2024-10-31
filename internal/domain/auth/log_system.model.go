package auth

import (
	"time"

	"github.com/gofrs/uuid"
)

type LogSystem struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Actions    string    `json:"actions" db:"actions"`
	Jam        time.Time `json:"jam" db:"jam"`
	Keterangan string    `json:"keterangan" db:"keterangan"`
	IdUser     uuid.UUID `json:"idUser" db:"id_user"`
	Platform   string    `json:"platform" db:"platform"`
	IpAddress  string    `json:"ipAddress" db:"ip_address"`
	UserAgent  string    `json:"userAgent" db:"user_agent"`
	Kode       string    `json:"kode" db:"kode"`
}

type RequestLogSystemFormat struct {
	Actions    string `json:"actions" db:"actions"`
	Keterangan string `json:"keterangan" db:"keterangan"`
	Kode       string `json:"kode" db:"kode"`
}

func (logSystem *LogSystem) NewLogSystemFormat(reqFormat RequestLogSystemFormat, userId uuid.UUID, ipAddress string, userAgent string) (newLogSystem LogSystem, err error) {
	newID, _ := uuid.NewV4()

	newLogSystem = LogSystem{
		ID:         newID,
		Actions:    reqFormat.Actions,
		Jam:        time.Now(),
		Keterangan: reqFormat.Keterangan,
		IdUser:     userId,
		Platform:   "WEB",
		IpAddress:  ipAddress,
		UserAgent:  userAgent,
		Kode:       reqFormat.Kode,
	}

	return
}
