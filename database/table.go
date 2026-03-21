// Package database contains functions that handles sqlite queries.
package database

// Table struct for easier manipulation of SQL tables during runtime.
type Table struct {
	TableName   string
	Columns     []Column
	PrimaryKeys []*Column
	ForeignKeys []*ForeignKey
}

type Column struct {
	ColumnName string
	ColumnType string // Name of the type in SQLITE.
}

type ForeignKey struct {
	ForeignColumn   *Column
	ReferenceTable  *Table
	ReferenceColumn *Column
}

// The following are getter methods for the tables in the database.
// They are in the form of functions to make them susceptible to garbage collection when out of scope.

func StudentTable() *Table {
	columns := []Column{
		{ColumnName: "student_id", ColumnType: "INTEGER"},
		{ColumnName: "student_name", ColumnType: "TEXT"},
		{ColumnName: "block", ColumnType: "TEXT"},
	}
	return &Table{
		TableName:   "student",
		Columns:     columns,
		PrimaryKeys: []*Column{&columns[0]},
	}
}

func CourseTable() *Table {
	columns := []Column{
		{ColumnName: "course_id", ColumnType: "INTEGER"},
		{ColumnName: "course_name", ColumnType: "TEXT"},
		{ColumnName: "units", ColumnType: "FLOAT"},
	}
	return &Table{
		TableName:   "course",
		Columns:     columns,
		PrimaryKeys: []*Column{&columns[0]},
	}
}

func ProfessorTable() *Table {
	columns := []Column{
		{ColumnName: "professor_id", ColumnType: "INTEGER"},
		{ColumnName: "professor_name", ColumnType: "TEXT"},
	}
	return &Table{
		TableName:   "professor",
		Columns:     columns,
		PrimaryKeys: []*Column{&columns[0]},
	}
}

func EnrollmentTable(student *Table, course *Table, professor *Table) *Table {
	columns := []Column{
		{ColumnName: "student_id", ColumnType: "INTEGER"},
		{ColumnName: "course_id", ColumnType: "INTEGER"},
		{ColumnName: "professor_id", ColumnType: "INTEGER"},
	}
	return &Table{
		TableName:   "enrollment",
		Columns:     columns,
		PrimaryKeys: []*Column{&columns[0], &columns[1]},
		ForeignKeys: []*ForeignKey{
			{ForeignColumn: &columns[0], ReferenceTable: student, ReferenceColumn: &student.Columns[0]},
			{ForeignColumn: &columns[1], ReferenceTable: course, ReferenceColumn: &course.Columns[0]},
			{ForeignColumn: &columns[2], ReferenceTable: professor, ReferenceColumn: &professor.Columns[0]},
		},
	}
}
