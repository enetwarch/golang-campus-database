package main

import (
	"bufio"
	"diproglang/input"
	"fmt"
	"os"
)

// Main entry point and menu of the program.
func main() {
	input := input.Input{Reader: bufio.NewReader(os.Stdin)}
	studentDatabase := StudentDatabase{students: []Student{}}
	fmt.Printf("Student database initialized.\n")

	for {
		input.Buffer()
		fmt.Printf("\n")
		fmt.Printf("1. Add Student\n")
		fmt.Printf("2. Edit Student\n")
		fmt.Printf("3. Search Student\n")
		fmt.Printf("4. Delete Student\n")
		fmt.Printf("5. Print Students\n")
		fmt.Printf("0. Exit Program\n")

		action := input.Int("Enter action (1/2/3/4/5/0): ", 0, 5)
		index := -1
		switch action {
		case 0:
			fmt.Printf("Thank you for using our program!\n")
			return
		case 2, 3, 4:
			number := input.Int("Student Number: ", 0, 99999999)
			index = studentDatabase.searchStudent(number)
			if index == -1 {
				fmt.Printf("Student not found.\n")
				continue
			}
		}

		switch action {
		case 1:
			studentDatabase.addStudent(inputStudent(input))
			fmt.Printf("Student successfully added to the database.\n")
		case 2:
			studentDatabase.editStudent(index, inputStudent(input))
			fmt.Printf("Student edited successfully.\n")
		case 3:
			fmt.Printf("Student information found.\n")
			studentDatabase.printStudent(index)
		case 4:
			studentDatabase.deleteStudent(index)
			fmt.Printf("Student removed successfully.\n")
		case 5:
			for _, student := range studentDatabase.students {
				fmt.Printf("%d, %s, %s\n", student.number, student.name, student.department)
			}
		}
	}
}

type Student struct {
	number     int
	name       string
	department string
}

type StudentDatabase struct {
	students []Student
}

// Adds a new student to the last index of the current database slice.
func (sdb *StudentDatabase) addStudent(student Student) {
	sdb.students = append(sdb.students, student)
}

// Replaces student in an index with a new student.
func (sdb *StudentDatabase) editStudent(index int, student Student) {
	sdb.students[index] = student
}

// Useful function that returns the index of a student.
// If no student is found, -1 will be returned as the index.
func (sdb *StudentDatabase) searchStudent(number int) int {
	for index, student := range sdb.students {
		if student.number == number {
			return index
		}
	}
	return -1
}

// Golang does not have a native delete method for its slice.
// The workaround is to append the next elements after the index.
// Example: [0, 1, 2, 3, 4], to delete 2, append([0, 1], 3, 4).
func (sdb *StudentDatabase) deleteStudent(index int) {
	sdb.students = append(sdb.students[:index], sdb.students[index+1:]...)
}

// Prints the student in the inserted index.
func (sdb *StudentDatabase) printStudent(index int) {
	fmt.Printf("Student Number: %d\n", sdb.students[index].number)
	fmt.Printf("Name: %s\n", sdb.students[index].name)
	fmt.Printf("Department: %s\n", sdb.students[index].department)
}

// Helper function to get user input on student data.
func inputStudent(input input.Input) Student {
	return Student{
		number:     input.Int("Student Number: ", 0, 99999999),
		name:       input.String("Name: "),
		department: input.String("Department: "),
	}
}
