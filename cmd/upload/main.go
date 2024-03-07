package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
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

	accessKey := os.Getenv("SPACES_KEY")
	secretKey := os.Getenv("SPACES_SECRET")
	region := "fra1"
	endpoint := "https://fra1.digitaloceanspaces.com"
	bucket := os.Getenv("S3_BUCKET")

	s3Uploader := upload.NewS3Uploader(accessKey, secretKey,
		region, endpoint, bucket)

	problem, err := readProblemToml(rootFolder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(problem)

	err = upload.ProcessTestsDir(rootFolder, s3Uploader)
	if err != nil {
		log.Fatal(err)
	}
}

type Problem struct {
	Specification string   `toml:"specification"`
	Code          string   `toml:"code"`
	Name          string   `toml:"name"`
	Authors       []string `toml:"authors"`
	Tags          []string `toml:"tags"`
	Type          string   `toml:"type"`
	Time          float64  `toml:"time"`
	Memory        int      `toml:"memory"`
	Difficulty    int      `toml:"difficulty"`
}

func readProblemToml(rootFolder string) (Problem, error) {
	problemTomlPath := filepath.Join(rootFolder, "problem.toml")
	var problem Problem
	_, err := toml.DecodeFile(problemTomlPath, &problem)
	return problem, err
}
