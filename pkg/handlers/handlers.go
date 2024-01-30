package handlers

import (
	"net/http"

	"github.com/Ekosetiawan993/go-bookings-web/pkg/config"
	"github.com/Ekosetiawan993/go-bookings-web/pkg/models"
	"github.com/Ekosetiawan993/go-bookings-web/pkg/render"
)

// new code
var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// create a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// newhandlers sets the repo for the handler
func NewHandlers(r *Repository) {
	Repo = r
}

//

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// 42, try session
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// old home handlers
// func Home1(w http.ResponseWriter, r *http.Request) {
// 	render.RenderTemplate(w, "home.page.html")
// }

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// vid 36
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	// 42
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
