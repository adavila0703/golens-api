package coverage

type CoverageMock struct {
}

func NewCoverageMock() *CoverageMock {
	return &CoverageMock{}
}

func (u *CoverageMock) GenerateCoverageAndHTMLFiles(path string) error {
	return nil
}

func (u *CoverageMock) GetCoveredLines(coverageName string) (int, int, error) {
	return 1000, 1000, nil
}

func (u *CoverageMock) IsGoDirectory(dirPath string) (bool, error) {
	return true, nil
}

func (u *CoverageMock) GetFileCoveragePercentage(coverageName string) (map[string][]map[string]any, error) {
	return nil, nil
}

func (u *CoverageMock) GetCoveredLinesByPackage(coverageName string) (map[string]map[string]int, error) {
	return nil, nil
}
