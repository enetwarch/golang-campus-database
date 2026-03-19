package main

import (
	"bufio"
	"campus/database"
	"campus/ui"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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
	fmt.Printf("Database Initialized...\n")

	in := ui.Input{Reader: bufio.NewReader(os.Stdin)}
	const ( // Makeshift Golang enums for switch case readability.
		StudentTable    = 1
		CourseTable     = 2
		ProfessorTable  = 3
		EnrollmentTable = 4
		ExitProgram     = 0
	)
	for {
		ui.PrintTable([][]string{
			{"1", "Student Table"},
			{"2", "Course Table"},
			{"3", "Professor Table"},
			{"4", "Enrollment Table"},
			{"0", "Exit Program"},
		})
		switch in.Int("Enter Action (0/1/2/3/4): ", 0, 4) {
		case StudentTable:
			{
				table(db, in, &TableData{
					menuName:   "Student",
					tableName:  "student",
					primaryKey: "student_id",
					columns:    []string{"student_id", "student_name", "block"},
					create: func() []any {
						studentID := in.Int("Student ID: ", 00000001, 99999999)
						studentName := in.String("Student Name: ")
						block := in.String("Block: ")
						return []any{studentID, studentName, block}
					},
					id: func() any {
						return in.Int("Student ID: ", 00000001, 99999999)
					},
				})
			}
		case CourseTable:
		case ProfessorTable:
		case EnrollmentTable:
		case ExitProgram:
			fmt.Printf("Thank you for using our program!")
			return
		}
	}
}

type TableData struct {
	menuName   string
	tableName  string
	primaryKey string
	columns    []string
	create     func() []any
	id         func() any
}

func table(db database.Database, in ui.Input, td *TableData) {
	const ( // Makeshift Golang enums for switch case readability.
		AddRow    = 1
		ViewTable = 2
		EditRow   = 3
		DeleteRow = 4
		ExitTable = 0
	)

	for {
		ui.PrintTable([][]string{
			{"1", fmt.Sprintf("Add %s", td.menuName)},
			{"2", fmt.Sprintf("View %ss", td.menuName)},
			{"3", fmt.Sprintf("Edit %s", td.menuName)},
			{"4", fmt.Sprintf("Delete %s", td.menuName)},
			{"0", "Exit Table"},
		})

		switch in.Int("Enter Action (0/1/2/3/4): ", 0, 4) {
		case AddRow:
			{
				insert := fmt.Sprintf("INSERT INTO %s (%s)", td.tableName, strings.Join(td.columns, ", "))
				placeholders := strings.Split(strings.Repeat("?", len(td.columns)), "")
				values := fmt.Sprintf("VALUES (%s)", strings.Join(placeholders, ", "))
				query := fmt.Sprintf("%s %s", insert, values)
				_, err := db.SQL.Exec(query, td.create()...)
				if err != nil {
					log.Fatal(err)
				}
			}

		case ViewTable:
			{
				rows, err := db.SQL.Query(fmt.Sprintf("SELECT * FROM %s", td.tableName))
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				table, err := database.StringifyRows(rows, len(td.columns))
				if err != nil {
					log.Fatal(err)
				}
				ui.PrintTable(append([][]string{td.columns}, table...))
			}

		case EditRow:
			{
				id := td.id()
				rows, _ := db.SQL.Query(fmt.Sprintf("SELECT * FROM %s", td.tableName))
				if !rows.Next() { // If inputted ID is not in the table.
					fmt.Printf("No %s found with ID %v.\n", td.tableName, id)
					break
				} else {
					fmt.Printf("%s with ID %v found.\n", td.menuName, id)
				}
				rows.Close()
				queryArguments := append(td.create(), id)
				query := fmt.Sprintf("Update %s SET ", td.tableName)
				for _, column := range td.columns {
					query += fmt.Sprintf("%s = ?, ", column)
				}
				query = query[:len(query)-2] // Truncate the last 2 characters to remove ", "
				query += fmt.Sprintf("WHERE %s = ?", td.primaryKey)
				_, err := db.SQL.Exec(query, queryArguments...)
				if err != nil {
					log.Fatal(err)
				}
			}

		case DeleteRow:
			{
				id := td.id()
				query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", td.tableName, td.primaryKey)
				result, _ := db.SQL.Exec(query, id)
				rowsAffected, _ := result.RowsAffected()
				if rowsAffected == 0 {
					fmt.Printf("No %s found with ID %v.\n", td.tableName, id)
				} else {
					fmt.Printf("Successfully deleted %s with ID %v.\n", td.tableName, id)
				}
			}

		case ExitTable:
			fmt.Printf("Exiting %s table.\n", td.tableName)
			return
		}
		in.Buffer()
	}
}
