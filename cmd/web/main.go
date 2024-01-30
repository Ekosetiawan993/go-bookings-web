package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ekosetiawan993/go-bookings-web/pkg/config"
	"github.com/Ekosetiawan993/go-bookings-web/pkg/handlers"
	"github.com/Ekosetiawan993/go-bookings-web/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8084"

// 41 all main package could access this
var app config.AppConfig

// 41 declare session here for middleware
var session *scs.SessionManager

func main() {

	// change this to true in production
	app.InProduction = false

	// vid 41 session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	// should the session persist after someone close the window?
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	// we set the secure and samesite in several places, it is better to store the value inside appConfig
	session.Cookie.Secure = app.InProduction // true if using https

	app.Session = session

	tc, err := render.CreateFullTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = true

	// create repo codes
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	//

	// call the NewTemplate function for render
	render.NewTemplate(&app)

	// handlefunc without using chi / third party middleware
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	n, err := fmt.Fprintf(w, "Hello Eko")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	fmt.Println(fmt.Sprintf("The bytes %d", n))

	// })
	fmt.Printf("Start server on port %v", portNumber)
	// http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
