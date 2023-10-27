package utils

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/cover"
)

var (
	// mocks
	IsGoDirectoryF                = IsGoDirectory
	GetCoveredLinesF              = GetCoveredLines
	GenerateCoverageAndHTMLFilesF = GenerateCoverageAndHTMLFiles
	GetWorkingDirectoryF          = GetWorkingDirectory
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

func CalculateCoverage(totalLines int, coveredLines int) float64 {
	coveragePercentage := float64(coveredLines) / float64(totalLines) * 100

	if math.IsNaN(coveragePercentage) {
		coveragePercentage = 0
	} else {
		coveragePercentage = math.Round(coveragePercentage*100) / 100
	}

	return coveragePercentage
}

func GetCoveredLinesByPackage(coverageName string) (map[string]map[string]int, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	coverageProfile := fmt.Sprintf("%s/data/coverage/%s.out", workingDir, coverageName)

	profiles, err := cover.ParseProfiles(coverageProfile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	coveredLinesByPackage := make(map[string]map[string]int)

	for _, profile := range profiles {
		packageName := GetPackageNameFromPath(profile.FileName)

		if _, ok := coveredLinesByPackage[packageName]; !ok {
			coveredLinesByPackage[packageName] = map[string]int{
				"totalLines":   0,
				"coveredLines": 0,
			}
		}

		for _, block := range profile.Blocks {
			coveredLinesByPackage[packageName]["totalLines"] += block.NumStmt

			if block.Count > 0 {
				coveredLinesByPackage[packageName]["coveredLines"] += block.NumStmt
			}
		}
	}

	return coveredLinesByPackage, nil
}

func GetFileCoveragePercentage(coverageName string) (map[string][]map[string]any, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	coverageProfile := fmt.Sprintf("%s/data/coverage/%s.out", workingDir, coverageName)

	profiles, err := cover.ParseProfiles(coverageProfile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fileMap := make(map[string][]map[string]any)

	for _, profile := range profiles {
		totalLines := 0
		coveredLines := 0
		fileName := GetProfileNameFromPath(profile.FileName)
		packageName := GetPackageNameFromPath(profile.FileName)

		for _, block := range profile.Blocks {
			totalLines += block.NumStmt

			if block.Count > 0 {
				coveredLines += block.NumStmt
			}
		}

		coverageMap := map[string]any{
			"fileName":     fileName,
			"totalLines":   totalLines,
			"coveredLines": coveredLines,
		}

		if _, ok := fileMap[packageName]; !ok {
			fileMap[packageName] = []map[string]any{coverageMap}
		} else {
			fileMap[packageName] = append(fileMap[packageName], coverageMap)
		}

	}

	return fileMap, nil
}

func GetCoveredLines(coverageName string) (int, int, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}
	coverageProfile := fmt.Sprintf("%s/data/coverage/%s.out", workingDir, coverageName)

	profiles, err := cover.ParseProfiles(coverageProfile)
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}

	totalLines := 0
	coveredLines := 0
	for _, profile := range profiles {
		profileTotalStatements := 0
		profileCoveredStatements := 0

		for _, block := range profile.Blocks {
			profileTotalStatements += block.NumStmt
			totalLines += block.NumStmt
			if block.Count > 0 {
				coveredLines += block.NumStmt
				profileCoveredStatements += block.NumStmt
			}
		}
	}

	return totalLines, coveredLines, nil
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

func IsGoDirectory(dirPath string) (bool, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			return true, nil
		}
	}

	return false, nil
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

func RemoveFile() {}
