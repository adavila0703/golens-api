package tasks_test

import (
	"net/http/httptest"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"golens-api/api/tasks"
	"golens-api/clients"
	"golens-api/utils"
)

type DeleteTaskCron struct {
	clients.ICron
}

func NewDeleteTaskCron() *DeleteTaskCron {
	return &DeleteTaskCron{}
}

func (d *DeleteTaskCron) RemoveCronJob(id cron.EntryID) {
}

var _ = Describe("DeleteTask", Ordered, func() {
	var mockClients *clients.GlobalClients
	var mock sqlmock.Sqlmock
	var closeDB func() error
	var err error
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	BeforeAll(func() {
		var db *gorm.DB
		db, mock, closeDB, err = clients.NewPostgresClientMock()
		cron := NewDeleteTaskCron()
		mockClients = clients.NewGlobalClients(db, cron, nil)
	})

	It("checks for errors on creating mock client", func() {
		Expect(err).To(BeNil())
	})

	It("deletes a task", func() {
		req := &tasks.DeleteTaskRequest{
			TaskID:       uuid.New(),
			ScheduleType: utils.EveryMinute,
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			UPDATE "task_schedules" SET "deleted_at"=$1 
			WHERE id = $2 AND "task_schedules"."deleted_at" IS NULL
		`)).WithArgs(
			sqlmock.AnyArg(),
			req.TaskID,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "task_schedules" 
			WHERE schedule_type = $1 AND "task_schedules"."deleted_at" IS NUL
		`)).WithArgs(
			req.ScheduleType,
		).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()),
		)

		res, err := tasks.DeleteTask(mockContext, req, mockClients)
		resMessage := res.(*tasks.DeleteTaskResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Message).To(Equal("Good!"))
	})

	It("will delete all cron jobs if there are no more tasks scheduled", func() {
		req := &tasks.DeleteTaskRequest{
			TaskID:       uuid.New(),
			ScheduleType: utils.EveryMinute,
		}

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			UPDATE "task_schedules" SET "deleted_at"=$1 
			WHERE id = $2 AND "task_schedules"."deleted_at" IS NULL
		`)).WithArgs(
			sqlmock.AnyArg(),
			req.TaskID,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "task_schedules" 
			WHERE schedule_type = $1 AND "task_schedules"."deleted_at" IS NUL
		`)).WithArgs(
			req.ScheduleType,
		).WillReturnRows(
			sqlmock.NewRows([]string{}),
		)

		expectedCronID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "cron_jobs" 
			WHERE schedule_type = $1 AND "cron_jobs"."deleted_at" IS NULL
		`)).WithArgs(
			req.ScheduleType,
		).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(expectedCronID),
		)

		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			UPDATE "cron_jobs" SET "deleted_at"=$1 
			WHERE id = $2 AND "cron_jobs"."deleted_at" IS NULL
		`)).WithArgs(
			sqlmock.AnyArg(),
			expectedCronID,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		res, err := tasks.DeleteTask(mockContext, req, mockClients)
		resMessage := res.(*tasks.DeleteTaskResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Message).To(Equal("Good!"))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
