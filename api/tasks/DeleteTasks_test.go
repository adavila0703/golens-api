package tasks_test

import (
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"golens-api/api"
	"golens-api/api/tasks"
	"golens-api/clients"
)

var _ = Describe("DeleteTasks", Ordered, func() {
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

	It("deletes all tasks", func() {
		req := &tasks.DeleteTasksRequest{}

		expectedID := uuid.New()
		mock.ExpectQuery(`
		SELECT "task_schedules"."id","task_schedules"."created_at","task_schedules"."updated_at","task_schedules"."deleted_at","task_schedules"."schedule_type","task_schedules"."directory_id","Directory"."id" AS "Directory__id","Directory"."created_at" AS "Directory__created_at","Directory"."updated_at" AS "Directory__updated_at","Directory"."deleted_at" AS "Directory__deleted_at","Directory"."path" AS "Directory__path","Directory"."coverage_name" AS "Directory__coverage_name" FROM "task_schedules" LEFT JOIN "directories" "Directory" ON "task_schedules"."directory_id" = "Directory"."id" AND "Directory"."deleted_at" IS NULL WHERE "task_schedules"."deleted_at" IS NULL
		`).
			WithArgs().
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).AddRow(expectedID),
			)

		tasks.DeleteTaskF = func(
			ctx *gin.Context,
			message *tasks.DeleteTaskRequest,
			clients *clients.GlobalClients,
		) (interface{}, *api.Error) {
			Expect(message.TaskID).To(Equal(expectedID))
			return nil, nil
		}

		res, err := tasks.DeleteTasks(mockContext, req, mockClients)
		resMessage := res.(*tasks.DeleteTasksResponse)

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
