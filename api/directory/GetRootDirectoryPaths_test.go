package directory_test

import (
	"golens-api/api/directory"
	"golens-api/clients"
	"golens-api/utils"
	"net/http/httptest"
	"regexp"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
)

type GetRootDirectoryPathsUtils struct {
	utils.IUtilsClient
}

func NewGetRootDirectoryPathsUtils() *GetRootDirectoryPathsUtils {
	return &GetRootDirectoryPathsUtils{}
}

func (g *GetRootDirectoryPathsUtils) IsGoDirectory(dirPath string) (bool, error) {
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
		utilsMock := NewGetRootDirectoryPathsUtils()
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
			SELECT * FROM "ignored_directories" 
			WHERE "ignored_directories"."deleted_at" IS NULL
		`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"directory_name"}).AddRow("test"),
			)

		res, err := directory.GetRootDirectoryPaths(mockContext, req, mockClients)

		Expect(err).To(BeNil())
		Expect(res.(*directory.GetRootDirectoryPathsResponse).Paths).To(Equal([]string{"dir1", "dir1"}))
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
