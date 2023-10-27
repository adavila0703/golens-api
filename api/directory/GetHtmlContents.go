package directory

import (
	"fmt"
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type GetHtmlContentsRequest struct {
	FileName    string    `json:"fileName" validate:"required"`
	DirectoryID uuid.UUID `json:"directoryId" validate:"required"`
}

type GetHtmlContentsResponse struct {
	Message     string `json:"message"`
	HtmlContent string `json:"htmlContent"`
}

func GetHtmlContents(
	ctx *gin.Context,
	message *GetHtmlContentsRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	directory, found, err := models.GetDirectory(ctx, clients.DB, message.DirectoryID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	if !found {
		return &GetHtmlContentsResponse{
			Message: "Directory not found",
		}, nil
	}

	htmlString, err := readHTMLFromFile(directory.CoverageName)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	optionsMap, err := getSelectOptionsMap(htmlString)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	var fileID string
	for optionsKey, optionsValue := range optionsMap {
		if strings.Contains(optionsValue, message.FileName) {
			fileID = optionsKey
		}
	}

	content, err := getElementContentByID(htmlString, fileID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	content = insertLineCount(content)

	htmlContent := fmt.Sprintf(`<div id="content"><pre class="file">%s</pre></div>`, content)

	return &GetHtmlContentsResponse{
		Message:     "Good!",
		HtmlContent: htmlContent,
	}, nil
}

func insertLineCount(htmlString string) string {
	lines := strings.Split(htmlString, "\n")
	result := ""

	for i, line := range lines {
		lineNumber := i + 1
		result += fmt.Sprintf("%d %s\n", lineNumber, line)
	}

	return result
}

func getElementContentByID(htmlString string, id string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return "", err
	}

	var getElementContent func(*html.Node) string
	getElementContent = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "pre" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == id {
					var content strings.Builder
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						html.Render(&content, c)
					}
					return content.String()
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if result := getElementContent(c); result != "" {
				return result
			}
		}

		return ""
	}

	return getElementContent(doc), nil
}

func readHTMLFromFile(name string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}
	fileName := fmt.Sprintf("%s/data/html/%s.html", workingDir, name)

	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getSelectOptionsMap(htmlString string) (map[string]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return nil, err
	}

	options := make(map[string]string)
	var parse func(*html.Node)
	parse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "option" {
			if len(n.Attr) > 0 {
				for _, attr := range n.Attr {
					if attr.Key == "value" {
						options[attr.Val] = n.FirstChild.Data
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parse(c)
		}
	}

	parse(doc)
	return options, nil
}
