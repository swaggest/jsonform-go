package jsonform

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
)

/**
 * @typedef formParams
 * @type {Object}
 * @property {String} title - Title of the form.
 * @property {String} schemaName - Schema name.
 * @property {String} valueUrl - URL to fetch value.
 * @property {String} submitUrl - URL to submit form.
 * @property {String} submitMethod - HTTP method to use on form submit.
 */

type FormParams struct {
	Title        string `json:"title,omitempty"`
	SchemaName   string `json:"schemaName,omitempty"`
	ValueURL     string `json:"valueUrl,omitempty"`
	SubmitURL    string `json:"submitUrl,omitempty"`
	SubmitMethod string `json:"submitMethod,omitempty"`

	Schema *FormSchema `json:"schema,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

var formTemplate = loadTemplate("form_tmpl.html")

func (r *Repository) RenderForm(params FormParams, w io.Writer) error {
	type pageData struct {
		Params  template.JS
		BaseURL string
	}

	var err error

	if params.Schema == nil && params.Value != nil {
		params.Schema, err = r.Schema(params.Value)
		if err != nil {
			return fmt.Errorf("render form: %w", err)
		}
	}

	j, err := json.Marshal(params)
	if err != nil {
		return err
	}

	d := pageData{
		Params:  template.JS(j),
		BaseURL: r.baseURL,
	}

	return formTemplate.Execute(w, d)
}
