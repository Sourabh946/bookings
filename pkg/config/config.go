package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the configuration config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLogger    *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
