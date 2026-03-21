// Package db contains functions that handles sqlite queries.
package db

import (
	"database/sql"
	"fmt"
	"strings"
)

type Database struct {
	SQL    *sql.DB
	Tables []*Table
}

const ( // Enums for type-safe index access.
	StudentTableIndex    = 0
	CourseTableIndex     = 1
	ProfessorTableIndex  = 2
	EnrollmentTableIndex = 3
)

func InitializeDatabase() (*Database, error) {
	sqlite, err := sql.Open("sqlite", "campus.db")
	if err != nil {
		return nil, err
	}
	if _, err := sqlite.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	student := StudentTable()
	course := CourseTable()
	professor := ProfessorTable()
	enrollment := EnrollmentTable(student, course, professor)
	database := Database{
		SQL:    sqlite,
		Tables: []*Table{student, course, professor, enrollment},
	}
	for _, table := range database.Tables {
		database.InitializeTable(table)
	}
	return &database, nil
}

// InitializeTable executes an initialize table query like this.
/*
	CREATE TABLE IF NOT EXISTS enrollments (
		student_id INTEGER,
		course_id INTEGER,
		professor_id INTEGER,
		PRIMARY KEY (student_id, course_id),
		FOREIGN KEY (course_id) REFERENCES course(course_id),
		FOREIGN KEY (professor_id) REFERENCES professor(professor_id),
		FOREIGN KEY (student_id) REFERENCES student(student_id)
	);
*/
func (db *Database) InitializeTable(table *Table) error {
	hasPK := len(table.PrimaryKeys) > 0
	hasFK := len(table.ForeignKeys) > 0
	var query strings.Builder
	fmt.Fprintf(&query, "CREATE TABLE IF NOT EXISTS %s (", table.TableName)
	for i, column := range table.Columns {
		fmt.Fprintf(&query, "%s %s", column.ColumnName, column.ColumnType)
		if hasPK || i < len(table.Columns)-1 {
			query.WriteString(", ")
		}
	}
	// CREATE TABLE IF NOT EXISTS enrollments (student_id INTEGER, course_id INTEGER, ...{, }
	if hasPK {
		fmt.Fprintf(&query, "PRIMARY KEY (%s)", strings.Join(table.PKColumnNames(), ", "))
		if hasFK {
			query.WriteString(", ")
		}
	}
	// PRIMARY KEY (student_id, course_id){, }
	if hasFK {
		for i, foreignKey := range table.ForeignKeys {
			fmt.Fprintf(&query, "FOREIGN KEY (%s) ", foreignKey.ForeignColumn.ColumnName)
			fmt.Fprintf(&query, "REFERENCES %s(%s)",
				foreignKey.ReferenceTable.TableName, foreignKey.ReferenceColumn.ColumnName)
			if i < len(table.ForeignKeys)-1 {
				query.WriteString(", ")
			}
		}
	}
	// FOREIGN KEY (course_id) REFERENCES course(course_id), ...
	query.WriteString(");")
	_, err := db.SQL.Exec(query.String())
	return err
}

// The following will be methods for the database to manipulate tables.
// Assumes the values are in the same order as columns but still validates them.
// Only CRUDs rows with ALL columns, this implementation cannot select specific rows.
// All have two return types, a bool indicating if the error is fatal or not, and the error itself.
// Most of the explanation for how they operate is in the Insert method.

func (db *Database) Insert(table *Table, values []any) (sql.Result, error) {
	// Stringbuilder to build up a query string.
	// This would be a lot shorter in Python.
	var query strings.Builder
	fmt.Fprintf(&query, "INSERT INTO %s (%s)", table.TableName,
		strings.Join(table.ColumnNames(), ", "))
	fmt.Fprintf(&query, " VALUES (%s)",
		strings.TrimSuffix(strings.Repeat("?, ", len(table.Columns)), ", "))
	query.WriteString(";")
	// INSERT INTO {tablename} ({column1, column2, column3, ...}) VALUES (?, ?, ?, ...);
	return db.SQL.Exec(query.String(), values...)
}

func (db *Database) Edit(table *Table, pkValues []any, values []any) (sql.Result, error) {
	var query strings.Builder
	fmt.Fprintf(&query, "UPDATE %s SET ", table.TableName)
	for i, column := range table.Columns {
		fmt.Fprintf(&query, "%s = ?", column.ColumnName)
		if i < len(table.Columns)-1 { // If that is not the last column yet...
			query.WriteString(", ")
		}
	}
	// UPDATE {tablename} SET {column1} = ?, {column2} = ?, ...
	query.WriteString(" WHERE ")
	for i, pkColumn := range table.PrimaryKeys {
		fmt.Fprintf(&query, "%s = ?", pkColumn.ColumnName)
		if i < len(table.PrimaryKeys)-1 { // If that is not the last primary key yet...
			query.WriteString(" AND ")
		}
	}
	query.WriteString(";")
	// WHERE {pkcolumn1} = ? AND {pkcolumn2} = ? AND ...;
	// Both UPDATE... and WHERE... are still in the same query.
	return db.SQL.Exec(query.String(), append(values, pkValues...)...)
}

func (db *Database) View(table *Table) (*sql.Rows, error) {
	var query strings.Builder
	fmt.Fprintf(&query, "SELECT * FROM %s;", table.TableName)
	// SELECT * FROM {tablename};
	return db.SQL.Query(query.String())
}

func (db *Database) Search(table *Table, pkValues []any) (*sql.Rows, error) {
	var query strings.Builder
	fmt.Fprintf(&query, "SELECT * FROM %s WHERE ", table.TableName)
	for i, pkColumn := range table.PrimaryKeys {
		fmt.Fprintf(&query, "%s = ?", pkColumn.ColumnName)
		if i < len(table.PrimaryKeys)-1 {
			query.WriteString(" AND ")
		}
	}
	query.WriteString(";")
	// SELECT * FROM {tablename} WHERE {pkcolumn1} = ? AND {pkcolumn2} = ? AND ...;
	return db.SQL.Query(query.String(), pkValues...)
}

func (db *Database) Delete(table *Table, pkValues []any) (sql.Result, error) {
	var query strings.Builder
	fmt.Fprintf(&query, "DELETE FROM %s WHERE ", table.TableName)
	for i, pkColumn := range table.PrimaryKeys {
		fmt.Fprintf(&query, "%s = ?", pkColumn.ColumnName)
		if i < len(table.PrimaryKeys)-1 {
			query.WriteString(" AND ")
		}
	}
	query.WriteString(";")
	// DELETE FROM {tablename} WHERE {pkcolumn1} = ? AND {pkcolumn2} = ? AND ...;
	return db.SQL.Exec(query.String(), pkValues...)
}
