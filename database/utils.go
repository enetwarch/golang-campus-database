// Package database contains functions that handles sqlite queries.
package database

import "database/sql"

func (table *Table) ColumnNames() []string {
	columnNames := make([]string, len(table.Columns))
	for i := range table.Columns {
		columnNames[i] = table.Columns[i].ColumnName
	}
	return columnNames
}

func (table *Table) PKColumnNames() []string {
	pkColumnNames := make([]string, len(table.PrimaryKeys))
	for i := range table.PrimaryKeys {
		pkColumnNames[i] = table.PrimaryKeys[i].ColumnName
	}
	return pkColumnNames
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
