package coverage

import (
	"fmt"
	"os"
	"strings"

	"golens-api/utils"

	"github.com/pkg/errors"
	"golang.org/x/tools/cover"
)

type ICoverage interface {
	GenerateCoverageAndHTMLFiles(path string) error
	GetCoveredLines(coverageName string) (int, int, error)
	IsGoDirectory(dirPath string) (bool, error)
	GetFileCoveragePercentage(coverageName string) (map[string][]map[string]any, error)
	GetCoveredLinesByPackage(coverageName string) (map[string]map[string]int, error)
}

type Coverage struct {
}

func NewCoverage() *Coverage {
	return &Coverage{}
}

func (u *Coverage) GenerateCoverageAndHTMLFiles(path string) error {
	err := utils.RunGitCommands(path)
	if err != nil {
		return errors.WithStack(err)
	}

	err = utils.GenerateCoverageOut(path)
	if err != nil {
		return errors.WithStack(err)
	}

	err = utils.GenerateCoverageHTML(path)
	if err != nil {
		return errors.WithStack(err)
	}

	err = utils.MoveCoverageAndHTMLFiles(path)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u *Coverage) GetCoveredLines(coverageName string) (int, int, error) {
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

func (u *Coverage) IsGoDirectory(dirPath string) (bool, error) {
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

func (u *Coverage) GetFileCoveragePercentage(coverageName string) (map[string][]map[string]any, error) {
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
		fileName := utils.GetProfileNameFromPath(profile.FileName)
		packageName := utils.GetPackageNameFromPath(profile.FileName)

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

func (u *Coverage) GetCoveredLinesByPackage(coverageName string) (map[string]map[string]int, error) {
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
		packageName := utils.GetPackageNameFromPath(profile.FileName)

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
