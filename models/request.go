package models

import "github.com/blinkinglight/pocketbase-mysql/tools/types"

var _ Model = (*Request)(nil)

// list with the supported values for `Request.Auth`
const (
	RequestAuthGuest = "guest"
	RequestAuthUser  = "user"
	RequestAuthAdmin = "admin"
)

type Request struct {
	BaseModel

	Url       string        `db:"url" json:"url"`
	Method    string        `db:"method" json:"method"`
	Status    int           `db:"status" json:"status"`
	Auth      string        `db:"auth" json:"auth"`
	Ip        string        `db:"ip" json:"ip"`
	Referer   string        `db:"referer" json:"referer"`
	UserAgent string        `db:"userAgent" json:"userAgent"`
	Meta      types.JsonMap `db:"meta" json:"meta"`
}

func (m *Request) TableName() string {
	return "_requests"
}
