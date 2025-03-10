package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/sourabh-go/bookings/pkg/config" // Ensure config package is imported correctly
	"github.com/sourabh-go/bookings/pkg/models"
)

var app *config.AppConfig // Ensure correct struct name capitalization

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template Cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// Get request template from cache
	t, ok := tc[tmpl]
	if !ok {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template %s not found in cache\n", tmpl)
		return
	}

	buf := new(bytes.Buffer)

	// Execute template and check for errors immediately
	err := t.Execute(buf, td)
	if err != nil {
		log.Printf("Error executing template: %v\n", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Write to response and check for errors
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Error writing to response: %v\n", err)
		http.Error(w, "Error sending response", http.StatusInternalServerError)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// Get all the files named *.page.tmpl from ./templates directory
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// Range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// Include layout templates
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			tmpl, err = tmpl.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = tmpl
	}

	return myCache, nil
}
