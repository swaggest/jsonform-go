package jsonform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/assertjson"
	"github.com/swaggest/jsonform-go"
	"github.com/swaggest/jsonschema-go"
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

func TestRepository_AddSchema(t *testing.T) {
	repo := jsonform.NewRepository(&jsonschema.Reflector{})

	assert.NoError(t, repo.Add(User{}, "user"))
	assertjson.EqMarshal(t, `
		{
		  "form":[
			{"key":"firstName"},{"key":"lastName"},{"key":"locale"},{"key":"age"},
			{"key":"status"},{"key":"bio","type":"textarea"},
			{"type":"submit","title":"Submit"}
		  ],
		  "schema":{
			"title":"User","description":"User is a sample entity.",
			"properties":{
			  "age":{"title":"Age","minimum":1,"type":"integer"},
			  "bio":{
				"title":"Bio","description":"A brief description of the person.",
				"type":"string"
			  },
			  "firstName":{"title":"First name","minLength":3,"type":"string","required":true},
			  "lastName":{"title":"Last name","minLength":3,"type":"string","required":true},
			  "locale":{"title":"User locale","enum":["ru-RU","en-US"],"type":"string"},
			  "status":{
				"title":"Status","enum":["new","approved","active","deleted"],
				"type":"string"
			  }
			},
			"type":"object"
		  }
		}`,
		repo.Schema(User{}))
}
