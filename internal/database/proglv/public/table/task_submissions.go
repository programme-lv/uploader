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

var TaskSubmissions = newTaskSubmissionsTable("public", "task_submissions", "")

type taskSubmissionsTable struct {
	postgres.Table

	// Columns
	UserID            postgres.ColumnInteger
	TaskID            postgres.ColumnInteger
	ProgrammingLangID postgres.ColumnString
	Submission        postgres.ColumnString
	CreatedAt         postgres.ColumnTimestampz
	Hidden            postgres.ColumnBool
	VisibleEvalID     postgres.ColumnInteger
	ID                postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type TaskSubmissionsTable struct {
	taskSubmissionsTable

	EXCLUDED taskSubmissionsTable
}

// AS creates new TaskSubmissionsTable with assigned alias
func (a TaskSubmissionsTable) AS(alias string) *TaskSubmissionsTable {
	return newTaskSubmissionsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TaskSubmissionsTable with assigned schema name
func (a TaskSubmissionsTable) FromSchema(schemaName string) *TaskSubmissionsTable {
	return newTaskSubmissionsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TaskSubmissionsTable with assigned table prefix
func (a TaskSubmissionsTable) WithPrefix(prefix string) *TaskSubmissionsTable {
	return newTaskSubmissionsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TaskSubmissionsTable with assigned table suffix
func (a TaskSubmissionsTable) WithSuffix(suffix string) *TaskSubmissionsTable {
	return newTaskSubmissionsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTaskSubmissionsTable(schemaName, tableName, alias string) *TaskSubmissionsTable {
	return &TaskSubmissionsTable{
		taskSubmissionsTable: newTaskSubmissionsTableImpl(schemaName, tableName, alias),
		EXCLUDED:             newTaskSubmissionsTableImpl("", "excluded", ""),
	}
}

func newTaskSubmissionsTableImpl(schemaName, tableName, alias string) taskSubmissionsTable {
	var (
		UserIDColumn            = postgres.IntegerColumn("user_id")
		TaskIDColumn            = postgres.IntegerColumn("task_id")
		ProgrammingLangIDColumn = postgres.StringColumn("programming_lang_id")
		SubmissionColumn        = postgres.StringColumn("submission")
		CreatedAtColumn         = postgres.TimestampzColumn("created_at")
		HiddenColumn            = postgres.BoolColumn("hidden")
		VisibleEvalIDColumn     = postgres.IntegerColumn("visible_eval_id")
		IDColumn                = postgres.StringColumn("id")
		allColumns              = postgres.ColumnList{UserIDColumn, TaskIDColumn, ProgrammingLangIDColumn, SubmissionColumn, CreatedAtColumn, HiddenColumn, VisibleEvalIDColumn, IDColumn}
		mutableColumns          = postgres.ColumnList{UserIDColumn, TaskIDColumn, ProgrammingLangIDColumn, SubmissionColumn, CreatedAtColumn, HiddenColumn, VisibleEvalIDColumn}
	)

	return taskSubmissionsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID:            UserIDColumn,
		TaskID:            TaskIDColumn,
		ProgrammingLangID: ProgrammingLangIDColumn,
		Submission:        SubmissionColumn,
		CreatedAt:         CreatedAtColumn,
		Hidden:            HiddenColumn,
		VisibleEvalID:     VisibleEvalIDColumn,
		ID:                IDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
