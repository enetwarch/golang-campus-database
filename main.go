package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Main entry point and menu of the program.
func main() {
	studentDatabase := StudentDatabase{students: []Student{}}
	fmt.Printf("Student database initialized.\n")

	for {
		fmt.Printf("Press enter to proceed. ")
		reader.ReadString('\n')
		fmt.Printf("\n")
		fmt.Printf("1. Add Student\n")
		fmt.Printf("2. Edit Student\n")
		fmt.Printf("3. Search Student\n")
		fmt.Printf("4. Delete Student\n")
		fmt.Printf("5. Print Students\n")
		fmt.Printf("0. Exit Program\n")

		action := inputInt("Enter action (1/2/3/4/5/0): ", 0, 5)
		index := -1
		switch action {
		case 0:
			fmt.Printf("Thank you for using our program!\n")
			return
		case 2, 3, 4:
			number := inputInt("Student Number: ", 0, 99999999)
			index = studentDatabase.searchStudent(number)
			if index == -1 {
				fmt.Printf("Student not found.\n")
				continue
			}
		}

		switch action {
		case 1:
			studentDatabase.addStudent(inputStudent())
			fmt.Printf("Student successfully added to the database.\n")
		case 2:
			studentDatabase.editStudent(index, inputStudent())
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
func inputStudent() Student {
	return Student{
		number:     inputInt("Student Number: ", 0, 99999999),
		name:       inputString("Name: "),
		department: inputString("Department: "),
	}
}

// Reusable scanner for helper input functions.
var reader = bufio.NewReader(os.Stdin)

// A simple helper function that gets the user string input after they press enter.
func inputString(prompt string) string {
	fmt.Printf("%s", prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Helper function to easily get user integer input with validation.
// When the user hits enter, the string input will be parsed into an integer.
// If it cannot be parsed into an integer or the integer is outside min and max, print input error.
// Otherwise, the parsed integer will be returned like normal.
func inputInt(prompt string, min, max int) int {
	for {
		fmt.Printf("%s", prompt)
		input, _ := reader.ReadString('\n')
		parsed, err := strconv.Atoi(strings.TrimSpace(input))

		if err != nil {
			fmt.Printf("INPUT ERROR. Please enter a valid integer.\n")
		} else if parsed < min || parsed > max {
			fmt.Printf("INPUT ERROR. Please enter an integer between %d and %d.\n", min, max)
		} else {
			return parsed
		}
	}
}
