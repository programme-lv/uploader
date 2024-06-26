//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type TaskSubmissions struct {
	UserID            int64
	TaskID            int64
	ProgrammingLangID string
	Submission        string
	CreatedAt         time.Time
	Hidden            bool
	VisibleEvalID     *int64
	ID                uuid.UUID `sql:"primary_key"`
}
