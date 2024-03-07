package upload

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

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
