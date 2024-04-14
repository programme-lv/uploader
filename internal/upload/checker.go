package upload

import (
	"log"
	"os"
	"path"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/uploader/internal/database/proglv/public/model"
	"github.com/programme-lv/uploader/internal/database/proglv/public/table"
)

func assignCheckerToTaskVersion(db qrm.Executable, checkerID, versionID int64) error {
	updateStmt := table.TaskVersions.UPDATE(
		table.TaskVersions.CheckerID,
	).SET(
		checkerID,
	).WHERE(
		table.TaskVersions.ID.EQ(postgres.Int(versionID)),
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

func ensureCheckerExists(db qrm.Queryable, code string) (int64, error) {
	log.Printf("Ensuring checker exists")
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

	return dest.ID, nil
}
