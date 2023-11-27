package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Rsvp struct {
	Name, Email, Phone string
	WillAttend         bool
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	templates["welcome"].Execute(w, nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	templates["list"].Execute(w, responses)
}

type formData struct {
	*Rsvp
	Errors []string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates["form"].Execute(w, formData{
			Rsvp: &Rsvp{}, Errors: []string{},
		})
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		responseData := Rsvp{
			Name:       r.Form["name"][0],
			Email:      r.Form["email"][0],
			Phone:      r.Form["phone"][0],
			WillAttend: r.Form["willattend"][0] == "true",
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Name is required")
		}
		if responseData.Email == "" {
			errors = append(errors, "Email is required")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Phone is required")
		}
		if len(errors) > 0 {
			templates["form"].Execute(w, formData{
				Rsvp: &responseData, Errors: errors,
			})
		} else {

			responses = append(responses, &responseData)

			if responseData.WillAttend {
				templates["thanks"].Execute(w, responseData.Name)
			} else {
				templates["sorry"].Execute(w, responseData.Name)
			}
		}
	}
}

func main() {
	loadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}
