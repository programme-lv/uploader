package upload

import (
	"log"
	"math"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/uploader/internal/database/proglv/public/model"
	"github.com/programme-lv/uploader/internal/database/proglv/public/table"
)

func UploadTask(taskDir string, uploader S3Uploader, sqlxDB *sqlx.DB) error {

	problem, err := readProblemToml(taskDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(problem)

	taskID, err := createTask(sqlxDB)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("TaskID:", taskID)

	versionID, err := createTaskVersion(sqlxDB, taskID, problem)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("VersionID:", versionID)
	// create a new task
	// create a new task version
	// fill task_version_tests, text_files
	// link relevant_version to task

	err = ProcessTestsDir(taskDir, uploader, sqlxDB)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func createTaskVersion(db qrm.Queryable, taskID int, problem Problem) (int, error) {
	insertStmt := table.TaskVersions.INSERT(
		table.TaskVersions.TaskID,
		table.TaskVersions.ShortCode,
		table.TaskVersions.FullName,
		table.TaskVersions.TimeLimMs,
		table.TaskVersions.MemLimKibibytes,
		table.TaskVersions.TestingTypeID,
		// table.TaskVersions.Origin,
		table.TaskVersions.CreatedAt,
		// table.TaskVersions.CheckerID,
		// table.TaskVersions.InteractorID,
	).VALUES(
		taskID,
		problem.Code,
		problem.Name,
		math.Round(problem.Time*1000),
		problem.Memory*1024,
		problem.Type,
		postgres.NOW(),
	).RETURNING(table.TaskVersions.ID)

	dest := model.TaskVersions{}
	err := insertStmt.Query(db, &dest)
	if err != nil {
		return 0, err
	}

	return int(dest.ID), nil
}

func createTask(db qrm.Queryable) (int, error) {
	insertStmt := table.Tasks.INSERT(
		table.Tasks.CreatedAt,
		table.Tasks.CreatedByID,
		table.Tasks.RelevantVersionID,
		table.Tasks.PublishedVersionID,
	).VALUES(
		postgres.NOW(),
		1,
		nil,
		nil,
	).RETURNING(table.Tasks.ID)

	dest := model.Tasks{}
	err := insertStmt.Query(db, &dest)
	if err != nil {
		return 0, err
	}

	return int(dest.ID), nil

}
