package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/swaggest/jsonform-go"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v4emb"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type userStatus string

func (us userStatus) Enum() []interface{} {
	return []interface{}{
		"new",
		"approved",
		"active",
		"deleted",
	}
}

// A demo app that receives data from http and stores it in memory.

type User struct {
	FirstName string     `json:"firstName" required:"true" title:"First name" minLength:"3"`
	LastName  string     `json:"lastName" required:"true" title:"Last name" minLength:"3"`
	Locale    string     `json:"locale" title:"User locale" enum:"ru-RU,en-US"`
	Age       int        `json:"age" title:"Age" minimum:"1"`
	Status    userStatus `json:"status" title:"Status"`
	Bio       string     `json:"bio" title:"Bio" description:"A brief description of the person." formType:"textarea"`
}

func (User) Title() string {
	return "User"
}

func (User) Description() string {
	return "User is a sample entity."
}

type userRepo struct {
	st []User
}

func (r *userRepo) create(u User) {
	r.st = append(r.st, u)
}

func (r *userRepo) update(i int, u User) {
	r.st[i-1] = u
}

func (r userRepo) list() []User {
	return r.st
}

func (r userRepo) get(i int) (User, error) {
	if i-1 > len(r.st) {
		return User{}, errors.New("user not found")
	}

	return r.st[i-1], nil
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

func main() {
	ur := &userRepo{}
	ur.create(User{
		FirstName: "John",
		LastName:  "Doe",
		Locale:    "en-US",
		Age:       30,
		Status:    "approved",
		Bio:       "whoa, I never existed!",
	})
	s := web.NewService(openapi31.NewReflector())

	// Init API documentation schema.
	s.OpenAPISchema().SetTitle("Users")
	s.OpenAPISchema().SetDescription("This app showcases a trivial REST API.")
	s.OpenAPISchema().SetVersion("v1.2.3")

	// Add use case handler to router.
	s.Post("/users", createUser(ur))
	s.Get("/users.json", listUsers(ur))
	s.Get("/user/{id}.json", getUser(ur))
	s.Put("/user/{id}.json", updateUser(ur))

	// Swagger UI endpoint at /docs.
	s.Docs("/docs", swgui.New)

	jf := jsonform.NewRepository(s.OpenAPIReflector().JSONSchemaReflector())
	_ = jf.AddWithName(User{}, "user")

	jf.Mount(s, "/json-form/")

	s.Method(http.MethodGet, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`
<div>
<a href="/json-form/form.html?title=Create%20user&amp;schemaUrl=user-schema.json&amp;submitUrl=/users&amp;submitMethod=POST">Create user</a>
</div>

<ul>
`))

		for i, u := range ur.list() {
			_, _ = w.Write([]byte(fmt.Sprintf(`

<li> %s %s
<a href="/json-form/form.html?title=Edit%%20user&amp;schemaUrl=user-schema.json&amp;valueUrl=/user/%d.json&amp;submitUrl=/user/%d.json&amp;submitMethod=PUT">Edit</a><br />
</li>
`, u.FirstName, u.LastName, i+1, i+1)))
		}

		_, _ = w.Write([]byte(`
</ul>
`))
	}))

	// Start server.
	log.Println("JSON Forms at http://localhost:8011/, SwaggerUI docs at http://localhost:8011/docs")

	if err := http.ListenAndServe("localhost:8011", s); err != nil {
		log.Fatal(err)
	}
}
