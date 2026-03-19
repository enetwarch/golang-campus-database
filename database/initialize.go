// Package database contains functions that handles sqlite queries.
package database

// InitializeTables creates the table in the db if it does not exist yet.
func (db *Database) InitializeTables() error {
	return db.Execute(
		`CREATE TABLE IF NOT EXISTS student (
			student_id INTEGER PRIMARY KEY,
			student_name TEXT,
			block TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS course (
			course_id INTEGER PRIMARY KEY,
			course_name TEXT,
			units FLOAT
		);`,
		`CREATE TABLE IF NOT EXISTS professor (
			professor_id INTEGER PRIMARY KEY,
			teacher_name TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS enrollments (
			student_id INTEGER, FOREIGN KEY(student_id) REFERENCES student(student_id),
			course_id INTEGER, FOREIGN KEY(course_id) REFERENCES course(course_id),
			professor_id INTEGER, FOREIGN KEY(professor_id) REFERENCES professor(professor_id)
		);`,
	)
}

// InitializeExampleRows creates sample values in various tables of the database.
func (db *Database) InitializeExampleRows() error {
	return db.Execute(
		"INSERT OR IGNORE INTO student VALUES (20000001, 'Duya, Arman D.', 'CS-201')",
		"INSERT OR IGNORE INTO student VALUES (20000002, 'Espinosa, Lord Raizen I.', 'CS-201')",
		"INSERT OR IGNORE INTO student VALUES (20000003, 'Molina, Hugo P.', 'CS-201')",
		"INSERT OR IGNORE INTO student VALUES (20000004, 'Panergo, Mikko Brandon B.', 'CS-201')",
		"INSERT OR IGNORE INTO course VALUES (2060, '6DIPROGLANG', 3.0)",
		"INSERT OR IGNORE INTO professor VALUES (20000001, 'Salenga, Ma. Louella')",
	)
}
