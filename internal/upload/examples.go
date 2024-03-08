package upload

import (
	"os"
	"path"
	"strings"

	set "github.com/deckarep/golang-set/v2"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/uploader/internal/database/proglv/public/table"
)

func processExamplesDir(versionID int, examplesDir string, db qrm.Executable) error {
	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		return err
	}

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
			return err
		}
		inStr := string(inBytes)

		ansPath := path.Join(examplesDir, name+".ans")
		ansBytes, err := os.ReadFile(ansPath)
		if err != nil {
			return err
		}
		ansStr := string(ansBytes)

		err = createExample(versionID, inStr, ansStr, db)
		if err != nil {
			return err
		}

	}
	return nil
}

func createExample(versionID int, in, ans string, db qrm.Executable) error {
	insertStmt := table.StatementExamples.INSERT(
		table.StatementExamples.TaskVersionID,
		table.StatementExamples.Input,
		table.StatementExamples.Answer,
	).VALUES(
		versionID,
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
