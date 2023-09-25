package main

import (
	"log"
	"net/http"

	"github.com/swaggest/jsonform-go"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
)

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

	jf := jsonform.NewRepository(s.OpenAPIReflector().JSONSchemaReflector())
	ur.schemaName, _ = jf.Name(User{})

	// Add use case handler to router.
	s.Post("/users", createUser(ur))
	s.Get("/users.json", listUsers(ur))
	s.Get("/user/{id}.json", getUser(ur))
	s.Put("/user/{id}.json", updateUser(ur))

	// Static forms.
	s.Get("/create-user", createUserForm(jf))
	s.Get("/edit-user/{id}", editUserForm(jf, ur))

	// Swagger UI endpoint at /docs.
	s.Docs("/docs", swgui.New)

	jf.Mount(s, "/json-form/")
	s.Method(http.MethodGet, "/", ur)

	// Start server.
	log.Println("JSON Forms at http://localhost:8011/, SwaggerUI docs at http://localhost:8011/docs")

	if err := http.ListenAndServe("localhost:8011", s); err != nil {
		log.Fatal(err)
	}
}
