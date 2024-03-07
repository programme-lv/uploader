package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ProcessTestsDir(rootFolder string, uploader S3Uploader, sqlx any) error {
	testDir := filepath.Join(rootFolder, "tests")
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			sha256Hex, err := getFileSHA256Hex(path)
			if err != nil {
				return err
			}

			log.Println("uploading", path, sha256Hex)

			testPathDOSpaces := fmt.Sprintf("tests/%s", sha256Hex)
			exists, err := uploader.Exists(testPathDOSpaces)
			if err != nil {
				return err
			}
			if !exists {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				err = uploader.Upload(testPathDOSpaces, file)
				if err != nil {
					return err
				}
				log.Println("file uploaded")
			} else {
				log.Println("file already exists")
			}
		}
		return nil
	})
	return err
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
