//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var TestSets = newTestSetsTable("public", "test_sets", "")

type testSetsTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	CreatedAt postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type TestSetsTable struct {
	testSetsTable

	EXCLUDED testSetsTable
}

// AS creates new TestSetsTable with assigned alias
func (a TestSetsTable) AS(alias string) *TestSetsTable {
	return newTestSetsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TestSetsTable with assigned schema name
func (a TestSetsTable) FromSchema(schemaName string) *TestSetsTable {
	return newTestSetsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TestSetsTable with assigned table prefix
func (a TestSetsTable) WithPrefix(prefix string) *TestSetsTable {
	return newTestSetsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TestSetsTable with assigned table suffix
func (a TestSetsTable) WithSuffix(suffix string) *TestSetsTable {
	return newTestSetsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTestSetsTable(schemaName, tableName, alias string) *TestSetsTable {
	return &TestSetsTable{
		testSetsTable: newTestSetsTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newTestSetsTableImpl("", "excluded", ""),
	}
}

func newTestSetsTableImpl(schemaName, tableName, alias string) testSetsTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		CreatedAtColumn = postgres.TimestampColumn("created_at")
		allColumns      = postgres.ColumnList{IDColumn, CreatedAtColumn}
		mutableColumns  = postgres.ColumnList{CreatedAtColumn}
	)

	return testSetsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		CreatedAt: CreatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
