package directory_test

import (
	"golens-api/api/directory"
	"golens-api/clients"
	"golens-api/models"
	"net/http/httptest"
	"regexp"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"golens-api/coverage"

	"github.com/DATA-DOG/go-sqlmock"
)

type GetRootDirectoryPathsCoverage struct {
	coverage.ICoverage
}

func NewGetRootDirectoryPathsCoverage() *GetRootDirectoryPathsCoverage {
	return &GetRootDirectoryPathsCoverage{}
}

func (g *GetRootDirectoryPathsCoverage) IsGoDirectory(dirPath string) (bool, error) {
	if dirPath == "happy" {
		return false, nil
	}

	return true, nil
}

var _ = Describe("GetRootDirectoryPaths", Ordered, func() {
	var mockClients *clients.GlobalClients
	var sqlMock sqlmock.Sqlmock
	var closeDB func() error
	var err error
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	BeforeAll(func() {
		var db *gorm.DB
		db, sqlMock, closeDB, err = clients.NewPostgresClientMock()
		utilsMock := NewGetRootDirectoryPathsCoverage()
		mockClients = clients.NewGlobalClients(db, nil, utilsMock)
	})

	It("checks for errors on creating mock client", func() {
		Expect(err).To(BeNil())
	})

	It("returns all none ignored go paths", func() {
		req := &directory.GetRootDirectoryPathsRequest{
			RootPath: "happy",
		}

		directory.GetDirPathsF = func(rootPath string) ([]string, error) {
			Expect(rootPath).To(Equal(req.RootPath))
			return []string{"dir1", "dir1", "test"}, nil
		}

		sqlMock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * 
			FROM "ignoreds" 
			WHERE type = $1 
			AND "ignoreds"."deleted_at" IS NULL
		`)).WithArgs(
			models.DirectoryType,
		).
			WillReturnRows(
				sqlmock.NewRows([]string{"name"}).AddRow("test"),
			)

		res, err := directory.GetRootDirectoryPaths(mockContext, req, mockClients)

		Expect(err).To(BeNil())
		Expect(res.(*directory.GetRootDirectoryPathsResponse).Paths).To(Equal([]struct {
			Path          string
			DirectoryName string
		}{
			{Path: "dir1", DirectoryName: "dir1"},
			{Path: "dir1", DirectoryName: "dir1"},
		}))
	})

	It("will return an error if your root path is a go project", func() {
		req := &directory.GetRootDirectoryPathsRequest{
			RootPath: "sad",
		}

		res, err := directory.GetRootDirectoryPaths(mockContext, req, mockClients)

		Expect(err.Err.Error()).To(Equal("Is a go directory"))
		Expect(res).To(BeNil())
	})

	AfterEach(func() {
		Expect(sqlMock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
