package view

import (
	"bytes"
	"html/template"
	"net/http"
)

// Renderer renders named views to an http.ResponseWriter.
type Renderer interface {
	Render(w http.ResponseWriter, name string, data any) error
}

type TemplateRenderer struct {
	tmpl *template.Template
}

func NewTemplateRenderer(tmpl *template.Template) *TemplateRenderer {
	return &TemplateRenderer{tmpl: tmpl}
}

func (r *TemplateRenderer) Render(w http.ResponseWriter, name string, data any) error {
	var buf bytes.Buffer
	if err := r.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return err
	}
	_, err := buf.WriteTo(w)
	return err
}
