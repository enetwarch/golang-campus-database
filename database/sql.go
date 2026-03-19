// Package database contains functions that handles sqlite queries.
package database

import "database/sql"

type Database struct {
	SQL *sql.DB
}

func (db *Database) Execute(queries ...string) error {
	for _, query := range queries {
		_, err := db.SQL.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) Query(query string) (*sql.Rows, error) {
	return db.SQL.Query(query)
}

func StringifyRows(rows *sql.Rows, columns int) ([][]string, error) {
	var table [][]string
	for rows.Next() {
		column := make([]string, columns)
		destination := make([]any, columns)
		for i := range column {
			destination[i] = &column[i]
		}
		err := rows.Scan(destination...)
		if err != nil {
			return nil, err
		}
		table = append(table, column)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return table, nil
}
