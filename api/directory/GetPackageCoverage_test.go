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
	"golens-api/coverage"
)

type GetPackageCoverageCoverage struct {
	coverage.ICoverage
}

func NewGetPackageCoverageCoverage() *GetPackageCoverageCoverage {
	return &GetPackageCoverageCoverage{}
}

func (g *GetPackageCoverageCoverage) GetCoveredLinesByPackage(
	coverageName string,
	ignoredFilesByPackage map[string]map[string]bool,
	ignoredPackages map[string]bool,
) (map[string]map[string]int, error) {
	return map[string]map[string]int{
		"test1": {
			"totalLines":   1000,
			"coveredLines": 500,
		},
	}, nil
}

var _ = Describe("GetPackageCoverage", Ordered, func() {
	var mockClients *clients.GlobalClients
	var mock sqlmock.Sqlmock
	var closeDB func() error
	var err error
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	BeforeAll(func() {
		var db *gorm.DB
		db, mock, closeDB, err = clients.NewPostgresClientMock()
		utilsMock := NewGetPackageCoverageCoverage()
		mockClients = clients.NewGlobalClients(db, nil, utilsMock)
	})

	It("checks for errors on creating mock client", func() {
		Expect(err).To(BeNil())
	})

	It("returns the package coverage", func() {
		expectedCoverageName := "test"
		req := &directory.GetPackageCoverageRequest{
			ID: uuid.New(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE id = $1 
			AND "directories"."deleted_at" IS NULL
	`)).WithArgs(
			req.ID,
		).WillReturnRows(
			sqlmock.NewRows([]string{"coverage_name"}).AddRow(expectedCoverageName),
		)

		res, err := directory.GetPackageCoverage(mockContext, req, mockClients)
		resMessage := res.(*directory.GetPackageCoverageResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.PackageCoverage[0]["packageName"]).To(Equal("test1"))
		Expect(resMessage.PackageCoverage[0]["totalLines"]).To(Equal(1000))
	})

	It("will not find the directory", func() {
		req := &directory.GetPackageCoverageRequest{
			ID: uuid.New(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE id = $1 
			AND "directories"."deleted_at" IS NULL
	`)).WithArgs(
			req.ID,
		).WillReturnRows(
			sqlmock.NewRows([]string{}),
		)

		res, err := directory.GetPackageCoverage(mockContext, req, mockClients)
		resMessage := res.(*directory.GetPackageCoverageResponse)

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
