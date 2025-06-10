package jsonform

import (
	"fmt"
	"html/template"
	"io"
	"strconv"
)

// Form describes form parameters.
type Form struct {
	// Name is used in form elements identifiers, form number is used for empty name.
	Name          string `json:"name,omitempty"`
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	SchemaName    string `json:"schemaName,omitempty"`
	ValueURL      string `json:"valueUrl,omitempty"`
	SubmitURL     string `json:"submitUrl,omitempty"`
	SubmitMethod  string `json:"submitMethod,omitempty"`
	SuccessStatus int    `json:"successStatus,omitempty"`

	// OnSuccess is a javascript callback that receives XMLHttpRequest value in case of successful response.
	OnSuccess template.JS `json:"-"`
	// OnFail is a javascript callback that receives XMLHttpRequest value in case of a failure response.
	OnFail template.JS `json:"-"`
	// OnError is a javascript callback that receives string HTML value in case of an error while processing the form.
	OnError template.JS `json:"-"`
	// OnBeforeSubmit is a javascript callback that receives form data after Submit button is clicked and before request is sent.
	OnBeforeSubmit template.JS `json:"-"`
	// OnRequestFinished is a javascript callback that receives XMLHttpRequest after request is finished.
	OnRequestFinished template.JS `json:"-"`

	Schema *FormSchema `json:"schema,omitempty"`
	Value  interface{} `json:"value,omitempty"`

	// SubmitText is an optional description of submit button.
	SubmitText string `json:"-"`

	// BeforeForm is injected before form container.
	BeforeForm template.HTML `json:"-"`
	// AfterForm is injected after form container.
	AfterForm template.HTML `json:"-"`
}

// Page allows page customizations.
type Page struct {
	// AppendHTMLHead is injected into the <head> of an HTML document.
	AppendHTMLHead template.HTML

	// PrependHTML is added before the forms.
	PrependHTML template.HTML

	// AppendHTML is added after the forms.
	AppendHTML template.HTML

	// Title is set to HTML document title.
	Title string
}

var formTemplate = loadTemplate("form_tmpl.html")

// Render renders forms as web page.
func (r *Repository) Render(w io.Writer, p Page, forms ...Form) error {
	type pageData struct {
		Page
		Params  []Form
		BaseURL string
	}

	d := pageData{
		Page:    p,
		BaseURL: r.baseURL,
	}

	for i, form := range forms {
		if d.Title == "" {
			d.Title = form.Title
		}

		if form.Name == "" {
			form.Name = strconv.Itoa(i)
		}

		if form.Schema == nil && form.Value != nil {
			s, err := r.formSchema(form.Value)
			if err != nil {
				return err
			}

			form.Schema = &FormSchema{}
			*form.Schema = *s

			submit := FormItem{FormType: "submit", FormTitle: "Submit"}

			if form.OnBeforeSubmit == "" && form.OnRequestFinished == "" {
				form.OnBeforeSubmit = "startSpinner"
				form.OnRequestFinished = "stopSpinner"
			}

			if form.SubmitText != "" {
				submit.FormTitle = form.SubmitText
			}

			form.Schema.Form = append(form.Schema.Form, submit)
		}

		d.Params = append(d.Params, form)
	}

	return formTemplate.Execute(w, d)
}

func (r *Repository) formSchema(value interface{}) (*FormSchema, error) {
	formSchema := r.Schema(value)
	if formSchema == nil {
		if r.Strict {
			return nil, fmt.Errorf("missing form schema for %T", value)
		}

		if err := r.Add(value); err != nil {
			return nil, err
		}

		formSchema = r.Schema(value)
	}

	return formSchema, nil
}
