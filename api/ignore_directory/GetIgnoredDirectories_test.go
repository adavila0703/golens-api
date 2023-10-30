package ignore_directory_test

import (
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"golens-api/api/ignore_directory"
	"golens-api/clients"
	"golens-api/coverage"
)

var _ = Describe("GetIgnoredDirectories", Ordered, func() {
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

	It("get ignored directory", func() {
		req := &ignore_directory.GetIgnoredDirectoriesRequest{}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "ignored_directories" 
			WHERE "ignored_directories"."deleted_at" IS NULL
		`)).
			WithArgs().
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New()),
			)

		res, err := ignore_directory.GetIgnoredDirectories(mockContext, req, mockClients)
		resMessage := res.(*ignore_directory.GetIgnoredDirectoriesResponse)

		Expect(err).To(BeNil())
		Expect(len(resMessage.Directories)).To(Equal(2))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
