package main

import (
	"bufio"
	"campus/db"
	"campus/ui"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite" // sql driver
)

const ( // Makeshift Golang enums for switch case readability.
	StudentTable    = 1
	CourseTable     = 2
	ProfessorTable  = 3
	EnrollmentTable = 4
	ExitProgram     = 0
)

func main() {
	programName := "Golang Campus Database"
	fmt.Printf("Welcome to %s!\n", programName)
	database, err := db.InitializeDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.SQL.Close()
	fmt.Printf("Database Initialized...\n")
	input := ui.Input{Reader: bufio.NewReader(os.Stdin)}
	input.Buffer()
	fmt.Printf("\n")

	for {
		fmt.Printf("[1] Student Table\n")
		fmt.Printf("[2] Course Table\n")
		fmt.Printf("[3] Professor Table\n")
		fmt.Printf("[4] Enrollment Table\n")
		fmt.Printf("[0] Exit Program\n")
		action := input.Int("Enter Action (0/1/2/3/4): ", 0, 4)
		fmt.Printf("\n")

		switch action {
		case StudentTable:
			handleTable(&input, database, database.Tables[db.StudentTableIndex], Methods{
				inputValues: func() []any {
					studentID := input.Int("Student ID: ", 1, 99999999)
					studentName := input.String("Student Name: ")
					block := input.String("Block: ")
					return []any{studentID, studentName, block}
				},
				inputPrimaryKeys: func() []any {
					studentID := input.Int("Student ID: ", 1, 99999999)
					return []any{studentID}
				},
			})

		case CourseTable:
			handleTable(&input, database, database.Tables[db.CourseTableIndex], Methods{
				inputValues: func() []any {
					courseID := input.Int("Course ID: ", 1, 9999)
					courseName := input.String("Course Name: ")
					units := input.Float("Units: ", 0, 168) // There's only 168 hours max in a week.
					return []any{courseID, courseName, units}
				},
				inputPrimaryKeys: func() []any {
					courseID := input.Int("Course ID: ", 1, 9999)
					return []any{courseID}
				},
			})

		case ProfessorTable:
			handleTable(&input, database, database.Tables[db.ProfessorTableIndex], Methods{
				inputValues: func() []any {
					professorID := input.Int("Professor ID: ", 1, 99999999)
					professorName := input.String("Professor Name: ")
					return []any{professorID, professorName}
				},
				inputPrimaryKeys: func() []any {
					professorID := input.Int("Professor ID: ", 1, 99999999)
					return []any{professorID}
				},
			})

		case EnrollmentTable:
			handleTable(&input, database, database.Tables[db.EnrollmentTableIndex], Methods{
				inputValues: func() []any {
					studentID := input.Int("Student ID: ", 1, 99999999)
					courseID := input.Int("Course ID: ", 1, 9999)
					professorID := input.Int("Professor ID: ", 1, 99999999)
					return []any{studentID, courseID, professorID}
				},
				inputPrimaryKeys: func() []any {
					studentID := input.Int("Student ID: ", 1, 99999999)
					courseID := input.Int("Course ID: ", 1, 9999)
					return []any{studentID, courseID}
				},
			})

		case ExitProgram:
			fmt.Printf("Thank you for using %s!", programName)
			os.Exit(0)
		}
		fmt.Printf("\n")
	}
}

type Methods struct {
	inputValues      func() []any
	inputPrimaryKeys func() []any
}

const ( // Makeshift Golang enums for switch case readability.
	AddRow    = 1
	ViewTable = 2
	EditRow   = 3
	DeleteRow = 4
	ExitTable = 0
)

func handleTable(input *ui.Input, database *db.Database, table *db.Table, methods Methods) {
	for {
		fmt.Printf("Currently in %s table.\n", table.TableName)
		fmt.Printf("[1] Add Record\n")
		fmt.Printf("[2] View Records\n")
		fmt.Printf("[3] Edit Record\n")
		fmt.Printf("[4] Delete Record\n")
		fmt.Printf("[0] Exit Table\n")

		switch input.Int("Enter Action (0/1/2/3/4): ", 0, 4) {
		case AddRow:
			fmt.Printf("Insert a record to %s table.\n", table.TableName)
			valuesToInsert := methods.inputValues()
			result, err := database.Insert(table, valuesToInsert)
			if err != nil {
				fmt.Printf("INSERT RECORD ERROR. %v\n", err)
			} else if affected, err := result.RowsAffected(); err != nil {
				fmt.Printf("ROW AFFECTED ERROR. %v\n", err)
			} else if affected == 0 {
				fmt.Printf("Record with primary key already exists in %s table.\n", table.TableName)
			} else {
				fmt.Printf("Successfully inserted record to %s table!\n", table.TableName)
			}

		case ViewTable:
			fmt.Printf("View all record in %s table.\n", table.TableName)
			rows, err := database.View(table)
			if err != nil {
				fmt.Printf("VIEW TABLE ERROR. %v\n", err)
			}
			defer rows.Close()
			stringifiedRows, err := db.StringifyRows(rows, len(table.Columns))
			if err != nil {
				log.Fatal(err)
			}
			ui.PrintTable(append([][]string{table.ColumnNames()}, stringifiedRows...))

		case EditRow:
			fmt.Printf("Edit a record in %s table.\n", table.TableName)
			primaryKeysToEdit := methods.inputPrimaryKeys()
			fmt.Printf("Input new values for the record.\n")
			valuesToReplaceEdit := methods.inputValues()
			result, err := database.Edit(table, primaryKeysToEdit, valuesToReplaceEdit)
			if err != nil {
				fmt.Printf("EDIT RECORD ERROR. %v\n", err)
			} else if affected, err := result.RowsAffected(); err != nil {
				fmt.Printf("ROW AFFECTED ERROR. %v\n", err)
			} else if affected == 0 {
				fmt.Printf("Record to edit not found in %s table.\n", table.TableName)
			} else {
				fmt.Printf("Successfully edited record in %s table!\n", table.TableName)
			}

		case DeleteRow:
			fmt.Printf("Delete a record in %s table.\n", table.TableName)
			primaryKeysToDelete := methods.inputPrimaryKeys()
			result, err := database.Delete(table, primaryKeysToDelete)
			if err != nil {
				fmt.Printf("DELETE RECORD ERROR. %v\n", err)
			} else if affected, err := result.RowsAffected(); err != nil {
				fmt.Printf("ROW AFFECTED ERROR. %v\n", err)
			} else if affected == 0 {
				fmt.Printf("Record to delete not found in %s table.\n", table.TableName)
			} else {
				fmt.Printf("Record successfully deleted in %s table!\n", table.TableName)
			}

		case ExitTable:
			fmt.Printf("Exiting %s table.\n", table.TableName)
			return
		}
		input.Buffer()
		fmt.Printf("\n")
	}
}
