package jsonform

import (
	"fmt"
	"strings"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/refl"
)

type FormItem struct {
	Key       string `json:"key,omitempty" example:"longmood"`
	FormType  string `json:"type,omitempty" examples:"[\"textarea\",\"password\",\"wysihtml5\",\"submit\",\"color\",\"checkboxes\",\"radios\",\"fieldset\", \"help\", \"hidden\"]"`
	FormTitle string `json:"title,omitempty" example:"Submit"`

	ReadOnly bool `json:"readonly,omitempty"`

	Prepend        string            `json:"prepend,omitempty" example:"I feel"`
	Append         string            `json:"append,omitempty" example:"today"`
	NoTitle        bool              `json:"notitle,omitempty"`
	HTMLClass      string            `json:"htmlClass,omitempty" example:"usermood"`
	HTMLMetaData   map[string]string `json:"htmlMetaData,omitempty" example:"{\"style\":\"border: 1px solid blue\",\"data-title\":\"Mood\"}"`
	FieldHTMLClass string            `json:"fieldHtmlClass,omitempty" example:"input-xxlarge"`
	Placeholder    string            `json:"placeholder,omitempty" example:"incredibly and admirably great"`
	InlineTitle    string            `json:"inlinetitle,omitempty" example:"Check this box if you are over 18"`
	TitleMap       map[string]string `json:"titleMap,omitempty" description:"Title mapping for enum."`
	ActiveClass    string            `json:"activeClass,omitempty" example:"btn-success" description:"Button mode for radio buttons."`
	HelpValue      string            `json:"helpvalue,omitempty" example:"<strong>Click me!</strong>"`
}
type FormSchema struct {
	Form   []FormItem        `json:"form,omitempty"`
	Schema jsonschema.Schema `json:"schema"`
}

type Repository struct {
	reflector *jsonschema.Reflector
	schemas   map[string]FormSchema
}

func NewRepository(reflector *jsonschema.Reflector) *Repository {
	r := Repository{}
	r.reflector = reflector
	r.schemas = make(map[string]FormSchema)

	return &r
}

func (r *Repository) AddSchema(name string, value interface{}) error {
	fs := FormSchema{}

	schema, err := r.reflector.Reflect(value, jsonschema.InlineRefs, jsonschema.InterceptProp(
		func(params jsonschema.InterceptPropParams) error {
			if !params.Processed {
				return nil
			}

			if params.PropertySchema.HasType(jsonschema.Object) {
				return nil
			}

			fi := FormItem{
				Key: strings.Join(append(params.Path[1:], params.Name), "."),
			}

			fi.Key = strings.ReplaceAll(fi.Key, ".[]", "[]")

			// println(fi.Key)
			if err := refl.PopulateFieldsFromTags(&fi, params.Field.Tag); err != nil {
				return err
			}

			fs.Form = append(fs.Form, fi)

			return nil
		},
	))
	if err != nil {
		return fmt.Errorf("reflecting %s schema: %w", name, err)
	}

	// Complying with Draft 3.
	for _, name := range schema.Required {
		if prop, ok := schema.Properties[name]; ok {
			prop.TypeObject.WithExtraPropertiesItem("required", true)
		}
	}

	schema.Required = nil

	fs.Form = append(fs.Form, FormItem{FormType: "submit", FormTitle: "Submit"})
	fs.Schema = schema
	r.schemas[name] = fs

	return nil
}

func (r *Repository) Schema(name string) FormSchema {
	return r.schemas[name]
}
