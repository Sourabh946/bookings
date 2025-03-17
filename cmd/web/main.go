package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sourabh-go/bookings/internal/config"

	"github.com/alexedwards/scs/v2"
	"github.com/sourabh-go/bookings/internal/handlers"
	"github.com/sourabh-go/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change it to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // Assuming this is a production environment
	app.Session = session

	// Assuming you have a render package with the correct function name
	tc, err := render.CreateTemplateCache() // Fixed the spelling of "createTemplateCache"
	if err != nil {
		panic(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting application on %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	// Use the already declared err variable here, no need to redeclare
	err = srv.ListenAndServe()
	log.Fatal(err)
}
