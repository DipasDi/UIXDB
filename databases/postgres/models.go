package database

import "database/sql"

var PostgreConnect *sql.DB

type ConnectionRequest struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password"`
	Database string `json:"database" binding:"required"`
	Schema   string `json:"schema"`
}

type ColumnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	PK   bool   `json:"pk"`
}

type Tables struct {
	TablesMap map[string][]ColumnInfo
}

type Key struct {
	FkName     string `json:"fk_name"`
	FromTable  string `json:"from_table"`
	FromColumn string `json:"from_column"`
	ToTable    string `json:"to_table"`
	ToColumn   string `json:"to_column"`
}

type Relation struct {
	ID   string `json:"id"`   // fk_name (имя констрейнта в Postgres)
	From string `json:"from"` // формат "orders.user_id"
	To   string `json:"to"`   // формат "users.id"
}

var ConnInfo ConnectionRequest
