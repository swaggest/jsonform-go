package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/swaggest/jsonform-go"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func createUserForm(r *jsonform.Repository) usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *usecase.OutputWithEmbeddedWriter) error {
		return r.RenderForm(jsonform.FormParams{
			Title:        "Create User",
			SubmitMethod: http.MethodPost,
			SubmitURL:    "/users",
			Value:        User{},
		}, output.Writer)
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

		return r.RenderForm(jsonform.FormParams{
			Title:        "Create User",
			SubmitMethod: http.MethodPut,
			SubmitURL:    "/user/" + strconv.Itoa(input.ID) + ".json",
			Value:        user,
		}, output.Writer)
	})

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
