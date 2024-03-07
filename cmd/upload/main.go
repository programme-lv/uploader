package main

import (
	"flag"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/programme-lv/uploader/internal/upload"
)

func main() {
	var rootFolder string
	flag.StringVar(&rootFolder, "folder", "./", "Path to the folder to process")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	sqlxConnString := os.Getenv("SQLX_CONN_STRING")
	accessKey := os.Getenv("SPACES_KEY")
	secretKey := os.Getenv("SPACES_SECRET")
	bucket := os.Getenv("S3_BUCKET")

	sqlxDB := sqlx.MustConnect("postgres", sqlxConnString)
	defer sqlxDB.Close()

	s3Uploader := upload.NewS3Uploader(accessKey, secretKey,
		"fra1", "https://fra1.digitaloceanspaces.com", bucket)

	upload.UploadTask(rootFolder, s3Uploader, sqlxDB)
}
