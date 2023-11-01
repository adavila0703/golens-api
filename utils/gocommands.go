package utils

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

var (
	// mocks
	GetWorkingDirectoryF = GetWorkingDirectory
	RemoveFileF          = RemoveFile
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

func RunGitCommands(path string) error {
	hasGitDir := hasGitDirectory(path)

	if hasGitDir {
		cmd := exec.Command("git", "checkout", "main")
		cmd.Dir = path
		cmd.Run()

		cmd = exec.Command("git", "checkout", "master")
		cmd.Dir = path
		cmd.Run()

		cmd = exec.Command("git", "fetch", "&&", "git", "pull")
		cmd.Dir = path
		cmd.Run()
	}

	return nil
}

func GenerateCoverageOut(path string) error {
	coverageName := GetCoverageNameFromPath(path)

	coverageProfile := fmt.Sprintf("%s/%s.out", path, coverageName)

	cmd := exec.Command("go", "test", "./...", "-coverprofile="+coverageProfile)
	cmd.Dir = path
	cmd.Run()

	return nil
}

func GenerateCoverageHTML(path string) error {
	coverageName := GetCoverageNameFromPath(path)

	coverageProfile := fmt.Sprintf("%s/%s.out", path, coverageName)

	destinationPath := fmt.Sprintf("%s/%s.html", path, coverageName)

	cmd := exec.Command("go", "tool", "cover", "-html="+coverageProfile, "-o", destinationPath)
	cmd.Dir = path
	cmd.Run()

	return nil
}

func CalculateCoverage(totalLines int, coveredLines int) float64 {
	coveragePercentage := float64(coveredLines) / float64(totalLines) * 100

	if math.IsNaN(coveragePercentage) {
		coveragePercentage = 0
	} else {
		coveragePercentage = math.Round(coveragePercentage*100) / 100
	}

	return coveragePercentage
}

func GetCoverageNameFromPath(path string) string {
	var pathStrings []string

	if runtime.GOOS == "windows" {
		pathStrings = strings.Split(path, "\\")
	} else {
		pathStrings = strings.Split(path, "/")
	}

	return pathStrings[len(pathStrings)-1]
}

func GetProfileNameFromPath(path string) string {
	pathStrings := strings.Split(path, "/")
	return pathStrings[len(pathStrings)-1]
}

func GetPackageNameFromPath(path string) string {
	packageString := strings.Split(path, "/")
	return packageString[len(packageString)-2]
}

func hasGitDirectory(path string) bool {
	gitPath := filepath.Join(path, ".git")
	_, err := os.Stat(gitPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetWorkingDirectory() (string, error) {
	return os.Getwd()
}

func RemoveFile(file string) error {
	return os.Remove(file)
}
