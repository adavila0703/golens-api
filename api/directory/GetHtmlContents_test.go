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
)

var _ = Describe("GetHtmlContents", Ordered, func() {
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

	It("returns html content of the given file name", func() {
		req := &directory.GetHtmlContentsRequest{
			FileName:    "test",
			DirectoryID: uuid.New(),
			PackageName: "test",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT * FROM "directories" 
		WHERE id = $1 
		AND "directories"."deleted_at" IS NULL
	`)).WithArgs(
			req.DirectoryID,
		).WillReturnRows(
			sqlmock.NewRows([]string{"coverage_name"}).AddRow(req.FileName),
		)

		directory.ReadHTMLFromFileF = func(name string) (string, error) {
			Expect(name).To(Equal(req.FileName))

			html := `
			<!DOCTYPE html>
			<html>
				<body>
					<div id="topbar">
						<div id="nav">
							<select id="files">
							
							<option value="file0">specs/test/test.go (50.0%)</option>
							
							</select>
						</div>
					</div>
					<div id="content">
					
					<pre class="file" id="file0" style="display: none">test</pre>
					
					</div>
				</body>
			</html>
			`
			return html, nil
		}

		res, err := directory.GetHtmlContents(mockContext, req, mockClients)
		resMessage := res.(*directory.GetHtmlContentsResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.HtmlContent).To(Equal(
			`<div id="content"><pre class="file">1 test
</pre></div>`,
		))
	})

	It("cannot find a directory under the given directory id", func() {
		req := &directory.GetHtmlContentsRequest{
			FileName:    "test",
			DirectoryID: uuid.New(),
			PackageName: "test",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
				SELECT * FROM "directories"
				WHERE id = $1
				AND "directories"."deleted_at" IS NULL
		`)).WithArgs(
			req.DirectoryID,
		).WillReturnRows(
			sqlmock.NewRows([]string{}),
		)

		res, err := directory.GetHtmlContents(mockContext, req, mockClients)
		resMessage := res.(*directory.GetHtmlContentsResponse)

		Expect(err).To(BeNil())
		Expect(resMessage.Message).To(Equal("Directory not found"))
	})

	AfterEach(func() {
		Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		closeDB()
	})
})
