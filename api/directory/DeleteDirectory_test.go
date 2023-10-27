package directory_test

import (
	"fmt"
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

var _ = Describe("DeleteDirectory", Ordered, func() {
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

	It("deletes a directory", func() {
		expectedID := uuid.New()

		req := &directory.DeleteDirectoryRequest{
			ID: expectedID,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE id = $1 AND "directories"."deleted_at" IS NULL
		`)).
			WithArgs(expectedID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()),
			)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			UPDATE "directories" 
			SET "deleted_at"=$1 
			WHERE id = $2 AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			sqlmock.AnyArg(),
			expectedID,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		utils.GetWorkingDirectoryF = func() (string, error) {
			return "./", nil
		}

		res, err := directory.DeleteDirectory(mockContext, req, mockClients)
		fmt.Println(res, err)

		// Expect(err).To(BeNil())
		// Expect(res.(*directory.CreateDirectoryResponse).Directory["path"]).To(Equal(expectedPath))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
