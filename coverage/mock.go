package coverage

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CoverageMock struct {
}

func NewCoverageMock() *CoverageMock {
	return &CoverageMock{}
}

func (u *CoverageMock) GenerateCoverageAndHTMLFiles(path string) error {
	return nil
}

func (u *CoverageMock) GetCoveredLines(coverageName string, ignoredPackages map[string]bool) (int, int, error) {
	return 1000, 1000, nil
}

func (u *CoverageMock) IsGoDirectory(dirPath string) (bool, error) {
	return true, nil
}

func (u *CoverageMock) GetFileCoveragePercentage(
	coverageName string,
	ignoredFilesByPackage map[string]map[string]bool,
) (map[string][]map[string]any, error) {
	return nil, nil
}

func (u *CoverageMock) GetCoveredLinesByPackage(
	coverageName string,
	ignoredFilesByPackage map[string]map[string]bool,
	ignoredPackages map[string]bool,
) (map[string]map[string]int, error) {
	return nil, nil
}

func (u *CoverageMock) GetIgnoredFilesByPackage(
	ctx *gin.Context,
	db *gorm.DB,
	directoryName string,
) map[string]map[string]bool {
	return nil
}

func (u *CoverageMock) GetIgnoredPackages(ctx *gin.Context, db *gorm.DB, directoryName string) map[string]bool {
	return nil
}
