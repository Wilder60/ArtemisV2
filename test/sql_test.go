package test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetKeyRing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec()

}
