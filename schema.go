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
	reflector     *jsonschema.Reflector
	schemasByName map[string]FormSchema
	namesByType   map[reflect.Type]string
}

func NewRepository(reflector *jsonschema.Reflector) *Repository {
	r := Repository{}
	r.reflector = reflector
	r.schemasByName = make(map[string]FormSchema)
	r.namesByType = make(map[reflect.Type]string)

	return &r
}

func (r *Repository) Add(value interface{}) error {
	name := strings.ToLower(path.Base(string(refl.GoType(refl.DeepIndirect(reflect.TypeOf(value)))))) // Just right amount of brackets.

	return r.AddWithName(value, name)
}

func (r *Repository) Name(value interface{}) string {
	return r.namesByType[refl.DeepIndirect(reflect.TypeOf(value))]
}

func (r *Repository) AddWithName(value interface{}, name string) error {
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

func (r *Repository) Schema(value interface{}) *FormSchema {
	return r.schemaByName(r.Name(value))
}

func (r *Repository) schemaByName(name string) *FormSchema {
	if s, ok := r.schemasByName[name]; ok {
		return &s
	}

	return nil
}
