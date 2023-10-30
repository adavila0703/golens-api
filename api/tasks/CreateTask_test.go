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
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CreateTaskCron struct {
	clients.ICron
}

func NewCreateTaskCron() *CreateTaskCron {
	return &CreateTaskCron{}
}

func (c *CreateTaskCron) CreateCronJob(schedule utils.CronJobScheduleType, handler func()) (cron.EntryID, error) {
	return 1, nil
}

var _ = Describe("CreateTask", Ordered, func() {
	var mockClients *clients.GlobalClients
	var mock sqlmock.Sqlmock
	var closeDB func() error
	var err error
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

	BeforeAll(func() {
		var db *gorm.DB
		db, mock, closeDB, err = clients.NewPostgresClientMock()
		cronMock := NewCreateTaskCron()
		mockClients = clients.NewGlobalClients(db, cronMock, nil)
	})

	It("checks for errors on creating mock client", func() {
		Expect(err).To(BeNil())
	})

	It("creates a task on an already created cron job", func() {
		req := &tasks.CreateTaskRequest{
			DirectoryID:  uuid.New(),
			ScheduleType: utils.EveryMinute,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "cron_jobs" 
			WHERE schedule_type = $1 
			AND "cron_jobs"."deleted_at" IS NULL
		`)).WithArgs(
			req.ScheduleType,
		).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(req.DirectoryID),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
			WHERE id = $1 
			AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			req.DirectoryID,
		).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(req.DirectoryID),
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
			req.DirectoryID,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		res, err := tasks.CreateTask(mockContext, req, mockClients)
		resMessage := res.(*tasks.CreateTaskResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Task.DirectoryID).To(Equal(req.DirectoryID))
	})

	It("creates a task and a cron job", func() {
		req := &tasks.CreateTaskRequest{
			DirectoryID:  uuid.New(),
			ScheduleType: utils.EveryMinute,
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "cron_jobs" 
			WHERE schedule_type = $1 
			AND "cron_jobs"."deleted_at" IS NULL
		`)).WithArgs(
			req.ScheduleType,
		).WillReturnRows(
			sqlmock.NewRows([]string{}),
		)

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "cron_jobs" 
			WHERE "cron_jobs"."schedule" = $1 
			AND "cron_jobs"."schedule_type" = $2 
			AND "cron_jobs"."entry_id" = $3 
			AND "cron_jobs"."deleted_at" IS NULL 
			ORDER BY "cron_jobs"."id" LIMIT 1
		`)).WithArgs(
			utils.GetCronSchedule(utils.EveryMinute),
			utils.EveryMinute,
			1,
		).WillReturnRows(
			sqlmock.NewRows([]string{}),
		)

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO "cron_jobs" ("id","created_at","updated_at","deleted_at","schedule","schedule_type","entry_id") 
			VALUES ($1,$2,$3,$4,$5,$6,$7)
		`)).WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			utils.GetCronSchedule(utils.EveryMinute),
			utils.EveryMinute,
			1,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT * FROM "directories" 
		 	WHERE id = $1 
			AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			req.DirectoryID,
		).WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(req.DirectoryID),
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
			req.DirectoryID,
		).WillReturnResult(
			sqlmock.NewResult(1, 1),
		)

		mock.ExpectCommit()

		res, err := tasks.CreateTask(mockContext, req, mockClients)
		resMessage := res.(*tasks.CreateTaskResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Task.DirectoryID).To(Equal(req.DirectoryID))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
