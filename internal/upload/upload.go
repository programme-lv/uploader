package upload

import (
	"log"
	"math"
	"os"
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

	// create a new task
	// create a new task version
	// fill task_version_tests, text_files
	// link relevant_version to task
	statementsDir := path.Join(taskDir, "statements")
	err = processStatementDir(versionID, statementsDir, sqlxDB)
	if err != nil {
		log.Fatalf("Error processing statement dir: %v", err)
	}

	testDir := filepath.Join(taskDir, "tests")
	err = processTestsDir(testDir, uploader, sqlxDB)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func assignCheckerToTaskVersion(db qrm.Executable, checkerID, versionID int) error {
	updateStmt := table.TaskVersions.UPDATE(
		table.TaskVersions.CheckerID,
	).SET(
		checkerID,
	).WHERE(
		table.TaskVersions.ID.EQ(postgres.Int(int64(versionID))),
	)

	_, err := updateStmt.Exec(db)
	return err
}

func readCheckerFile(taskDir string) string {
	checkerCpp, err := os.ReadFile(path.Join(taskDir, "evaluation", "checker.cpp"))
	if err != nil {
		log.Fatalf("Error reading checker file: %v", err)
	}
	return string(checkerCpp)
}

func ensureCheckerExists(db qrm.Queryable, code string) (int, error) {
	insertStmt := table.TestlibCheckers.INSERT(
		table.TestlibCheckers.Code,
	).VALUES(code).
		ON_CONFLICT(table.TestlibCheckers.Code).
		DO_UPDATE(postgres.SET(table.TestlibCheckers.Code.SET(postgres.String(code)))).
		RETURNING(table.TestlibCheckers.ID)

	dest := model.TestlibCheckers{}
	err := insertStmt.Query(db, &dest)
	if err != nil {
		return 0, err
	}

	return int(dest.ID), nil
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
