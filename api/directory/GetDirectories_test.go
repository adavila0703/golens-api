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

var _ = Describe("GetDirectories", Ordered, func() {
	var mockClients *clients.GlobalClients
	var mock sqlmock.Sqlmock
	var closeDB func() error
	var err error
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	BeforeAll(func() {
		var db *gorm.DB
		db, mock, closeDB, err = clients.NewPostgresClientMock()
		utilsMock := utils.NewMockUtilsClient()
		mockClients = clients.NewGlobalClients(db, nil, utilsMock)
	})

	It("checks for errors on creating mock client", func() {
		Expect(err).To(BeNil())
	})

	It("gets all directory", func() {
		expectedDirectories := []struct {
			ID           uuid.UUID
			CoveragePath string
			CoverageName string
		}{
			{ID: uuid.New(), CoveragePath: "C:\\test1", CoverageName: "test"},
			{ID: uuid.New(), CoveragePath: "C:\\test2", CoverageName: "test2"},
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" WHERE "directories"."deleted_at" IS NULL
		`)).
			WithArgs().
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "coverage_path", "coverage_name"}).
					AddRow(&expectedDirectories[0].ID, &expectedDirectories[0].CoveragePath, &expectedDirectories[0].CoverageName).
					AddRow(&expectedDirectories[1].ID, &expectedDirectories[1].CoveragePath, &expectedDirectories[1].CoverageName),
			)

		res, err := directory.GetDirectories(mockContext, nil, mockClients)
		resMessage := res.(*directory.GetDirectoriesResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Directories[0]["id"]).To(
			BeElementOf(
				expectedDirectories[0].ID.String(),
				expectedDirectories[1].ID.String(),
			),
		)
		Expect(resMessage.Directories[1]["id"]).To(
			BeElementOf(
				expectedDirectories[0].ID.String(),
				expectedDirectories[1].ID.String(),
			),
		)
		Expect(len(resMessage.Directories)).To(Equal(2))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
