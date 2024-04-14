package upload

import (
	"log"
	"os"
	"path"
	"strings"

	set "github.com/deckarep/golang-set/v2"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/uploader/internal/database/proglv/public/model"
	"github.com/programme-lv/uploader/internal/database/proglv/public/table"
)

// ProcessExamplesDir processes the examples directory and creates examples set
// returns example_set_id
func processExamplesDir(examplesDir string, db qrm.DB) (int64, error) {
	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		return 0, err
	}

	// create new example set
	createExampleSetStmt := table.ExampleSets.INSERT(
		table.ExampleSets.CreatedAt,
	).VALUES(postgres.NOW()).RETURNING(table.ExampleSets.ID)

	var exampleSetRecord model.ExampleSets
	err = createExampleSetStmt.Query(db, &exampleSetRecord)
	if err != nil {
		return 0, err
	}

	exampleSetID := exampleSetRecord.ID

	exampleNames := set.NewSet[string]()
	for _, entry := range entries {
		if !entry.IsDir() {
			baseName := entry.Name()
			nameNoExt := filenameWithoutExt(baseName)
			exampleNames.Add(nameNoExt)
		}
	}

	for _, name := range exampleNames.ToSlice() {
		inPath := path.Join(examplesDir, name+".in")
		inBytes, err := os.ReadFile(inPath)
		if err != nil {
			return 0, err
		}
		inStr := string(inBytes)

		ansPath := path.Join(examplesDir, name+".ans")
		ansBytes, err := os.ReadFile(ansPath)
		if err != nil {
			return 0, err
		}
		ansStr := string(ansBytes)

		err = createExample(exampleSetID, inStr, ansStr, db)
		if err != nil {
			return 0, err
		}

	}
	return exampleSetID, nil
}

func createExample(exampleSetID int64, in, ans string, db qrm.Executable) error {
	log.Printf("Creating example for task version %v", exampleSetID)
	insertStmt := table.StatementExamples.INSERT(
		table.StatementExamples.ExampleSetID,
		table.StatementExamples.Input,
		table.StatementExamples.Answer,
	).VALUES(
		exampleSetID,
		in,
		ans,
	)

	_, err := insertStmt.Exec(db)
	return err
}

func filenameWithoutExt(filename string) string {
	if pos := strings.LastIndexByte(filename, '.'); pos != -1 {
		return filename[:pos]
	}
	return filename
}
