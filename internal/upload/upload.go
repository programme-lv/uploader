package upload

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/uploader/internal/database/proglv/public/table"
)

func UploadTask(taskDir string, uploader S3Uploader, sqlxDB *sqlx.DB) error {

	problem, err := readProblemToml(taskDir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(problem)

	taskID, err := createANewTask(problem, sqlxDB)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(taskID)

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

func createANewTask(problem Problem, sqlxDB *sqlx.DB) (int, error) {
	insertStmt := table.Tasks.INSERT(
		table.Tasks.CreatedAt,
		table.Tasks.CreatedByID,
		table.Tasks.RelevantVersionID,
		table.Tasks.PublishedVersionID,
	).VALUES(
		"now()",
		1,
		nil,
		nil,
	).RETURNING(table.Tasks.ID)

	dest := struct {
		ID int `db:"id"`
	}{}

	err := insertStmt.Query(sqlxDB.DB, &dest)
	if err != nil {
		return 0, err
	}

	return dest.ID, nil

}
