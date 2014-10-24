package demo_test

import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/jmoiron/sqlx"
import _ "github.com/mattn/go-sqlite3"

func TestSomething(t *testing.T) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Ping()
	assert.Nil(t, err)
}
