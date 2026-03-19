package main

import (
	"campus/database"
	"campus/ui"
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // sql driver
)

func main() {
	sqliteDB, err := sql.Open("sqlite", "campus.db")
	if err != nil {
		log.Fatal(err)
	}
	db := database.Database{SQL: sqliteDB}
	defer db.SQL.Close()

	db.InitializeTables()
	db.InitializeExampleRows()

	rows, err := db.Query("SELECT * FROM student")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	table, err := database.StringifyRows(rows, 3)
	if err != nil {
		log.Fatal(err)
	}
	ui.PrintTable(append([][]string{{"student_id", "student_name", "block"}}, table...))
}
