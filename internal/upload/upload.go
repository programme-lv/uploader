package upload

import (
	"log"
	"math"
	"path"
	"path/filepath"

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
		log.Fatalf("Error creating task: %v", err)
	}

	log.Println("TaskID:", taskID)

	versionID, err := createTaskVersion(sqlxDB, taskID, problem)
	if err != nil {
		log.Fatalf("Error creating task version: %v", err)
	}

	log.Println("VersionID:", versionID)

	err = updateTaskRelevantVersionID(sqlxDB, taskID, versionID)
	if err != nil {
		log.Fatalf("Error updating task relevant version ID: %v", err)
	}

	checkerCpp := readCheckerFile(taskDir)
	checkerID, err := ensureCheckerExists(sqlxDB, checkerCpp)
	if err != nil {
		log.Fatalf("Error ensuring checker exists: %v", err)
	}
	log.Println("CheckerID:", checkerID)

	err = assignCheckerToTaskVersion(sqlxDB, checkerID, versionID)
	if err != nil {
		log.Fatalf("Error assigning checker to task version: %v", err)
	}
	log.Println("Checker assigned to task version")

	statementsDir := path.Join(taskDir, "statements")
	err = processStatementDir(versionID, statementsDir, sqlxDB)
	if err != nil {
		log.Fatalf("Error processing statement dir: %v", err)
	}

	testDir := filepath.Join(taskDir, "tests")
	err = processTestsDir(testDir, uploader, sqlxDB, versionID)
	if err != nil {
		log.Fatalf("Error processing tests dir: %v", err)
	}

	examplesDir := path.Join(taskDir, "examples")
	err = processExamplesDir(versionID, examplesDir, sqlxDB)
	if err != nil {
		log.Fatalf("Error processing examples dir: %v", err)
	}

	return nil
}

func updateTaskRelevantVersionID(db qrm.Executable, taskID, versionID int) error {
	updateStmt := table.Tasks.UPDATE(
		table.Tasks.RelevantVersionID,
	).SET(
		versionID,
	).WHERE(table.Tasks.ID.EQ(postgres.Int64(int64(taskID))))

	_, err := updateStmt.Exec(db)
	return err
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
