package provider

import (
	"database/sql"
	"fmt"

	"github.com/Wilder60/KeyRing/configs"
	"github.com/Wilder60/KeyRing/internal/interfaces"
)

type SQLProvider interface {
	CreateDriver() (interfaces.SQLDriver, error)
}

const fmtStr = "host=%s:%s:%s user=%s dbname=%s password=%s sslmode=disable"

type CloudSQLProvider struct {
}

func (p *CloudSQLProvider) CreateDriver() (*sql.DB, error) {
	cfg := configs.Get()
	dsn := fmt.Sprintf(fmtStr,
		cfg.Database.SQL.Project,
		cfg.Database.SQL.Region,
		cfg.Database.SQL.Instance,
		cfg.Database.SQL.User,
		cfg.Database.SQL.Dbname,
		cfg.Database.SQL.Password,
	)

	// Since I can have different types of sql databases, I should be using
	// dependency Inject for this

	return sql.Open("cloudsqlpostgres", dsn)
}
