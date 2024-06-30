package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	set "github.com/deckarep/golang-set/v2"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/uploader/internal/database/proglv/public/model"
	"github.com/programme-lv/uploader/internal/database/proglv/public/table"
)

func processTestsDir(testDir string, uploader S3Uploader, db qrm.DB) (int64, error) {
	tests, err := getTestNamesNoExt(testDir)
	if err != nil {
		return 0, err
	}

	// create new test set
	createTestSetStmt := table.TestSets.INSERT(
		table.TestSets.CreatedAt,
	).VALUES(postgres.NOW()).RETURNING(table.TestSets.ID)

	var testSetRecord model.TestSets
	err = createTestSetStmt.Query(db, &testSetRecord)
	if err != nil {
		// failed to create test set
		log.Println("Failed to create test set")
		return 0, err
	}

	testSetID := testSetRecord.ID

	for _, test := range tests.ToSlice() {
		inPath := filepath.Join(testDir, test+".in")
		ansPath := filepath.Join(testDir, test+".ans")

		inID, err := ensureTextFileExistsDBAndS3(inPath, uploader, db)
		if err != nil {
			return 0, err
		}

		ansID, err := ensureTextFileExistsDBAndS3(ansPath, uploader, db)
		if err != nil {
			return 0, err
		}

		err = linkTaskVersionToTest(db, test, testSetID, inID, ansID)
		if err != nil {
			return 0, err
		}

	}
	return testSetID, nil
}

func linkTaskVersionToTest(db qrm.Executable, name string, testSetID, inTextFileID, ansTextFileID int64) error {
	log.Printf("Linking test %v to task version %v", name, testSetID)
	insertStmt := table.TestSetTests.INSERT(
		table.TestSetTests.TestFilename,
		table.TestSetTests.TestSetID,
		table.TestSetTests.InputTextFileID,
		table.TestSetTests.AnswerTextFileID,
	).VALUES(
		name,
		testSetID,
		inTextFileID,
		ansTextFileID,
	)

	_, err := insertStmt.Exec(db)
	return err

}

func ensureTextFileExistsDBAndS3(path string, uploader S3Uploader, db qrm.DB) (int64, error) {
	sha256Hex, err := getFileSHA256Hex(path)
	if err != nil {
		return 0, err
	}

	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	err = ensureTextFileExistsS3(uploader, sha256Hex, file)
	if err != nil {
		return 0, err
	}
	file.Close()

	textFileId, err := ensureTextFileRowExistsDB(db, sha256Hex)
	if err != nil {
		return 0, err
	}

	return textFileId, nil
}

func getTestNamesNoExt(testDir string) (set.Set[string], error) {
	testNames := set.NewSet[string]()
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			testNames.Add(filenameWithoutExt(info.Name()))
		}
		return nil
	})
	return testNames, err

}

func ensureTextFileExistsS3(uploader S3Uploader, sha256Hex string, content io.ReadSeeker) error {
	testS3Key := fmt.Sprintf("tests/%s", sha256Hex)

	log.Printf("Ensuring %v exists in S3", testS3Key)

	exists, err := uploader.Exists(testS3Key)
	if err != nil {
		return err
	}
	if !exists {
		err = uploader.Upload(testS3Key, content)
		if err != nil {
			return err
		}
		log.Printf("Uploaded %v to S3", testS3Key)
	} else {
		log.Printf("%v already exists in S3", testS3Key)
	}

	return nil
}

func ensureTextFileRowExistsDB(db qrm.Queryable, sha256Hex string) (int64, error) {
	log.Printf("Ensuring %v text file exists in DB", sha256Hex)

	insertStmt := table.TextFiles.INSERT(
		table.TextFiles.Sha256,
	).VALUES(sha256Hex).
		ON_CONFLICT(table.TextFiles.Sha256).
		DO_UPDATE(
			postgres.SET(table.TextFiles.Sha256.SET(postgres.String(sha256Hex)))).
		RETURNING(table.TextFiles.ID)

	dest := model.TextFiles{}
	err := insertStmt.Query(db, &dest)
	if err != nil {
		return 0, err
	}

	return dest.ID, nil
}

func getFileSHA256Hex(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}
