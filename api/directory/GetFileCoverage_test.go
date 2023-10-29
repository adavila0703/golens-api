package directory_test

import (
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"golens-api/api/directory"
	"golens-api/clients"
	"golens-api/utils"
)

var _ = Describe("GetFileCoverage", Ordered, func() {
	var mockClients *clients.GlobalClients
	var mock sqlmock.Sqlmock
	var closeDB func() error
	var err error
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	BeforeAll(func() {
		var db *gorm.DB
		db, mock, closeDB, err = clients.NewPostgresClientMock()
		mockClients = clients.NewGlobalClients(db, nil)
	})

	It("checks for errors on creating mock client", func() {
		Expect(err).To(BeNil())
	})

	It("gets file code coverage", func() {
		expectedCoverageName := "test"
		expectedDirectoryID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE id = $1 
			AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			expectedDirectoryID,
		).WillReturnRows(
			sqlmock.NewRows([]string{"coverage_name"}).AddRow(expectedCoverageName),
		)

		expectedCoverageMap := []map[string]any{
			{
				"fileName":     "file1",
				"totalLines":   1000,
				"coveredLines": 1000,
			},
			{
				"fileName":     "file2",
				"totalLines":   1000,
				"coveredLines": 1000,
			},
		}

		expectedPackageName := "testPackage"

		utils.GetFileCoveragePercentageF = func(coverageName string) (map[string][]map[string]any, error) {
			Expect(coverageName).To(Equal(expectedCoverageName))
			return map[string][]map[string]any{
				expectedPackageName: expectedCoverageMap,
			}, nil
		}

		req := &directory.GetFileCoverageRequest{
			RepoID:      expectedDirectoryID,
			PackageName: expectedPackageName,
		}

		res, err := directory.GetFileCoverage(mockContext, req, mockClients)
		resMessage := res.(*directory.GetFileCoverageResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.FileCoverage).To(Equal(expectedCoverageMap))
	})

	It("returns no directory found message", func() {
		expectedDirectoryID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE id = $1 
			AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			expectedDirectoryID,
		).WillReturnRows(
			sqlmock.NewRows([]string{""}),
		)

		req := &directory.GetFileCoverageRequest{
			RepoID: expectedDirectoryID,
		}

		res, err := directory.GetFileCoverage(mockContext, req, mockClients)
		resMessage := res.(*directory.GetFileCoverageResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Message).To(Equal("Directory not found"))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
