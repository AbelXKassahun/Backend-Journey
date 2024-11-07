package main

import (
	"book-collection/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func rootHandler(w http.ResponseWriter, r *http.Request){
	path := strings.ReplaceAll(r.URL.Path[1:], "-", " ")
	fmt.Fprintf(w, "Through out heaven and earth i alone am %s", path)
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	title := strings.Split(strings.ReplaceAll(r.URL.Path[len("/view/"):], "-", " "), "/")
	
	page, err := utils.LoadFile(title[0])
	
	if err != nil {
		http.Redirect(w, r, "/edit/"+title[0], http.StatusFound)
		// fmt.Fprintf(w, "<h1>404 file not found<h2>")
	}
	
	renderTemplate(w, "../html/View", page)

}

func editHandler(w http.ResponseWriter, r *http.Request){
	title := strings.Split(strings.ReplaceAll(r.URL.Path[len("/edit/"):], "-", " "), "/")

	page, err := utils.LoadFile(title[0])
	
	// create a page if the page is not found in a file
	if err != nil {
		page = &utils.Page{Title: title[0]}
	}

	renderTemplate(w, "../html/Edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request){
	title := strings.ReplaceAll(r.URL.Path[len("/save/"):], "-", " ")
	body := r.FormValue("body")

	page := &utils.Page{Title: title, Body: []byte(body)}
	
	err := page.Save()
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	http.Redirect(w, r, "/view/" + title, http.StatusFound)
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
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	fmt.Println("listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}