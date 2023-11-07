package ignored_test

import (
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"golens-api/api/ignored"
	"golens-api/clients"
	"golens-api/coverage"
)

var _ = Describe("GetIgnored", Ordered, func() {
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
		req := &ignored.GetIgnoredRequest{}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "ignored" 
			WHERE "ignored"."deleted_at" IS NULL
		`)).
			WithArgs().
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()).AddRow(uuid.New()),
			)

		res, err := ignored.GetIgnored(mockContext, req, mockClients)
		resMessage := res.(*ignored.GetIgnoredResponse)

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
