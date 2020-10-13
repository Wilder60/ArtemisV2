package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/Wilder60/KeyRing/configs"
	"github.com/Wilder60/KeyRing/internal/interfaces"
	"go.uber.org/fx"
)

type SQLProvider interface {
	CreateDriver() (interfaces.SQLDriver, error)
}

const fmtStr = "host=%s:%s:%s user=%s dbname=%s password=%s sslmode=disable"

func CreateCloudSQLDriver(config *configs.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(fmtStr,
		config.Database.SQL.Project,
		config.Database.SQL.Region,
		config.Database.SQL.Instance,
		config.Database.SQL.User,
		config.Database.SQL.Dbname,
		config.Database.SQL.Password,
	)
	return sql.Open("cloudsqlpostgres", dsn)
}

var Module = fx.Option(
	fx.Provide(CreateCloudSQLDriver),
)
