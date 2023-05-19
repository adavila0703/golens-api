package utils

import (
	"fmt"
	"golens-api/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/cover"
)

//	go tool cover -html=coverage.out

func MoveFile(path string, extension string, directory string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return errors.WithStack(err)
	}

	coverageName := GetCoverageNameFromPath(path)

	newPath := fmt.Sprintf("%s/data/%s/%s", workingDir, directory, coverageName+extension)

	oldPath := fmt.Sprintf("%s/%s%s", path, coverageName, extension)

	err = os.Rename(oldPath, newPath)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func MoveCoverageAndHTMLFiles(path string) error {
	err := MoveFile(path, ".out", "coverage")
	if err != nil {
		return errors.WithStack(err)
	}

	err = MoveFile(path, ".html", "html")
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GenerateCoverageAndHTMLFiles(path string) error {
	err := runGitCommands(path)
	if err != nil {
		return errors.WithStack(err)
	}

	err = GenerateCoverageOut(path)
	if err != nil {
		return errors.WithStack(err)
	}

	err = GenerateCoverageHTML(path)
	if err != nil {
		return errors.WithStack(err)
	}

	err = MoveCoverageAndHTMLFiles(path)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func runGitCommands(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			cmd := exec.Command("git", "checkout", config.Cfg.HeadBranch)
			cmd.Dir = path
			err = cmd.Run()
			if err != nil {
				return errors.WithStack(err)
			}

			cmd = exec.Command("git", "fetch", "&&", "git", "pull")
			cmd.Dir = path
			err = cmd.Run()
			if err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GenerateCoverageOut(path string) error {
	coverageName := GetCoverageNameFromPath(path)

	coverageProfile := fmt.Sprintf("%s/%s.out", path, coverageName)

	cmd := exec.Command("go", "test", "./...", "-coverprofile="+coverageProfile)
	cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GenerateCoverageHTML(path string) error {
	coverageName := GetCoverageNameFromPath(path)

	coverageProfile := fmt.Sprintf("%s/%s.out", path, coverageName)

	destinationPath := fmt.Sprintf("%s/%s.html", path, coverageName)

	cmd := exec.Command("go", "tool", "cover", "-html="+coverageProfile, "-o", destinationPath)
	cmd.Dir = path
	err := cmd.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func ParseCoverageOut() {
	coverageProfile := "./coverage.out"

	// Parse the coverage profile
	profiles, err := cover.ParseProfiles(coverageProfile)
	if err != nil {
		log.Fatal(err)
	}
	// Access coverage information from profiles
	for _, profile := range profiles {
		for _, block := range profile.Blocks {
			fmt.Println(block.EndLine)
		}
	}
}

func ParseCoveragePercentage(coverageName string) (float64, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	coverageProfile := fmt.Sprintf("%s/data/coverage/%s.out", workingDir, coverageName)

	// Parse the coverage profile
	profiles, err := cover.ParseProfiles(coverageProfile)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the coverage percentage
	totalStatements := 0
	coveredStatements := 0
	for _, profile := range profiles {
		for _, block := range profile.Blocks {
			totalStatements += block.NumStmt
			if block.Count > 0 {
				coveredStatements += block.NumStmt
			}
		}
	}

	coveragePercentage := float64(coveredStatements) / float64(totalStatements) * 100

	return coveragePercentage, nil
}

func GetCoverageNameFromPath(path string) string {
	pathStrings := strings.Split(path, "\\")
	return pathStrings[len(pathStrings)-1]
}
