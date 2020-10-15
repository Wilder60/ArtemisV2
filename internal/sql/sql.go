package sql

import (
	"database/sql"

	"go.uber.org/fx"

	"github.com/Wilder60/KeyRing/internal/domain"
)

// fmtStr is the connection string that cloudsql

// SQL stores a *sql database connection for processing requests
//
// It implements the interfaces.database interface so we can use it for
// our dependency injection
type SQL struct {
	db *sql.DB
}

// CreateSQL will
func CreateSQL(dbCtn *sql.DB) *SQL {
	s := &SQL{db: dbCtn}
	s.init()
	return s
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
			e := domain.KeyEntry{}
			var userid string //we don't need to return this to the user
			err := rows.Scan(&e.ID, &userid, &e.URL, &e.SiteName, &e.Folder,
				&e.Username, &e.SitePassword, &e.Notes, &e.Favorite)
			if err == nil {
				pagedData = append(pagedData, e)
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

var KeyRingSQLModule = fx.Option(
	fx.Provide(CreateSQL),
)
