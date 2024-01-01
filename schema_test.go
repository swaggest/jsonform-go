package jsonform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

type UserWithNeighbors struct {
	User      User   `json:"user" title:"User" description:"The user."`
	Neighbors []User `json:"neighbors,omitempty" title:"Neighbors" description:"A list of neighbors."`
}

func (User) Title() string {
	return "User"
}

func (User) Description() string {
	return "User is a sample entity."
}

func TestRepository_AddSchema(t *testing.T) {
	repo := jsonform.NewRepository(&jsonschema.Reflector{})

	assert.NoError(t, repo.Add(UserWithNeighbors{}))
	assertjson.EqMarshal(t, `{
	  "form":[
		{"key":"user.firstName"},{"key":"user.lastName"},{"key":"user.locale"},
		{"key":"user.age"},{"key":"user.status"},
		{"key":"user.bio","type":"textarea"},
		{
		  "key":"neighbors","type":"array",
		  "items":[
			{
			  "type":"section",
			  "items":[
				{"key":"neighbors[].firstName"},{"key":"neighbors[].lastName"},
				{"key":"neighbors[].locale"},{"key":"neighbors[].age"},
				{"key":"neighbors[].status"},
				{"key":"neighbors[].bio","type":"textarea"}
			  ]
			}
		  ]
		}
	  ],
	  "schema":{
		"properties":{
		  "neighbors":{
			"title":"Neighbors","description":"A list of neighbors.",
			"items":{
			  "title":"User","description":"User is a sample entity.",
			  "required":["firstName","lastName"],
			  "properties":{
				"age":{"title":"Age","minimum":1,"type":"integer"},
				"bio":{
				  "title":"Bio","description":"A brief description of the person.",
				  "type":"string"
				},
				"firstName":{"title":"First name","minLength":3,"type":"string"},
				"lastName":{"title":"Last name","minLength":3,"type":"string"},
				"locale":{"title":"User locale","enum":["ru-RU","en-US"],"type":"string"},
				"status":{
				  "title":"Status","enum":["new","approved","active","deleted"],
				  "type":"string"
				}
			  },
			  "type":"object"
			},
			"type":"array"
		  },
		  "user":{
			"title":"User","description":"The user.",
			"required":["firstName","lastName"],
			"properties":{
			  "age":{"title":"Age","minimum":1,"type":"integer"},
			  "bio":{
				"title":"Bio","description":"A brief description of the person.",
				"type":"string"
			  },
			  "firstName":{"title":"First name","minLength":3,"type":"string"},
			  "lastName":{"title":"Last name","minLength":3,"type":"string"},
			  "locale":{"title":"User locale","enum":["ru-RU","en-US"],"type":"string"},
			  "status":{
				"title":"Status","enum":["new","approved","active","deleted"],
				"type":"string"
			  }
			},
			"type":"object"
		  }
		},
		"type":"object"
	  }
	}`,
		repo.Schema(UserWithNeighbors{}))
}

func TestRepository_Add_arrays(t *testing.T) {
	type BarItem struct {
		Bar string `json:"bar" title:"Bar"`
	}

	type ObjectItem struct {
		Foo  string    `json:"foo" title:"Foo" formType:"textarea"`
		More []string  `json:"more" title:"More"`
		Bars []BarItem `json:"bars" title:"Bars"`
	}

	type My struct {
		Objects []ObjectItem `json:"objects,omitempty" title:"Objects"`
		Strings []string     `json:"strings,omitempty" title:"Strings"`
	}

	repo := jsonform.NewRepository(&jsonschema.Reflector{})
	require.NoError(t, repo.Add(My{}))

	assertjson.EqMarshal(t, `{
	  "form":[
		{
		  "key":"objects","type":"array",
		  "items":[
			{
			  "type":"section",
			  "items":[
				{"key":"objects[].foo","type":"textarea"},{"key":"objects[].more"},
				{
				  "key":"objects[].bars","type":"array",
				  "items":[{"type":"section","items":[{"key":"objects[].bars[].bar"}]}]
				}
			  ]
			}
		  ]
		},
		{"key":"strings"}
	  ],
	  "schema":{
		"properties":{
		  "objects":{
			"title":"Objects",
			"items":{
			  "properties":{
				"bars":{
				  "title":"Bars",
				  "items":{
					"properties":{"bar":{"title":"Bar","type":"string"}},
					"type":"object"
				  },
				  "type":["array","null"]
				},
				"foo":{"title":"Foo","type":"string"},
				"more":{"title":"More","items":{"type":"string"},"type":["array","null"]}
			  },
			  "type":"object"
			},
			"type":"array"
		  },
		  "strings":{"title":"Strings","items":{"type":"string"},"type":"array"}
		},
		"type":"object"
	  }
	}`, repo.Schema(My{}))
}
