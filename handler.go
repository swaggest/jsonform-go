package jsonform

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/swaggest/rest/web"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

// Mount attaches handlers to web service.
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
	ss := strings.Split(string(s), ",")
	enum := make([]interface{}, 0, len(ss))

	for _, v := range ss {
		enum = append(enum, v)
	}

	return enum
}

// GetSchema returns JSONForm schema.
func (r *Repository) GetSchema() usecase.Interactor {
	in := schemaReq{
		Name: schemaName(strings.Join(r.Names(), ",")),
	}

	u := usecase.NewIOI(in, new(FormSchema), func(ctx context.Context, in, out interface{}) error {
		input, ok := in.(schemaReq)
		if !ok {
			return fmt.Errorf("unexpected input: %T", in)
		}

		output, ok := out.(*FormSchema)
		if !ok {
			return fmt.Errorf("unexpected output: %T", out)
		}

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
