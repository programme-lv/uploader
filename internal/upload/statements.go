package upload

import (
	"os"
	"path"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/uploader/internal/database/proglv/public/table"
)

func processStatementDir(versionID int, statementsDir string, db qrm.Executable) error {
	// list folders in taskDir
	entries, err := os.ReadDir(statementsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "md" {
				err := processMdDir(versionID, path.Join(statementsDir, "md"), db)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func processMdDir(versionID int, mdDir string, db qrm.Executable) error {
	entries, err := os.ReadDir(mdDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := path.Join(mdDir, entry.Name())
			processMdLangDir(versionID, dirPath, entry.Name(), db)
		}
	}
	return nil
}

func processMdLangDir(versionID int,
	mdLangDir string, lang string, db qrm.Executable) error {
	entries, err := os.ReadDir(mdLangDir)
	if err != nil {
		return err
	}

	// create a new markdown statement
	var story, input, output string
	var notes, scoring *string

	for _, entry := range entries {
		if !entry.IsDir() {
			path := path.Join(mdLangDir, entry.Name())
			content, err := readMDFile(path)
			if err != nil {
				return err
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
		table.MarkdownStatements.TaskVersionID,
		table.MarkdownStatements.LangIso6391,
	).VALUES(
		story,
		input,
		output,
		notes,
		scoring,
		versionID,
		lang,
	)

	insertStmt.Exec(db)

	return nil
}

func readMDFile(mdFilePath string) (string, error) {
	bytes, err := os.ReadFile(mdFilePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
