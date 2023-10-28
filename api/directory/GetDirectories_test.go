package directory_test

import (
	"fmt"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"golens-api/api/directory"
	"golens-api/clients"
)

var _ = Describe("GetDirectories", Ordered, func() {
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

	It("gets all directory", Focus, func() {
		expectedIDs := []string{}


		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" WHERE "directories"."deleted_at" IS NULL
		`)).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{}).
			AddRow()
		)

		res, err := directory.GetDirectories(mockContext, nil, mockClients)
		fmt.Println(res)

		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
