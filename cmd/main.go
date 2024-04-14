package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Block struct {
    Id int
}

type Blocks struct {
    Start int
    Next int
    More bool
    Blocks []Block
}

var id = 0

type Contact struct {
    Name string
    Email string
    Id int
}

type Contacts = []Contact

type Data struct {
    Contacts Contacts
}

func newContact(name string, email string) Contact {
    id++
    return Contact{
        Name: name,
        Email: email,
        Id: id,
    }
}

func newData() Data {
    return Data{
        Contacts: []Contact{
            newContact("Jerry Phillip", "JP@gmail.com"),
            newContact("Aidan Lows", "AL@gmail.com"),
        },
    }
}

func (d Data) hasEmail(email string) bool {
    for _, contact := range d.Contacts {
        if (contact.Email == email) {
            return true
        }
    }
    return false
}

type FormData struct {
    Values map[string]string
    Errors map[string]string
} 

func newFormData() FormData {
    return FormData{
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

type Page struct {
    Data Data
    Form FormData
}

func newPage() Page {
    return Page {
        Data: newData(),
        Form: newFormData(),
    }
}

func (d *Data) indexOf(id int) int {
    for i, contact := range d.Contacts {
        if contact.Id == id {
            return i
        }
    }
    return -1
}

func main() {
	e := echo.New()
    e.Renderer = NewTemplates()
    e.Use(middleware.Logger())

    page := newPage()

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", page)
    })

    e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        if (page.Data.hasEmail(email)) {
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email

            formData.Errors["email"] = "Email already exists"
            return c.Render(422, "form", formData)
        }

        contact := newContact(name, email)

        page.Data.Contacts = append(page.Data.Contacts, contact)

        c.Render(200, "form", newFormData())
        return c.Render(200, "oob-contact", contact)
    })

    e.GET("/blocks", func(c echo.Context) error {
        startStr := c.QueryParam("start")
        start, err := strconv.Atoi(startStr)
        if err != nil {
            start = 0
        }

        blocks := []Block{}
        for i := start; i < start + 10; i++ {
            blocks = append(blocks, Block{Id: i})
        }

        template := "blocks"
        if start == 0 {
            template = "blocks-index"
        }
        return c.Render(http.StatusOK, template, Blocks{
            Start: start,
            Next: start + 10,
            More: start + 10 < 100,
            Blocks: blocks,
        });
    });

    e.DELETE("/contacts/:id", func(c echo.Context) error {
        time.Sleep(1 * time.Second)
        idStr := c.Param("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            return c.String(400, "Invalid Id")
        }
        index := page.Data.indexOf(id)
        if index == -1 {
            return c.String(404, "Contact not found")
        }
        
        page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)
        return c.NoContent(200)
    })


    e.Static("/css", "css")
    e.Static("/images", "images")
    e.Static("/", "fem-htmx-proj")

    e.Logger.Fatal(e.Start(":42069"))
}
