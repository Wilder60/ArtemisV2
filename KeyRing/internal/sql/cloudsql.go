package sql

import (
	"database/sql"
	"fmt"

	"github.com/Wilder60/KeyRing/configs"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// const fmtStr = "postgres://%v:%v@%v:%v/%v?sslmode=disable"
const fmtStr = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"

func CreateCloudSQLDriver(config *configs.Config, logger *zap.Logger) (*sql.DB, error) {

	logger.Info(fmt.Sprintf(fmtStr,
		config.Database.Postgres.Port,
		config.Database.Postgres.User,
		config.Database.Postgres.Password,
		config.Database.Postgres.Dbname),
	)
	// config.Database.Postgres.Hostname,
	dsn := fmt.Sprintf(fmtStr,
		config.Database.Postgres.Hostname,
		config.Database.Postgres.Port,
		config.Database.Postgres.User,
		config.Database.Postgres.Password,
		config.Database.Postgres.Dbname)

	client, err := sql.Open("postgres", dsn)
	pingErr := client.Ping()
	if pingErr != nil {
		print(err.Error())
	}
	return client, pingErr
}

var ModuleCloudSql = fx.Option(
	fx.Provide(CreateCloudSQLDriver),
)
