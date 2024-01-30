package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Ekosetiawan993/go-bookings-web/pkg/config"
	"github.com/Ekosetiawan993/go-bookings-web/pkg/models"
)

var app *config.AppConfig

// function for accessing app config from main file
func NewTemplate(a *config.AppConfig) {
	app = a
}

// always parse and render in each request : inefficient
func RenderTemplateTest1(w http.ResponseWriter, tmpl string) {
	// parsefiles := can parse multiple files, INCLUDING the layout file
	parsedTemplate, err := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")

	if err != nil {
		fmt.Println("error parsing template and layout:", err)
		return
	}

	err = parsedTemplate.Execute(w, nil)

	if err != nil {
		fmt.Println("error executing parsed template:", err)
		return
	}

}

var tempateCache = make(map[string]*template.Template)

// a better (one level) rendering function
// parse new template when needed
func RenderTemplateTest2(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	_, inMap := tempateCache[t]

	// if template on cache
	if !inMap {
		log.Println("Create template and save it to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// using cached template
		log.Println("Reusing cached template")
	}

	// first created template
	tmpl = tempateCache[t]
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("error executing first created template", err)
	}

}

func createTemplateCache(tmpName string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", tmpName),
		"./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println("Error parsing templates", err)
		return err
	}
	tempateCache[tmpName] = tmpl
	return nil
}

// third way of parsing template
func RenderTemplateWithoutConfig(w http.ResponseWriter, tmpl string) {
	// create template cache
	tc, err := CreateFullTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// access cahced template
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(err)
	}

	// use buffer?
	buf := new(bytes.Buffer)

	err = t.Execute(buf, nil)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

func CreateFullTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	// create entire template cache at once
	// get all files named *.html
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	// loop through all html files
	for _, page := range pages {
		name := filepath.Base(page)
		// set the name and parse each page's template
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		log.Println("parse tempalte of: ", name)

		// base layout
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		// add base template to each page
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

// vid 36 : pass default data to every pages
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// the tempate data params from vid 36
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// vid 35 code : use cache or rebuild everytime
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache
		tc = app.TemplateCache
	} else {
		tc, _ = CreateFullTemplateCache()
	}

	// create template cache
	// tc := app.TemplateCache

	// access cahced template
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("cannot get template from cached")
	}

	// use buffer?
	buf := new(bytes.Buffer)

	// vid 36
	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}
