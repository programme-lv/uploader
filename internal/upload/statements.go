package upload

import (
	"log"
	"os"
	"path"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/uploader/internal/database/proglv/public/model"
	"github.com/programme-lv/uploader/internal/database/proglv/public/table"
)

func processStatementDir(statementsDir string, db qrm.DB) ([]int64, error) {
	// list folders in taskDir
	entries, err := os.ReadDir(statementsDir)
	if err != nil {
		return nil, err
	}

	res := make([]int64, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "md" {
				x, err := processMdDir(path.Join(statementsDir, "md"), db)
				if err != nil {
					return nil, err
				}
				res = append(res, x...)
			}
		}
	}
	return res, nil
}

func processMdDir(mdDir string, db qrm.DB) ([]int64, error) {
	entries, err := os.ReadDir(mdDir)
	if err != nil {
		return nil, err
	}

	res := make([]int64, 0)

	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := path.Join(mdDir, entry.Name())
			stmtID, err := processMdLangDir(dirPath, entry.Name(), db)
			if err != nil {
				return nil, err
			}
			res = append(res, stmtID)
		}
	}
	return res, nil
}

func processMdLangDir(mdLangDir string, lang string, db qrm.DB) (int64, error) {
	log.Printf("Processing lang dir: %v", mdLangDir)

	entries, err := os.ReadDir(mdLangDir)
	if err != nil {
		return 0, err
	}

	// create a new markdown statement
	var story, input, output string
	var notes, scoring *string

	for _, entry := range entries {
		if !entry.IsDir() {
			path := path.Join(mdLangDir, entry.Name())
			content, err := readMDFile(path)
			if err != nil {
				return 0, err
			}
			switch entry.Name() {
			case "story.md":
				story = content
			case "input.md":
				input = content
			case "output.md":
				output = content
			case "notes.md":
				notes = &content
			case "scoring.md":
				scoring = &content
			}
		}
	}

	insertStmt := table.MarkdownStatements.INSERT(
		table.MarkdownStatements.Story,
		table.MarkdownStatements.Input,
		table.MarkdownStatements.Output,
		table.MarkdownStatements.Notes,
		table.MarkdownStatements.Scoring,
		table.MarkdownStatements.LangIso6391,
	).VALUES(
		story,
		input,
		output,
		notes,
		scoring,
		lang,
	).RETURNING(table.MarkdownStatements.ID)

	var stmtRecord model.MarkdownStatements
	err = insertStmt.Query(db, &stmtRecord)
	if err != nil {
		return 0, err
	}

	return stmtRecord.ID, nil
}

func readMDFile(mdFilePath string) (string, error) {
	bytes, err := os.ReadFile(mdFilePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
