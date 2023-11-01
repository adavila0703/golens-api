package tasks_test

import (
	"golens-api/api/tasks"
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

var _ = Describe("CreateTasks", Ordered, func() {
	var mockClients *clients.GlobalClients
	var mock sqlmock.Sqlmock
	var closeDB func() error
	var err error
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	BeforeAll(func() {
		var db *gorm.DB
		db, mock, closeDB, err = clients.NewPostgresClientMock()
		mockClients = clients.NewGlobalClients(db, nil, nil)
	})

	It("checks for errors on creating mock client", func() {
		Expect(err).To(BeNil())
	})

	It("creates a tasks for all current directories", func() {
		req := &tasks.CreateTasksRequest{
			ScheduleType: utils.EveryMinute,
		}

		expectedDirID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "cron_jobs" 
			WHERE schedule_type = $1 
			AND "cron_jobs"."deleted_at" IS NULL
		`)).WithArgs(
			req.ScheduleType,
		).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()),
		)

		expectedCoverageName := "test"
		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" WHERE "directories"."deleted_at" IS NULL
		`)).
			WithArgs().
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "coverage_name"}).
					AddRow(expectedDirID, expectedCoverageName),
			)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO "task_schedules" ("id","created_at","updated_at","deleted_at","schedule_type","directory_id") 
			VALUES ($1,$2,$3,$4,$5,$6)
		`)).WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			utils.EveryMinute,
			expectedDirID,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		res, err := tasks.CreateTasks(mockContext, req, mockClients)
		resMessage := res.(*tasks.CreateTasksResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Tasks[0]["DirectoryID"]).To(Equal(expectedDirID))
		Expect(resMessage.Tasks[0]["id"]).To(Equal(1))
		Expect(resMessage.Tasks[0]["coverageName"]).To(Equal(expectedCoverageName))
		Expect(resMessage.Tasks[0]["scheduleTypeName"]).To(Equal("Daily"))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
