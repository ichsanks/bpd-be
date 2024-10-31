package master

import (
	"github.com/gofrs/uuid"
)

type Log struct {
	Id         uuid.UUID `db:"id" json:"id"`
	Actions    *string   `db:"actions" json:"actions"`
	Jam        *string   `db:"jam" json:"jam"`
	Kode       *string   `db:"kode" json:"kode"`
	Keterangan *string   `db:"keterangan" json:"keterangan"`
	IdUser     uuid.UUID `db:"id_user" json:"idUser"`
	Username   *string   `db:"username" json:"username"`
	Platform   *string   `db:"platform" json:"platform"`
	IpAddress  *string   `db:"ip_address" json:"ipAddress"`
	UserAgent  *string   `db:"user_agent" json:"userAgent"`
}

var ColumnMappLog = map[string]interface{}{
	"id":         "id",
	"actions":    "actions",
	"jam":        "jam",
	"kode":       "kode",
	"keterangan": "keterangan",
	"idUser":     "id_user",
	"username":   "username",
	"platform":   "platform",
	"ipAddress":  "ip_address",
	"userAgent":  "user_agent",
}
