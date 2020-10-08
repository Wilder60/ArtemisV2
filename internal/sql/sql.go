package sql

import (
	"database/sql"
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/Wilder60/KeyRing/configs"
	"github.com/Wilder60/KeyRing/internal/domain"
	"github.com/google/uuid"
)

const fmtStr = "host=%s:%s:%s user=%s dbname=%s password=%s sslmode=disable"

type txFn func(*sql.Tx) error

type SQL struct {
	db *sql.DB
}

func New() SQL {
	cfg := configs.Get()
	dsn := fmt.Sprintf(fmtStr,
		cfg.Database.SQL.Project,
		cfg.Database.SQL.Region,
		cfg.Database.SQL.Instance,
		cfg.Database.SQL.User,
		cfg.Database.SQL.Dbname,
		cfg.Database.SQL.Password,
	)

	db, err := sql.Open("cloudsqlpostgres", dsn)
	if err != nil {
		panic(err)
	}

	s := SQL{db}
	err = s.initTable(db)
	return SQL{db}
}

func (s *SQL) Close() error {
	return s.db.Close()
}

// GetKeyRing will take the id from the user how is spending a request
// limit and offest are used for pagination
func (s *SQL) GetKeyRing(id uuid.UUID, limit int64, offest int64) ([]domain.KeyEntry, error) {
	var pagedData []domain.KeyEntry = []domain.KeyEntry{}

	err := s.withTransaction(func(tx *sql.Tx) error {
		rows, err := tx.Query(selectKeyEntry, id, limit, offest)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			entity := domain.KeyEntry{}
			err := rows.Scan(&entity)
			if err == nil {
				pagedData = append(pagedData, entity)
			}
		}
		return rows.Err()
	})

	return pagedData, err
}

// AddKeyRing will take
func (s *SQL) AddKeyRing(entry domain.KeyEntry) (int64, error) {
	var rows int64

	err := s.withTransaction(func(tx *sql.Tx) (err error) {
		insertResult, err := tx.Exec(insertKeyEntry,
			entry.UserID,
			entry.URL,
			entry.Username,
			entry.SiteName,
			entry.SitePassword,
			entry.Folder,
			entry.Notes,
			entry.Favorite,
		)
		if err != nil{
			return 
		}
		rows, err = insertResult.LastInsertId()
		return 
	})

	return rows, err
}

func (s *SQL) UpdateKeyRing(entry domain.KeyEntry) (int64, error) {
	var updateRow int64

	err := s.withTransaction(func(tx *sql.Tx) error {
		updateResult, err := tx.Exec(updateKeyEntry,
			entry.URL,
			entry.Username,
			entry.SiteName,
			entry.SitePassword,
			entry.Folder,
			entry.Folder,
			entry.Favorite,
			entry.ID,
			entry.UserID,
		)
		if err != nil 
		updateResult.RowsAffected()
		return nil
	})

	return updateRow, err
}

func (s *SQL) DeleteKeyRing(int64) (int64, error) {
	err := s.withTransaction(func(tx *sql.Tx) error {
		return nil
	})

	return 0, err
}

func (s *SQL) initTable(db *sql.DB) error {
	err := s.withTransaction(func(tx *sql.Tx) error {
		var err error
		_, err = tx.Exec(createExtension)
		_, err = tx.Exec(createTable)
		return err
	})
	return err
}

func (s *SQL) withTransaction(fn txFn) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		p := recover()
		if p != nil || err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
