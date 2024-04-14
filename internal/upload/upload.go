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

	err = updateTaskCurrentVersionID(sqlxDB, taskID, versionID)
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
	statementIDs, err := processStatementDir(statementsDir, sqlxDB)
	if err != nil {
		log.Fatalf("Error processing statement dir: %v", err)
	}

	if len(statementIDs) != 1 {
		log.Fatalf("Expected 1 statement, got %v", len(statementIDs))
	}

	statementID := statementIDs[0]
	err = assignMarkdownStatementToTaskVersion(sqlxDB, statementID, versionID)
	if err != nil {
		log.Fatalf("Error assigning markdown statement to task version: %v", err)
	}

	testDir := filepath.Join(taskDir, "tests")
	testSetID, err := processTestsDir(testDir, uploader, sqlxDB)
	if err != nil {
		log.Fatalf("Error processing tests dir: %v", err)
	}

	err = assignTestSetToTaskVersion(sqlxDB, testSetID, versionID)
	if err != nil {
		log.Fatalf("Error assigning test set to task version: %v", err)
	}

	examplesDir := path.Join(taskDir, "examples")
	exampleSetID, err := processExamplesDir(examplesDir, sqlxDB)
	if err != nil {
		log.Fatalf("Error processing examples dir: %v", err)
	}

	// assign example set to task version
	err = assignExampleSetToTaskVersion(sqlxDB, exampleSetID, versionID)
	if err != nil {
		log.Fatalf("Error assigning example set to task version: %v", err)
	}

	return nil
}

func updateTaskCurrentVersionID(db qrm.Executable, taskID, versionID int64) error {
	updateStmt := table.Tasks.UPDATE(
		table.Tasks.CurrentVersionID,
	).SET(
		versionID,
	).WHERE(table.Tasks.ID.EQ(postgres.Int64(taskID)))

	_, err := updateStmt.Exec(db)
	return err
}

func assignExampleSetToTaskVersion(db qrm.Executable, exampleSetID, versionID int64) error {
	updateStmt := table.TaskVersions.UPDATE(
		table.TaskVersions.ExampleSetID,
	).SET(
		exampleSetID,
	).WHERE(table.TaskVersions.ID.EQ(postgres.Int64(int64(versionID))))

	_, err := updateStmt.Exec(db)
	return err
}

func assignTestSetToTaskVersion(db qrm.Executable, testSetID, versionID int64) error {
	updateStmt := table.TaskVersions.UPDATE(
		table.TaskVersions.TestSetID,
	).SET(
		testSetID,
	).WHERE(table.TaskVersions.ID.EQ(postgres.Int64(int64(versionID))))

	_, err := updateStmt.Exec(db)
	return err
}

func assignMarkdownStatementToTaskVersion(db qrm.Executable, statementID, versionID int64) error {
	updateStmt := table.TaskVersions.UPDATE(
		table.TaskVersions.MdStatementID,
	).SET(
		statementID,
	).WHERE(table.TaskVersions.ID.EQ(postgres.Int64(int64(versionID))))

	_, err := updateStmt.Exec(db)
	return err
}

func createTaskVersion(db qrm.Queryable, taskID int64, problem Problem) (int64, error) {
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

	return dest.ID, nil
}

func createTask(db qrm.Queryable) (int64, error) {
	insertStmt := table.Tasks.INSERT(
		table.Tasks.CreatedAt,
		table.Tasks.CreatedByID,
		table.Tasks.CurrentVersionID,
		table.Tasks.StableVersionID,
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

	return dest.ID, nil

}
