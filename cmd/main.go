package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

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

type Count struct {
    Count int
}

func main() {
	e := echo.New()
    e.Renderer = NewTemplates()
    e.Use(middleware.Logger())

    count := Count { Count: 0 }

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", count)
    })

    e.POST("/count", func(c echo.Context) error {
        count.Count++;
        return c.String(http.StatusOK, strconv.Itoa(count.Count))
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

    e.Logger.Fatal(e.Start(":42069"))
}
