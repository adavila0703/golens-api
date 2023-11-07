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

var _ = Describe("DeleteIgnored", Ordered, func() {
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

	It("deletes an ignored", func() {
		req := &ignored.DeleteIgnoredRequest{
			ID: uuid.New(),
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			UPDATE "ignored" 
			SET "deleted_at"=$1 
			WHERE id = $2 AND "ignored"."deleted_at" IS NULL
		`)).WithArgs(
			sqlmock.AnyArg(),
			req.ID,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		_, err := ignored.DeleteIgnored(mockContext, req, mockClients)

		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
