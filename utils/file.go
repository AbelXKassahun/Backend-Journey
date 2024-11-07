package utils

import "os"

type Page struct {
	Title string
	Body  []byte
}


// save to file
func (p *Page) Save() error {
	fileName := p.Title + ".txt"
	return os.WriteFile("../data/" + fileName, p.Body, 0600)
}

// load a files content
func LoadFile(title string) (*Page, error) {
	filename := title + ".txt"
	body ,err := os.ReadFile("../data/" + filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}