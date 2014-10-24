package demo_test

import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/jmoiron/sqlx"
import _ "github.com/mattn/go-sqlite3"

type Country struct {
	Code string
	Name string
}

type City struct {
	Code        string
	Name        string
	CountryCode string `db:"country_code"`
}

func Test(t *testing.T) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Ping()
	assert.Nil(t, err)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Country (
		code  TEXT  PRIMARY KEY  NOT NULL,
		name  TEXT
	);
	CREATE TABLE IF NOT EXISTS City (
		code          TEXT  PRIMARY KEY  NOT NULL,
		name          TEXT,
		country_code  TEXT
	);`)
	assert.Nil(t, err)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO Country (code, name) VALUES (?, ?)", "PT", "Portugal")
	tx.MustExec("INSERT INTO Country (code, name) VALUES (?, ?)", "ES", "Spain")
	tx.MustExec("INSERT INTO City (code, name, country_code) VALUES (?, ?, ?)", "OPO", "Porto", "PT")
	tx.MustExec("INSERT INTO City (code, name, country_code) VALUES (?, ?, ?)", "LIS", "Lisbon", "PT")
	tx.MustExec("INSERT INTO City (code, name, country_code) VALUES (?, ?, ?)", "BAR", "Barcelona", "SP")
	tx.MustExec("INSERT INTO City (code, name, country_code) VALUES (?, ?, ?)", "MAD", "Madrid", "SP")
	err = tx.Commit()
	assert.Nil(t, err)

	rows, err := db.Queryx("SELECT code, name FROM Country")
	assert.Nil(t, err)
	for rows.Next() {
		country := Country{}
		err = rows.StructScan(&country)
		assert.Nil(t, err)
	}

	country := Country{}
	err = db.Get(&country, "SELECT * FROM Country WHERE code = ?", "PT")
	assert.Nil(t, err)
	assert.Equal(t, Country{Code: "PT", Name: "Portugal"}, country)

	cities := []City{}
	err = db.Select(&cities, "SELECT * FROM City")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(cities))

	tx = db.MustBegin()
	tx.MustExec("UPDATE City SET name = ? WHERE code = ?", "Lisboa", "LIS")
	err = tx.Commit()
	assert.Nil(t, err)

	lisbon := City{}
	err = db.Get(&lisbon, "SELECT name FROM City WHERE code = ?", "LIS")
	assert.Nil(t, err)
	assert.Equal(t, "Lisboa", lisbon.Name)
	assert.Equal(t, "", lisbon.Code) // because code column is not part of the select statement
}
