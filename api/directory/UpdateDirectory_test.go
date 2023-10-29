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

var _ = Describe("UpdateDirectory", Ordered, func() {
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

	It("updates a directory", Focus, func() {
		expectedPath := "path"
		expectedCoverageName := "test"
		req := &directory.UpdateDirectoryRequest{
			ID: uuid.New(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" WHERE id = $1 AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			req.ID,
		).WillReturnRows(
			sqlmock.NewRows([]string{"id", "path", "coverage_name"}).
				AddRow(req.ID, expectedPath, expectedCoverageName),
		)

		utils.GenerateCoverageAndHTMLFilesF = func(path string) error {
			Expect(path).To(Equal(expectedPath))
			return nil
		}

		utils.GetCoveredLinesF = func(coverageName string) (int, int, error) {
			Expect(coverageName).To(Equal(expectedCoverageName))
			return 1000, 1000, nil
		}

		res, err := directory.UpdateDirectory(mockContext, req, mockClients)
		resMessage := res.(*directory.UpdateDirectoryResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Directory["id"]).To(Equal(req.ID))
		Expect(resMessage.Directory["path"]).To(Equal(expectedPath))
		Expect(resMessage.Directory["coverageName"]).To(Equal(expectedCoverageName))
		Expect(resMessage.Directory["totalLines"]).To(Equal(1000))
		Expect(resMessage.Directory["coveredLines"]).To(Equal(1000))
	})

	It("will not find the directory", func() {
		req := &directory.UpdateDirectoryRequest{
			ID: uuid.New(),
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE id = $1 AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			req.ID,
		).WillReturnRows(
			sqlmock.NewRows([]string{}),
		)

		res, err := directory.UpdateDirectory(mockContext, req, mockClients)

		Expect(err).To(BeNil())
		Expect(res).To(Equal(nil))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
