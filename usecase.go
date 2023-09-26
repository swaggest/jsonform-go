package jsonform

import (
	"encoding/json"
	"html/template"
	"io"
)

// Form describes form parameters.
type Form struct {
	Title         string `json:"title,omitempty"`
	SchemaName    string `json:"schemaName,omitempty"`
	ValueURL      string `json:"valueUrl,omitempty"`
	SubmitURL     string `json:"submitUrl,omitempty"`
	SubmitMethod  string `json:"submitMethod,omitempty"`
	SuccessStatus int    `json:"successStatus,omitempty"`

	Schema *FormSchema `json:"schema,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

// Page allows page customizations.
type Page struct {
	PrependHTML template.HTML
	AppendHTML  template.HTML
	Title       string
}

var formTemplate = loadTemplate("form_tmpl.html")

// Render renders forms as web page.
func (r *Repository) Render(w io.Writer, p Page, forms ...Form) error {
	type pageData struct {
		Page
		Params  []template.JS
		BaseURL string
	}

	d := pageData{
		Page:    p,
		BaseURL: r.baseURL,
	}

	for _, form := range forms {
		if d.Title == "" {
			d.Title = form.Title
		}

		if form.Schema == nil && form.Value != nil {
			form.Schema = r.Schema(form.Value)
			if form.Schema == nil {
				if err := r.Add(form.Value); err != nil {
					return err
				}
			}
		}

		j, err := json.Marshal(form)
		if err != nil {
			return err
		}

		d.Params = append(d.Params, template.JS(j)) //nolint:gosec
	}

	return formTemplate.Execute(w, d)
}
