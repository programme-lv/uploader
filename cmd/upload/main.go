package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	var rootFolder string
	flag.StringVar(&rootFolder, "folder", "./", "Path to the folder to process")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	problem, err := readProblemToml(rootFolder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(problem)

	processTestsDir(rootFolder)

	key := os.Getenv("SPACES_KEY")
	secret := os.Getenv("SPACES_SECRET")

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String("https://fra1.digitaloceanspaces.com"),
		Region:           aws.String("fra1"),
		S3ForcePathStyle: aws.Bool(false),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatalf("Error creating new session: %v", err)
	}
	s3Client := s3.New(newSession)

	// Step 4: Define the parameters of the object you want to upload.
	object := s3.PutObjectInput{
		Bucket: aws.String("proglv-dev"),                  // The path to the directory you want to upload the object to, starting with your Space name.
		Key:    aws.String("folder-path/hello-world.txt"), // Object key, referenced whenever you want to access this file later.
		Body:   strings.NewReader("Hello, World!"),        // The object's contents.
		ACL:    aws.String("private"),                     // Defines Access-control List (ACL) permissions, such as private or public.
		Metadata: map[string]*string{ // Required. Defines metadata tags.
			"x-amz-meta-my-key": aws.String("your-value"),
		},
	}

	// Step 5: Run the PutObject function with your parameters, catching for errors.
	_, err = s3Client.PutObject(&object)
	if err != nil {
		fmt.Println(err.Error())
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

func processTestsDir(rootFolder string) error {
	testDir := filepath.Join(rootFolder, "tests")
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			log.Println("found test", path)
		}
		return nil
	})
	return err
}
