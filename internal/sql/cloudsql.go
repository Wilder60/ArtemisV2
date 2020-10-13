package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/Wilder60/KeyRing/configs"
	"go.uber.org/fx"
)

const fmtStr = "host=%s:%s:%s user=%s dbname=%s password=%s sslmode=disable"

func CreateCloudSQLDriver(config *configs.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(fmtStr,
		config.Database.Postgres.Project,
		config.Database.Postgres.Region,
		config.Database.Postgres.Instance,
		config.Database.Postgres.User,
		config.Database.Postgres.Dbname,
		config.Database.Postgres.Password,
	)
	return sql.Open(config.Database.Postgres.Type, dsn)
}

var ModuleCloudSql = fx.Option(
	fx.Provide(CreateCloudSQLDriver),
)
