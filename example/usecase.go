package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/swaggest/jsonform-go"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func createUserForm(r *jsonform.Repository) usecase.Interactor {
	type another struct {
		Foo  string `json:"foo" required:"true" title:"Foo" minLength:"3"`
		Bar  string `json:"bar" required:"true" title:"Bar" maxLength:"3"`
		Baz1 []int  `json:"baz1" required:"true" items.title:"Item Title" title:"Baz1" items.enum:"123,456,789" minimum:"1"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *usecase.OutputWithEmbeddedWriter) error {
		return r.Render(output.Writer,
			jsonform.Page{
				Title: "Create User and some more",
				AppendHTMLHead: `
<script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.37.1/ace.min.js" type="text/javascript" charset="utf-8"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.37.1/ext-inline_autocomplete.min.js" integrity="sha512-99FE+tBv3oH1/pBRMytEllCvV2Web1lvyeunqcKnJHjiGMFLiwQ6WK6rknV/HYq/esz9i6JuprfgV4senYwtlA==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/ace/1.37.1/ext-language_tools.min.js" integrity="sha512-WPFoecnG6+TbAah6uWLOr+/36IDeXpa9klWh+SFWRzQeVC6x/n7rzODbzctFYL7rqxaDK2pbKj3psH7sws70Ng==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
<script>
window.jsonform_ace_setup = function(setup){
console.log("jsonform_ace_setup", setup)
setup()
}
</script>
`,
				PrependHTML: `<div><h2>prepended</h2></div>`,
				AppendHTML:  `<div><h2>appended</h2></div>`,
			},
			jsonform.Form{
				Title:         "Create User",
				SubmitMethod:  http.MethodPost,
				SubmitURL:     "/users",
				Value:         User{},
				SuccessStatus: http.StatusCreated,
			},
			jsonform.Form{
				Title:        "Another random form",
				SubmitMethod: http.MethodPut,
				SubmitURL:    "/nowhere",
				Value:        another{},
			},
			jsonform.Form{
				Title:        "More random forms",
				SubmitMethod: http.MethodPut,
				SubmitURL:    "/nowhere",
				Value:        another{},
			},
		)
	})

	return u
}

func editUserForm(r *jsonform.Repository, ur *userRepo) usecase.Interactor {
	type in struct {
		ID int `path:"id"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input in, output *usecase.OutputWithEmbeddedWriter) error {
		user, err := ur.get(input.ID)
		if err != nil {
			return err
		}

		err = r.Render(output.Writer, jsonform.Page{}, jsonform.Form{
			Title:         "Update User",
			SubmitMethod:  http.MethodPut,
			SubmitURL:     "/user/" + strconv.Itoa(input.ID) + ".json",
			Value:         user,
			SuccessStatus: http.StatusNoContent,
			OnSuccess:     `function(x){console.log(x);alert("response status: " + x.status)}`,
		})
		if err != nil {
			log.Println(err.Error())
		}

		return err
	})

	return u
}

func nestedImageForm(r *jsonform.Repository) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *usecase.OutputWithEmbeddedWriter) error {
		img := Image{
			Path:   "album/demo/sample.jpg",
			Width:  2629,
			Height: 1351,
			Settings: ImageSettings{
				Description: "Canal Saint-Martin",
				HTTPSources: []string{"https://example.com/sample.jpg"},
			},
		}

		return r.Render(output.Writer, jsonform.Page{
			Title: "Nested Object Reproducer",
			PrependHTML: `<div style="margin:2em">
<p>This form reproduces nested object rendering. With the fix, only nested fields such as <code>settings.description</code> should be rendered, and there should be no extra plain <code>settings</code> input.</p>
<p><a href="/nested-image-schema.json">Open raw generated form schema JSON</a></p>
</div>`,
		}, jsonform.Form{
			Title:         "Edit Image",
			SubmitMethod:  http.MethodPut,
			SubmitURL:     "/nowhere",
			Value:         img,
			SuccessStatus: http.StatusNoContent,
		})
	})

	return u
}

func nestedImageSchema(r *jsonform.Repository) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *json.RawMessage) error {
		s := r.Schema(Image{})
		if s == nil {
			if err := r.Add(Image{}); err != nil {
				return err
			}

			s = r.Schema(Image{})
		}

		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return err
		}

		raw := json.RawMessage(b)
		*output = raw

		return nil
	})

	u.SetTitle("Nested Image Schema")

	return u
}

func createUser(ur *userRepo) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input User, output *struct{}) error {
		ur.create(input)

		return nil
	})
	// Describe use case interactor.
	u.SetTitle("Create User")
	u.SetExpectedErrors(status.InvalidArgument)

	return u
}

func listUsers(ur *userRepo) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *[]User) (err error) {
		*output = ur.list()

		return err
	})
	// Describe use case interactor.
	u.SetTitle("List Users")

	return u
}

func getUser(ur *userRepo) usecase.Interactor {
	type getUserInput struct {
		ID int `path:"id"`
	}

	u := usecase.NewInteractor(func(ctx context.Context, input getUserInput, output *User) (err error) {
		*output, err = ur.get(input.ID)

		return err
	})
	// Describe use case interactor.
	u.SetExpectedErrors(status.InvalidArgument)

	return u
}

func updateUser(ur *userRepo) usecase.Interactor {
	type updateUserInput struct {
		ID int `path:"id"`

		User
	}

	u := usecase.NewInteractor(func(ctx context.Context, input updateUserInput, output *struct{}) error {
		ur.update(input.ID, input.User)

		return nil
	})
	// Describe use case interactor.
	u.SetTitle("Create User")
	u.SetExpectedErrors(status.InvalidArgument)

	return u
}
