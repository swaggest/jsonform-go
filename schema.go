package jsonform

import (
	"fmt"
	"path"
	"reflect"
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

	schemasByName map[string]FormSchema
	namesByType   map[reflect.Type]string

	baseURL string
}

func NewRepository(reflector *jsonschema.Reflector) *Repository {
	r := Repository{}
	r.reflector = reflector
	r.schemasByName = make(map[string]FormSchema)
	r.namesByType = make(map[reflect.Type]string)

	return &r
}

func (r *Repository) Name(value interface{}) (string, error) {
	t := refl.DeepIndirect(reflect.TypeOf(value))

	if name, ok := r.namesByType[t]; ok {
		return name, nil
	}

	name := strings.TrimPrefix(strings.ToLower(path.Base(string(refl.GoType(t)))), "main.")

	return name, r.Add(value, name)
}

// Add registers schema with custom name, this is not needed if default name is good enough.
func (r *Repository) Add(value interface{}, name string) error {
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
	r.schemasByName[name] = fs
	r.namesByType[refl.DeepIndirect(reflect.TypeOf(value))] = name

	return nil
}

func (r *Repository) Schema(value interface{}) (*FormSchema, error) {
	if name, err := r.Name(value); err == nil {
		return r.SchemaByName(name), nil
	} else {
		return nil, err
	}
}

func (r *Repository) SchemaByName(name string) *FormSchema {
	if s, ok := r.schemasByName[name]; ok {
		return &s
	}

	return nil
}

func (r *Repository) Names() []string {
	names := make([]string, 0, len(r.schemasByName))

	for name := range r.schemasByName {
		names = append(names, name)
	}

	return names
}
