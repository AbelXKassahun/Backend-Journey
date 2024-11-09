package main

import (
	"book-collection/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var validURLPattern = "^/(edit|save|view)/([a-zA-Z0-9]+)$"
// var validURLPattern = "^/(edit|save|view)$"
var validPath = regexp.MustCompile(validURLPattern)

var templates = template.Must(template.ParseFiles("../html/Edit.html", "../html/View.html"))

func rootHandler(w http.ResponseWriter, r *http.Request){
	path := strings.ReplaceAll(r.URL.Path[1:], "-", " ")
	fmt.Fprintf(w, "Through out heaven and earth i alone am %s", path)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string){
	page, err := utils.LoadFile(title)
	
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		// fmt.Fprintf(w, "<h1>404 file not found<h2>")
	}
	
	renderTemplate(w, "View", page)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string){
	page, err := utils.LoadFile(title)
	
	// create a page if the page is not found in a file
	if err != nil {
		page = &utils.Page{Title: title}
	}

	renderTemplate(w, "Edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string){
	body := r.FormValue("body")
	page := &utils.Page{Title: title, Body: []byte(body)}	
	err := page.Save()
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}


func makeHandlder(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		
		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *utils.Page) {
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main(){
	p1 := &utils.Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    p1.Save()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandlder(viewHandler))
	http.HandleFunc("/edit/", makeHandlder(editHandler))
	http.HandleFunc("/save/", makeHandlder(saveHandler))

	fmt.Println("listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}