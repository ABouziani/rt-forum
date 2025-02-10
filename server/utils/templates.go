package utils

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	"forum/server/models"
)

type GlobalData struct {
	IsAuthenticated bool
	Data            any
	UserName        string
	Categories      []models.Category
}

type Error struct {
	Code    int
	Message string
	Details string
}

func ParseTemplates(tmpl string) (*template.Template, error) {
	// Parse the template files
	t, err := template.New(tmpl).Parse(HtmlTemplates[tmpl])
	if err != nil {
		return nil, fmt.Errorf("error parsing template files: %w", err)
	}
	t, err = t.Parse(HtmlTemplates["header"])
	if err != nil {
		return nil, fmt.Errorf("error parsing template files: %w", err)
	}
	t, err = t.Parse(HtmlTemplates["navbar"])
	if err != nil {
		return nil, fmt.Errorf("error parsing template files: %w", err)
	}
	t, err = t.Parse(HtmlTemplates["footer"])
	if err != nil {
		return nil, fmt.Errorf("error parsing template files: %w", err)
	}

	return t, nil
}

// RenderError handles error responses
func RenderError(db *sql.DB, w http.ResponseWriter, r *http.Request, statusCode int, isauth bool, username string) {
	typeError := Error{
		Code:    statusCode,
		Message: http.StatusText(statusCode),
	}
	if err := RenderTemplate(db, w, r, "error", statusCode, typeError, isauth, username); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(ErrorPageContents))
	}
}

func RenderTemplate(db *sql.DB, w http.ResponseWriter, r *http.Request, tmpl string, statusCode int, data any, isauth bool, username string) error {
	t, err := ParseTemplates(tmpl)
	if err != nil {
		return err
	}
	// Fetch categories for the navigation bar
	var categories []models.Category
	if db != nil {
		categories, err = models.FetchCategories(db)
		if err != nil {
			categories = nil
		}
	}

	globalData := GlobalData{
		IsAuthenticated: isauth,
		Data:            data,
		UserName:        username,
		Categories:      categories,
	}
	w.WriteHeader(statusCode)

	var buf bytes.Buffer
	// Execute the template with the provided data

	err = t.ExecuteTemplate(&buf, tmpl, globalData)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	w.Header().Set("Content-Type", "text/html")
	buf.WriteTo(w)
	return nil
}
