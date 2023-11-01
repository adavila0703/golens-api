package ignore_directory_test

import (
	"golens-api/api/ignore_directory"
	"golens-api/clients"
	"golens-api/coverage"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("CreateIgnoredDirectory", Ordered, func() {
	var mockClients *clients.GlobalClients
	var mock sqlmock.Sqlmock
	var closeDB func() error
	var err error
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	BeforeAll(func() {
		var db *gorm.DB
		db, mock, closeDB, err = clients.NewPostgresClientMock()
		utilsMock := coverage.NewCoverageMock()
		mockClients = clients.NewGlobalClients(db, nil, utilsMock)
	})

	It("checks for errors on creating mock client", func() {
		Expect(err).To(BeNil())
	})

	It("creates ignored directory", func() {
		req := &ignore_directory.CreateIgnoredDirectoryRequest{
			DirectoryName: "test",
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO "ignored_directories" 
			("id","created_at","updated_at","deleted_at","directory_name") 
			VALUES ($1,$2,$3,$4,$5)
		`)).WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			req.DirectoryName,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		_, err := ignore_directory.CreateIgnoredDirectory(mockContext, req, mockClients)

		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
