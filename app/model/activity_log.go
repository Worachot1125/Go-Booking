package model

import "github.com/uptrace/bun"

type ActivityLog struct {
	bun.BaseModel `bun:"table:activity_logs"`

	ID         string      `json:"id" bun:",pk,type:uuid,default:gen_random_uuid()" form:"id"`
	Section    string      `json:"section" form:"section"`
	EventType  string      `json:"event_type" bun:",notnull,type:varchar(50)" form:"event_type"`
	StatusCode int         `json:"status_code" bun:",notnull" form:"status_code"`
	Parameters interface{} `json:"parameters" bun:"type:jsonb" form:"parameters"`
	Responses  interface{} `json:"responses" bun:"type:jsonb" form:"responses"`
	Query      interface{} `json:"query" bun:"type:jsonb" form:"query"`
	IpAddress  string      `json:"ip_address" bun:",notnull,type:varchar(50)" form:"ip_address"`
	UserAgent  string      `json:"user_agent" form:"user_agent"`
	CreatedBy  int64       `json:"created_by" form:"created_by"`
	CreatedAt  int64       `json:"created_at" bun:"default:EXTRACT(EPOCH FROM NOW())" form:"created_at"`
}
