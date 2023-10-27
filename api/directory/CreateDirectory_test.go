package directory_test

import (
	"golens-api/api/directory"
	"golens-api/clients"
	"golens-api/utils"
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("CreateDirectory", Ordered, func() {
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

	It("creates a directory", func() {
		expectedPath := "C:\\Dir\\path"
		expectedCoverageName := "path"

		req := &directory.CreateDirectoryRequest{
			Path: expectedPath,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE path = $1 AND "directories"."deleted_at" IS NULL
		`)).
			WithArgs(expectedPath).
			WillReturnRows(
				sqlmock.NewRows([]string{}),
			)

		utils.IsGoDirectoryF = func(dirPath string) (bool, error) {
			Expect(dirPath).To(Equal(expectedPath))
			return true, nil
		}

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE "directories"."path" = $1
				AND "directories"."coverage_name" = $2 
				AND "directories"."deleted_at" IS NULL 
			ORDER BY "directories"."id" 
			LIMIT 1
		`)).WithArgs(
			expectedPath, expectedCoverageName,
		).WillReturnRows(
			sqlmock.NewRows([]string{}),
		)

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO "directories" ("id","created_at","updated_at","deleted_at","path","coverage_name") 
			VALUES ($1,$2,$3,$4,$5,$6)
		`)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				expectedPath,
				expectedCoverageName,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		utils.GenerateCoverageAndHTMLFilesF = func(path string) error {
			Expect(path).To(Equal(expectedPath))

			return nil
		}

		mock.ExpectCommit()

		utils.GetCoveredLinesF = func(coverageName string) (int, int, error) {
			Expect(coverageName).To(Equal(expectedCoverageName))
			return 1000, 1000, nil
		}

		res, err := directory.CreateDirectory(mockContext, req, mockClients)

		Expect(err).To(BeNil())
		Expect(res.(*directory.CreateDirectoryResponse).Directory["path"]).To(Equal(expectedPath))
	})

	It("fails when the directory is not a go project", func() {
		expectedPath := "C:\\Dir\\path"

		req := &directory.CreateDirectoryRequest{
			Path: expectedPath,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE path = $1 AND "directories"."deleted_at" IS NULL
		`)).
			WithArgs(expectedPath).
			WillReturnRows(
				sqlmock.NewRows([]string{}),
			)

		utils.IsGoDirectoryF = func(dirPath string) (bool, error) {
			Expect(dirPath).To(Equal(expectedPath))
			return false, nil
		}

		res, err := directory.CreateDirectory(mockContext, req, mockClients)

		Expect(err.Err.Error()).To(Equal("Is not a go directory"))
		Expect(res).To(BeNil())
	})

	It("returns nil if the directory already exists", func() {
		expectedPath := "C:\\Dir\\path"

		req := &directory.CreateDirectoryRequest{
			Path: expectedPath,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE path = $1 AND "directories"."deleted_at" IS NULL
		`)).
			WithArgs(expectedPath).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()),
			)

		res, err := directory.CreateDirectory(mockContext, req, mockClients)

		Expect(err).To(BeNil())
		Expect(res).To(BeNil())
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
