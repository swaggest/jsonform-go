package jsonform

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/refl"
)

// FormItem defines form item rendering parameters.
type FormItem struct {
	Key       string `json:"key,omitempty" example:"longmood"`
	FormType  string `json:"type,omitempty" examples:"[\"textarea\",\"password\",\"wysihtml5\",\"submit\",\"color\",\"checkboxes\",\"radios\",\"fieldset\", \"help\", \"hidden\"]"`
	FormTitle string `json:"title,omitempty" example:"Submit"`

	ReadOnly bool `json:"readonly,omitempty"`

	Prepend        string            `json:"prepend,omitempty" example:"I feel"`
	Append         string            `json:"append,omitempty" example:"today"`
	NoTitle        bool              `json:"notitle,omitempty"`
	HtmlClass      string            `json:"htmlClass,omitempty" example:"usermood"`
	HtmlMetaData   map[string]string `json:"htmlMetaData,omitempty" example:"{\"style\":\"border: 1px solid blue\",\"data-title\":\"Mood\"}"`
	FieldHtmlClass string            `json:"fieldHtmlClass,omitempty" example:"input-xxlarge"`
	Placeholder    string            `json:"placeholder,omitempty" example:"incredibly and admirably great"`
	InlineTitle    string            `json:"inlinetitle,omitempty" example:"Check this box if you are over 18"`
	TitleMap       map[string]string `json:"titleMap,omitempty" description:"Title mapping for enum."`
	ActiveClass    string            `json:"activeClass,omitempty" example:"btn-success" description:"Button mode for radio buttons."`
	HelpValue      string            `json:"helpvalue,omitempty" example:"<strong>Click me!</strong>"`
}

// FormSchema describes form elements.
type FormSchema struct {
	Form   []FormItem        `json:"form,omitempty"`
	Schema jsonschema.Schema `json:"schema"`
}

// Repository manages form schemas and provides integration helpers.
type Repository struct {
	reflector *jsonschema.Reflector

	schemasByName map[string]FormSchema
	namesByType   map[reflect.Type]string

	baseURL string
}

// NewRepository creates schema repository.
func NewRepository(reflector *jsonschema.Reflector) *Repository {
	r := Repository{}
	r.reflector = reflector
	r.schemasByName = make(map[string]FormSchema)
	r.namesByType = make(map[reflect.Type]string)

	return &r
}

// Name returns schema name by sample value.
func (r *Repository) Name(value interface{}) string {
	t := refl.DeepIndirect(reflect.TypeOf(value))

	if name, ok := r.namesByType[t]; ok {
		return name
	}

	name := strings.TrimPrefix(strings.ToLower(path.Base(string(refl.GoType(t)))), "main.")

	return name
}

// Add adds schemas of value samples.
// It stops on the first error.
func (r *Repository) Add(values ...interface{}) error {
	for _, v := range values {
		if err := r.AddNamed(v, r.Name(v)); err != nil {
			return err
		}
	}

	return nil
}

// AddNamed registers schema with custom name, this is not needed if default name is good enough.
func (r *Repository) AddNamed(value interface{}, name string) error {
	if _, ok := r.schemasByName[name]; ok {
		return fmt.Errorf("schema for %s (%T) is already added", name, value)
	}

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
	r.schemasByName[name] = fs
	r.namesByType[refl.DeepIndirect(reflect.TypeOf(value))] = name

	return nil
}

// Schema returns previously added schema by its sample value.
// It returns nil for unknown schema.
func (r *Repository) Schema(value interface{}) *FormSchema {
	return r.SchemaByName(r.Name(value))
}

// SchemaByName return previously added schema by its name.
// It returns nil for unknown schema.
func (r *Repository) SchemaByName(name string) *FormSchema {
	if s, ok := r.schemasByName[name]; ok {
		return &s
	}

	return nil
}

// Names returns names of added schemas.
func (r *Repository) Names() []string {
	names := make([]string, 0, len(r.schemasByName))

	for name := range r.schemasByName {
		names = append(names, name)
	}

	return names
}
