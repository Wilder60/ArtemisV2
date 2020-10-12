package sql

import (
	"database/sql"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/Wilder60/KeyRing/internal/domain"
	"github.com/Wilder60/KeyRing/internal/interfaces"
)

// fmtStr is the connection string that cloudsql

// SQL stores a *sql database connection for processing requests
//
// It implements the interfaces.database interface so we can use it for
// our dependency injection
type SQL struct {
	db interfaces.SQLDriver
}

// NOTE: FIX this shit... tomorrow
// New returns a new instance of the SQL struct with a connection to the cloudsql
func New(dbCtn interfaces.SQLDriver) SQL {
	s := SQL{db}
	s.init()
	return SQL{db}
}

func (s *SQL) init() {
	err := s.withTransaction(func(tx *sql.Tx) error {
		var err error
		_, err = tx.Exec(createExtension)
		_, err = tx.Exec(createTable)
		return err
	})
	if err != nil {
		panic(err)
	}
	return
}

// Close will close the connection to the http server
func (s *SQL) Close() error {
	return s.db.Close()
}

// GetKeyRing will take the id from the user how is spending a request
// limit and offest are used for pagination
func (s *SQL) GetKeyRing(id string, limit int64, offest int64) ([]domain.KeyEntry, error) {
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
func (s *SQL) AddKeyRing(entry domain.KeyEntry, userID string) (int64, error) {
	var rows int64
	err := s.withTransaction(func(tx *sql.Tx) (err error) {
		insertResult, err := tx.Exec(insertKeyEntry,
			userID, entry.URL, entry.Username,
			entry.SiteName, entry.SitePassword,
			entry.Folder, entry.Notes, entry.Favorite,
		)
		if err != nil {
			return
		}
		rows, err = insertResult.RowsAffected()
		return
	})
	return rows, err
}

// UpdateKeyRing will take a KeyEntry
func (s *SQL) UpdateKeyRing(entry domain.KeyEntry, userID string) (int64, error) {
	var updateRow int64
	err := s.withTransaction(func(tx *sql.Tx) (err error) {
		updateResult, err := tx.Exec(updateKeyEntry,
			entry.URL, entry.Username, entry.SiteName,
			entry.SitePassword, entry.Folder, entry.Folder,
			entry.Favorite, entry.ID, userID,
		)
		if err != nil {
			return
		}
		updateRow, err = updateResult.RowsAffected()
		return
	})
	return updateRow, err
}

// DeleteKeyRing will take a string id related to a request
func (s *SQL) DeleteKeyRing(eventID string) (int64, error) {
	var rowDeleted int64
	err := s.withTransaction(func(tx *sql.Tx) (err error) {
		deleteRequest, err := tx.Exec(deleteKeyEntry, eventID)
		if err != nil {
			return
		}
		rowDeleted, err = deleteRequest.RowsAffected()
		return
	})
	return rowDeleted, err
}

func (s *SQL) withTransaction(fn func(*sql.Tx) error) error {
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
