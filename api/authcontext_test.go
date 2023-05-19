package api_test

// import (
// 	"blockbro-api/api"
// 	"blockbro-api/clients"
// 	"net/http/httptest"
// 	"regexp"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/gin-gonic/gin"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// 	"gorm.io/gorm"
// )

// var _ = Describe("Authcontext", func() {
// 	var mockClients *clients.GlobalClients
// 	var mock sqlmock.Sqlmock
// 	var closeDB func() error
// 	var err error
// 	mockRecorder := httptest.NewRecorder()
// 	mockContext, _ := gin.CreateTestContext(mockRecorder)

// 	expectedDeviceUUID := "1234567890"
// 	expectedUsername := "adavila0703"

// 	BeforeEach(func() {
// 		var db *gorm.DB
// 		db, mock, closeDB, err = clients.NewPostgresClientMock()
// 		mockClients = clients.NewGlobalClients(db, nil)
// 	})

// 	It("checks for errors on creating mock client", func() {
// 		Expect(err).To(BeNil())
// 	})

// 	It("creates the auth context", func() {

// 		mock.ExpectQuery(regexp.QuoteMeta(`
// 			SELECT device_uuid, username
// 			FROM "users"
// 			WHERE device_uuid = $1
// 			AND "users"."deleted_at" IS NULL
// 			ORDER BY "users"."id" LIMIT 1
// 		`)).
// 			WithArgs(expectedDeviceUUID).
// 			WillReturnRows(
// 				sqlmock.NewRows(
// 					[]string{"device_uuid", "username"}).
// 					AddRow(expectedDeviceUUID, expectedUsername),
// 			)

// 		authContext := api.GetAuthContext(mockContext, expectedDeviceUUID, mockClients)

// 		Expect(err).To(BeNil())
// 		Expect(authContext.Username).To(Equal(expectedUsername))
// 		Expect(authContext.DeviceUUID).To(Equal(expectedDeviceUUID))
// 	})

// 	It("creates a user if non exists", func() {
// 		expectedDeviceUUID := "1234567890"

// 		mock.ExpectQuery(regexp.QuoteMeta(`
// 			SELECT device_uuid, username
// 			FROM "users"
// 			WHERE device_uuid = $1
// 			AND "users"."deleted_at" IS NULL
// 			ORDER BY "users"."id" LIMIT 1
// 		`)).
// 			WithArgs(expectedDeviceUUID).
// 			WillReturnRows(
// 				sqlmock.NewRows(
// 					[]string{"device_uuid", "username"}),
// 			)

// 		authContext := api.GetAuthContext(mockContext, mockClients)

// 		Expect(authContext.Username).To(Equal("TempUserName"))
// 	})

// 	AfterEach(func() {
// 		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
// 		closeDB()
// 	})
// })
