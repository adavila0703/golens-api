package ignored_test

import (
	"golens-api/api/ignored"
	"golens-api/clients"
	"golens-api/coverage"
	"golens-api/models"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("CreateIgnored", Ordered, func() {
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

	It("creates ignored", func() {
		expectedUUID := uuid.New()
		req := &ignored.CreateIgnoredRequest{
			Name:        "test",
			IgnoreType:  "1",
			DirectoryID: expectedUUID,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE id = $1 AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			expectedUUID,
		).WillReturnRows(
			sqlmock.NewRows([]string{"coverage_name"}).AddRow(req.Name),
		)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO "ignoreds" ("id","created_at","updated_at","deleted_at","directory_name","name","type") 
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		`)).WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			req.Name,
			req.Name,
			models.DirectoryType,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		_, err := ignored.CreateIgnored(mockContext, req, mockClients)

		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
