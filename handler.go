package jsonform

import (
	"context"
	"net/http"
	"strings"

	"github.com/swaggest/rest/web"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Repository) Handler() http.Handler {
	return nil
}

func (r *Repository) Mount(s *web.Service, prefix string) {
	r.baseURL = prefix

	s.Get(prefix+"{name}-schema.json", r.GetSchema())
	s.Mount(prefix, http.StripPrefix(prefix, staticServer))
}

type schemaReq struct {
	Name schemaName `path:"name"`
}

type schemaName string

func (s schemaName) Enum() []interface{} {
	var enum []interface{}

	for _, v := range strings.Split(string(s), ",") {
		enum = append(enum, v)
	}

	return enum
}

func (r *Repository) GetSchema() usecase.Interactor {
	in := schemaReq{
		Name: schemaName(strings.Join(r.Names(), ",")),
	}

	u := usecase.NewIOI(in, new(FormSchema), func(ctx context.Context, in, out interface{}) error {
		input, _ := in.(schemaReq)
		output, _ := out.(*FormSchema)

		if fs, found := r.schemasByName[string(input.Name)]; found {
			*output = fs

			return nil
		}

		return status.NotFound
	})

	u.SetTitle("Get JSONForm Schema")
	u.SetExpectedErrors(status.NotFound)

	return u
}
